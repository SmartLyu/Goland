#!/bin/bash
# ==================
# Description: 巡查调用nat使用的程序
# Created By: 于志远
# Version: 1.2
# Last Modified: 2020-01-06
# ==================

# 全局变量
TIMEOUT='2'
USER=work

STEPLENTH='60'
PATROLIP="patrol.ijunhai.com"
URL="${PATROLIP}:8686/monitor/collect"
HOSTSURL="${PATROLIP}:8686/monitor/nat"
DOWNLOAD_URL="${PATROLIP}:8686/shell/monitor"
HOST=${HOSTNAME}
DATE=$(TZ=Asia/Shanghai date "+%Y-%m-%d %H:%M")
TIMETEMPLE=""

# 日志信息
LOGDIR="/work/logs/openfalcon"
LOGFILE="${LOGDIR}/patrol-nat-monitor.log"
RETRYLOGFILE="${LOGDIR}/retry-nat.log"
TMPLOGFILE="${LOGDIR}/patrol-nat-detail.$(TZ=Asia/Shanghai date "+%Y-%m-%d-%H-%M").log"
TMPSTATUSLOGFILE="/tmp/.patrol-nat.status.tmp.log"

# 海外加速配置
OVERSEA_PATROLIP="120.24.170.21"
OVERSEA_URL="${PATROLIP}:8686/monitor/collect"
OVERSEA_HOSTSURL="${PATROLIP}:8686/monitor/nat"
OVERSEA_DOWNLOAD_URL="${PATROLIP}:8686/shell/monitor"

# 日志裁剪阈值
MAXSIZE=10485760
NATINFO=""

# 发送数据给监控的API
# 参数1：检查项  参数2：状态
PostToApi(){
    local ip_info
    local tag_name
    local status
    local data

    ip_info=${NATINFO}" nat boot "
    tag_name=$1
    status=$2

    # 判断脚本是否需要保留日志
    if [[ $status == false ]];then
        echo false > ${TMPSTATUSLOGFILE}
    fi

    data="{ \
        \"time\": \"${DATE}\", \
        \"IP\": \"${ip_info}\", \
        \"hostname\": \"${HOST}\", \
        \"info\": \"${tag_name}\", \
        \"status\": ${status} \
        }"
    echo "$(date) ${DATE} ${ip_info}:${tag_name}  ${status}  " >> ${LOGFILE}
    curl -s --connect-timeout ${TIMEOUT} -X POST -d "${data}" $URL
    if [[ $? -ne 0 ]];then
        sleep 2
        echo First $(date) ${data} >> ${RETRYLOGFILE}
        curl -s --connect-timeout ${TIMEOUT} -X POST -d "${data}" $URL 2>> ${LOGFILE}
        if [[ $? -ne 0 ]];then
            sleep 2
            echo Second $(date) ${data} >> ${RETRYLOGFILE}
            curl -s --connect-timeout ${TIMEOUT} -X POST -d "${data}" ${OVERSEA_URL} 2>> ${LOGFILE}
        fi
    fi
    echo
}

# 发送数据给监控的API
# 参数1： 后端IP信息   参数2： 主机名   参数3：检查项  参数4：状态
PostTelnetToApi(){
    local ip_info
    local tag_name
    local status
    local data

    ip_info=${NATINFO}"=}"${1}
    tag_name=$3
    status=$4

    data="{ \
        \"time\": \"${DATE}\", \
        \"IP\": \"${ip_info}\", \
        \"hostname\": \"${2}\", \
        \"info\": \"${tag_name}\", \
        \"status\": ${status} \
        }"
    echo "$(date) ${DATE} ${ip_info}:${tag_name}  ${status}  " >> ${LOGFILE}
    curl -s --connect-timeout ${TIMEOUT} -X POST -d "${data}" $URL
    if [[ $? -ne 0 ]];then
        sleep 2
        echo First $(date) ${data} >> ${RETRYLOGFILE}
        curl -s --connect-timeout ${TIMEOUT} -X POST -d "${data}" $URL 2>> ${LOGFILE}
        if [[ $? -ne 0 ]];then
            sleep 2
            echo Second $(date) ${data} >> ${RETRYLOGFILE}
            curl -s --connect-timeout ${TIMEOUT} -X POST -d "${data}" ${OVERSEA_URL} 2>> ${LOGFILE}
        fi
    fi
}

# 发送hosts数据给监控的API
# 参数1：后端服务器IP  参数2：后端服务器简单标识名
PostHostsToApi(){
    local ip_info
    local data

    ip_info=${NATINFO}"=}"${1}
    hosname_info=$2

    data="{ \
        \"IP\": \"${ip_info}\", \
        \"hostname\": \"${hosname_info}\" \
        }"
    echo "$(date) ${DATE} post ${ip_info}  " >> ${LOGFILE}
    curl -s --connect-timeout ${TIMEOUT} -X POST -d "${data}" $HOSTSURL
    if [[ $? -ne 0 ]];then
        echo First $(date) ${data} >> ${RETRYLOGFILE}
        curl -s --connect-timeout ${TIMEOUT} -X POST -d "${data}" $HOSTSURL 2>> ${LOGFILE}
        if [[ $? -ne 0 ]];then
            echo Third $(date) ${data} >> ${RETRYLOGFILE}
            curl -s --connect-timeout ${TIMEOUT} -X POST -d "${data}" ${OVERSEA_HOSTSURL} 2>> ${LOGFILE}
            if [[ $? -ne 0 ]];then
                curl -s --connect-timeout ${TIMEOUT} -X POST -d "${data}" ${OVERSEA_HOSTSURL} 2>> ${LOGFILE}
                if [[ $? != 0 ]];then
                    PatrolStatus="false"
                fi
            fi
        fi
    fi
    echo
}

# 检查服务器的指定端口内网状态
# 参数1：IP  参数2： 端口
MonitorPort(){
    local ip
    local status
    local hostname
    local port
    ip=$1
    port=$2
    hostname=$(ssh $ip hostname)
    if [[ ${hostname} == '' ]];then
        hostname=$(ssh -p 45678 $ip hostname)
    fi

    # 错开巡查机器上传数据的时间
    sleep $[RANDOM%5].$[RANDOM%999]

    echo "" | timeout 3 telnet ${ip} ${port} | grep Connected
    if [[ $? -eq 0 ]];then
        status=true
    else
        echo "" | timeout 3 telnet ${ip} ${port} | grep Connected
        if [[ $? -eq 0 ]];then
            status=true
        else
            status=false
        fi
    fi

    PostTelnetToApi ${ip} ${hostname} port_${port} ${status}
}

# 错误日志记录并警报退出
# 参数1：记录异常数
MyError(){
    local error_num
    error_num=$1
    PostToApi "monitor-${error_num}"  false
    echo $(date) ${error_num} >> ${LOGFILE}
    exit ${error_num}
}

# 下载是否异常判断
# 参数1：异常数
DownloadError(){
    PostToApi "download_sshfile" $1
    if [[ $1 != 'true' ]];then
        exit 5
    fi
}

# 简单测试环境并完成初始化
Init(){
    if [[ s${NATINFO} == "s" ]] ;then
         MyError 1
    fi
    if [[ $? -ne 0 ]];then
        echo you have to init enviroment!
        MyError 2
    fi
    which telnet
    if [[ $? -ne 0 ]];then
        echo you have to install telnet!
        PostToApi telnet false
        MyError 3
    fi
    touch ${LOGFILE}
    if [[ ! -w ${LOGFILE} ]];then
        LOGFILE='/tmp/tmp-test.log'
        PostToApi ${LOGFILE}_permission_deny false
        MyError 4
    fi

}

# 遍历hosts进行监控
Monitor(){
    local all_ip
    local special_tag
    wget "${DOWNLOAD_URL}" --timeout 3 -O /tmp/patrol-tmp.sh ||\
    wget "${DOWNLOAD_URL}" --timeout 5 -O /tmp/patrol-tmp.sh ||\
    wget "${OVERSEA_DOWNLOAD_URL}" --timeout 4 -O /tmp/patrol-tmp.sh ||\
    DownloadError false

    DownloadError true
    all_ip=$(awk '{print $1}' /etc/hosts | grep -v ^# |  grep -v ^$ | grep -E ^\(10\)\\..*\|\(172\)\\..*\|\(192\)\\..* | sort | uniq)

    for i in ${all_ip}
    do
        echo $i
        {
            # 检查集群身份
            special_tag=""
            i_hostname=$(grep $i /etc/hosts | awk '{print $2}')

            # 记录后端服务器是否检查所有
            is_check_all=1

            # 不属于集群服务器
            grep $i /etc/hosts | grep PATROL_PASS
            if [[ $? -eq 0 ]];then
                continue
            fi

            echo $i "/bin/bash /tmp/patrol-tmp.sh --nat "${NATINFO}" --ip "${i}" -af"${special_tag}
            # 告知巡查机本次检查后端机器信息
            PostHostsToApi ${i} ${i_hostname}

            # 不具体检查服务器
            grep $i' ' /etc/hosts | grep PATROL_IGNORE
            if [[ $? -eq 0 ]];then
                special_tag=${special_tag}"f"
                is_check_all=0
            fi

            # 不具体检查监控信息
            grep $i' ' /etc/hosts | grep PATROL_JUST
            if [[ $? -eq 0 ]];then
                special_tag=${special_tag}"a"
                is_check_all=0
            fi

            # nginx服务器
            grep $i' ' /etc/hosts | grep PATROL_NGINX
            if [[ $? -eq 0 ]];then
                MonitorPort $i 80 &
                special_tag=${special_tag}"x"
            fi

            # php服务器
            grep $i' ' /etc/hosts | grep PATROL_PHP
            if [[ $? -eq 0 ]];then
                special_tag=${special_tag}"p"
            fi

            # mysql数据库(非从库)
            grep $i' ' /etc/hosts | grep PATROL_MYSQL
            if [[ $? -eq 0 ]];then
                MonitorPort $i 3306 &
                special_tag=${special_tag}"m"
            fi

            # mysql数据库(从库)
            grep $i' ' /etc/hosts | grep PATROL_MYSQL_SLAVE
            if [[ $? -eq 0 ]];then
                special_tag=${special_tag}"s"
            fi

            # 端口检查
            grep $i' ' /etc/hosts | grep PATROL_PORT
            if [[ $? -eq 0 ]];then
                local num=2
                while :
                do
                    port=$(grep $i' ' /etc/hosts | awk -F"PATROL_PORT_" '{print $'${num}'}' | awk '{print $1}')
                    if [[ 's'${port} == 's' ]];then
                        break
                    fi
                    MonitorPort $i ${port} &
                    num=$[num+1]
                done
            fi

            # 判断是否检查基础信息和监控状态
            if [[ $is_check_all == 1 ]];then
                special_tag=${special_tag}"af"
            fi

            # 控制后端机器进行下载监控脚本
            scp /tmp/patrol-tmp.sh $i:/tmp/patrol-tmp.sh || \
            scp -P 45678 /tmp/patrol-tmp.sh $i:/tmp/patrol-tmp.sh
            # 检查脚本是否存在
            # ssh $i "[ ! -f /tmp/patrol-tmp.sh ] && wget ${DOWNLOAD_URL} --timeout 3 -O /tmp/patrol-tmp.sh" || \
            # ssh -p 45678 $i "[ ! -f /tmp/patrol-tmp.sh ] && wget ${DOWNLOAD_URL} --timeout 3 -O /tmp/patrol-tmp.sh"
            # 执行脚本
            ssh $i "/bin/bash /tmp/patrol-tmp.sh --nat "${NATINFO}" --ip "${i}" --time "${TIMETEMPLE}" -"${special_tag} || \
            ssh -p 45678 $i "/bin/bash /tmp/patrol-tmp.sh --nat "${NATINFO}" --ip "${i}" --time "${TIMETEMPLE}" -"${special_tag} || \
            echo false > ${TMPSTATUSLOGFILE}
        } &
        sleep 0.01
    done
    MonitorSelf &
}

# 检查nat机器自身
MonitorSelf(){
    /bin/bash /tmp/patrol-tmp.sh --nat "${NATINFO}" --ip "127.0.0.1" --time "${TIMETEMPLE}" -a;
}

Main(){
    echo true > ${TMPSTATUSLOGFILE}
    echo $(date) $@
    echo "$(date) prepare to patrol all hosts." >> ${LOGFILE}
    local short_opts="hn:t:"
    local long_opts="help,nat:,time:"
    local args
    # 将规范化后的命令行参数分配至位置参数（$1,$2,...)
    args=$(getopt -o ${short_opts} --long ${long_opts} -- "$@" 2>/dev/null)

    if [[ $? -ne 0 || $# -eq 0 ]]
    then
        Usage
        exit 1
    fi
    eval set -- "${args}"
    while true
    do
        case "$1" in
            -h|--help)
                Usage
                shift
                ;;
            -n|--nat)
                NATINFO=$2"-"${HOSTNAME}
                shift 2
                ;;
            -t|--time)
                TIMETEMPLE=$2
                DATE=$(TZ=Asia/Shanghai date  -d "@${TIMETEMPLE}" "+%Y-%m-%d %H:%M")
                shift 2
                ;;
            --)
                shift
                break
                ;;
            *)
                Usage
                exit 1
                ;;
    esac
    done

    # 获取基本参数后，开始监控模块
    Init
    find ${LOGDIR} -name 'patrol-nat-detail.*.log' -mmin +1000 -exec rm -f {} \;
    Monitor
    wait
    if [[ $(cat ${TMPSTATUSLOGFILE}) == 'true' ]];then
        rm -f ${TMPLOGFILE}
    fi
    rm -f ${TMPSTATUSLOGFILE}
}

# 显示脚本用法
Usage(){
    cat <<EOF

USAGE:$0 [OPTIONS] [work_password]

选择安装模式：
    -h | --help          查看帮助信息
    -n | --nat           指定本机IP地址
    -t | --time          指定当前时间信息
EOF
}

mkdir -p ${LOGDIR}
Main $@ &> ${TMPLOGFILE} &
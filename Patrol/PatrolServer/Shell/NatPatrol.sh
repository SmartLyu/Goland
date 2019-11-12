#!/bin/bash
# ==================
# Description: 巡查调用nat使用的程序
# Created By: 于志远
# Version: 0.1
# Last Modified: 2019-5-14
# ==================

# 全局变量
TIMEOUT='5'
USER=work

STEPLENTH='60'
URL="patrol.ijunhai.com:8686/monitor/collect"
HOSTSURL="patrol.ijunhai.com:8686/monitor/nat"
HOST=${HOSTNAME}
LOGDIR="/work/logs/openfalcon"
LOGFILE="${LOGDIR}/patrol-nat-monitor.log"
DOWNLOAD_URL="134.175.50.184:8686/shell/monitor"
DATE=$(TZ=Asia/Shanghai date "+%Y-%m-%d %H:%M")
MAXSIZE=10485760

# 详细日志位置
DETAILLOG="${LOGDIR}/patrol.detail.log"
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

    data="{ \
        \"time\": \"${DATE}\", \
        \"IP\": \"${ip_info}\", \
        \"hostname\": \"${HOST}\", \
        \"info\": \"${tag_name}\", \
        \"status\": ${status} \
        }"
    echo "$(date) ${ip_info}:${tag_name}  ${status}  " >> ${LOGFILE}
    curl -s --connect-timeout ${TIMEOUT} -X POST -d "${data}" $URL
    if [[ $? -ne 0 ]];then
        sleep 2
        echo First $(date) ${data} >> /tmp/retry.log
        curl -s --connect-timeout ${TIMEOUT} -X POST -d "${data}" $URL 2>> ${LOGFILE}
        if [[ $? -ne 0 ]];then
            sleep 2
            echo Second $(date) ${data} >> /tmp/retry.log
            curl -s --connect-timeout ${TIMEOUT} -X POST -d "${data}" $URL 2>> ${LOGFILE}
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
    echo "$(date) ${ip_info}:${tag_name}  ${status}  " >> ${LOGFILE}
    curl -s --connect-timeout ${TIMEOUT} -X POST -d "${data}" $URL
    if [[ $? -ne 0 ]];then
        sleep 2
        echo First $(date) ${data} >> /tmp/retry.log
        curl -s --connect-timeout ${TIMEOUT} -X POST -d "${data}" $URL 2>> ${LOGFILE}
        if [[ $? -ne 0 ]];then
            sleep 2
            echo Second $(date) ${data} >> /tmp/retry.log
            curl -s --connect-timeout ${TIMEOUT} -X POST -d "${data}" $URL 2>> ${LOGFILE}
        fi
    fi
}

# 发送hosts数据给监控的API
# 参数1：后端服务器IP  参数2：后端服务器简单标识名
PostHostsToApi(){
    local ip_info
    local data

    ip_info=${NATINFO}"=}"${1}

    data="{ \
        \"IP\": \"${ip_info}\" \
        }"
    echo "$(date) post ${ip_info}  " >> ${LOGFILE}
    curl -s --connect-timeout ${TIMEOUT} -X POST -d "${data}" $HOSTSURL
    if [[ $? -ne 0 ]];then
        echo First $(date) ${data} >> /tmp/retry.log
        curl -s --connect-timeout ${TIMEOUT} -X POST -d "${data}" $HOSTSURL 2>> ${LOGFILE}
        if [[ $? -ne 0 ]];then
            echo Second $(date) ${data} >> /tmp/retry.log
            curl -s --connect-timeout ${TIMEOUT} -X POST -d "${data}" $HOSTSURL 2>> ${LOGFILE}
            if [[ $? -ne 0 ]];then
                PostToApi "nat_post_host=${1}-${2}" false &
                echo ; return
            fi
        fi
    fi
    PostToApi "nat_post_host=${1}-${2}" true &
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
        echo $(date) 3 >> ${LOGFILE}
    fi
    touch ${LOGFILE}
}

# 遍历hosts进行监控
Monitor(){
    local all_ip
    local special_tag
    wget "${DOWNLOAD_URL}" --timeout 5 -O /tmp/patrol-tmp.sh ||\
    wget "${DOWNLOAD_URL}" --timeout 5 -O /tmp/patrol-tmp.sh ||\
    DownloadError false

    DownloadError true
    all_ip=$(awk '{print $1}' /etc/hosts | grep -v ^# |  grep -v ^$ | grep -E ^\(10\)\\..*\|\(172\)\\..*\|\(192\)\\..*)

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

            # 控制后端机器进行下载监控脚本并执行
            scp /tmp/patrol-tmp.sh $i:/tmp/patrol-tmp.sh || \
            scp -P 45678 /tmp/patrol-tmp.sh $i:/tmp/patrol-tmp.sh
            ssh $i "mkdir -p ${LOGDIR} ;\
                    /bin/bash /tmp/patrol-tmp.sh --nat "${NATINFO}" --ip "${i}" -"${special_tag}" &> ${DETAILLOG};\
                    rm -f /tmp/patrol-tmp.sh" || \
            ssh -p 45678 $i "/bin/bash /tmp/patrol-tmp.sh --nat "${NATINFO}" --ip "${i}" -"${special_tag}" &> ${DETAILLOG};\
                    rm -f /tmp/patrol-tmp.sh"
        } &
        sleep 0.01
    done
    MonitorSelf &
}

# 检查nat机器自身
MonitorSelf(){
    /bin/bash /tmp/patrol-tmp.sh --nat "${NATINFO}" --ip "127.0.0.1" -a &> ${DETAILLOG};
}

Main(){
    echo "$(date) prepare to patrol all hosts." >> ${LOGFILE}
    local short_opts="hn:"
    local long_opts="help,nat:"
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
    if [[ $(ls -l ${LOGFILE} | awk '{print $5}') -ge ${MAXSIZE} ]];then
        echo "$(tail -10000 ${LOGFILE})" > ${LOGFILE}
    fi
    Monitor
}

# 显示脚本用法
Usage(){
    cat <<EOF

USAGE:$0 [OPTIONS] [work_password]

选择安装模式：
    -h | --help          查看帮助信息
    -n | --nat           指定本机IP地址
EOF
}

mkdir -p ${LOGDIR}
Main $@ &> /work/logs/openfalcon/patrol-nat.detail.log &
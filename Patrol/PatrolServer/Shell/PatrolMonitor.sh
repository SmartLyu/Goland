#!/bin/bash
# ==================
# Description: 巡查使用监控程序
# Created By: 于志远
# Version: 1.2
# Last Modified: 2020-01-06
# ==================

# 全局变量
TIMEOUT='2'
USER=work

STEPLENTH='60'
URL="patrol.ijunhai.com:8686/monitor/collect"
HOST=${HOSTNAME}
DATE=$(TZ=Asia/Shanghai date "+%Y-%m-%d %H:%M")
MAXSIZE=10485760

# 日志信息
LOGDIR="/work/logs/openfalcon"
LOGFILE="${LOGDIR}/patrol-monitor.log"
RETRYLOGFILE="${LOGDIR}/retry.log"
TMPLOGFILE="${LOGDIR}/patrol-detail.$(TZ=Asia/Shanghai date "+%Y-%m-%d-%H-%M").log"
TMPCPULOGFILE="/tmp/.patrol.cpu.tmp.$(TZ=Asia/Shanghai date "+%s").log"
TMPMYSQLLOGFILE="/tmp/.patrol.mysql.tmp.$(TZ=Asia/Shanghai date "+%s").log"
TMPSTATUSLOGFILE="/tmp/.patrol.status.tmp.$(TZ=Asia/Shanghai date "+%s").log"

# 海外加速连接通道
OVERSEA_URL="120.24.170.21:8686/monitor/collect"

#设置报警阈值
POLICE=90

NATINFO=""
IPINFO=""

# 发送数据给监控的API
# 参数1：检查项  参数2：状态
PostToApi(){
    local ip_info
    local tag_name
    local status
    local data

    ip_info=${NATINFO}"=}"${IPINFO}
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
    echo "$(date) : ${DATE} ${ip_info}:${tag_name}  ${status}  " >> ${LOGFILE}
    curl -s --connect-timeout ${TIMEOUT} -X POST -d "${data}" $URL
    if [[ $? -ne 0 ]];then
        sleep 1
        echo First $(date) ${data} >> ${RETRYLOGFILE}
        curl -s --connect-timeout ${TIMEOUT} -X POST -d "${data}" $URL
        if [[ $? -ne 0 ]];then
            sleep 1
            echo Second $(date)  ${data} >> ${RETRYLOGFILE}
            curl -s --connect-timeout ${TIMEOUT} -X POST -d "${data}" $URL
            if [[ $? -ne 0 ]];then
                echo Third $(date)  ${data} >> ${RETRYLOGFILE}
                curl -s --connect-timeout ${TIMEOUT} -X POST -d "${data}" ${OVERSEA_URL}
                if [[ $? -ne 0 ]];then
                    echo $(date) Push data ${tag_name} Error ! >> ${LOGFILE}
                    return
                fi
            fi
        fi
    fi
}

CheckCpu(){
    local status
    local cpu
    local cpu_number

    top -bn 1 &> ${TMPCPULOGFILE}
    cat ${TMPCPULOGFILE}
    ps -aux
    cpu_number=$(cat /proc/cpuinfo| grep "processor"| wc -l)
    if [[ ${cpu_number} == 1 ]] || [[ ${cpu_number} == 2 ]];then
        return
    fi
    if [[ -f ${TMPCPULOGFILE} ]];then
        cpu=$(cat ${TMPCPULOGFILE} |grep Cpu | tail -1 | cut -d "," -f 1 | cut -d ":" -f 2 | awk '{print $1}');
    else
        top -bn 1
        cpu=$(top -bn 1 |grep Cpu | tail -1 | cut -d "," -f 1 | cut -d ":" -f 2 | awk '{print $1}');
    fi

    echo $cpu | grep '%'
    if [[ $? -eq 0 ]];then
        cpu=$(echo $cpu | awk -F% '{print $1}')
    fi

    echo $cpu | grep '\.'
    if [[ $? -eq 0 ]];then
        cpu=$(echo $cpu | awk -F. '{print $1}')
    fi
    echo cpu=$cpu

    if [[ $cpu -ge $POLICE ]] || [[ s$cpu == 's' ]];then
        status=false
    else
        status=true
    fi

    PostToApi cpu_used=${cpu}% ${status}
}

CheckLoad(){
    local status
    local load
    local cpu_number

    uptime
    cpu_number=$(cat /proc/cpuinfo| grep "processor"| wc -l)
    if [[ ${cpu_number} == 1 ]] || [[ ${cpu_number} == 2 ]];then
        return
    fi
    load_1=$(uptime | awk '{printf ("%d\n",$(NF-2)/'${cpu_number}'*6)}')
    load_5=$(uptime | awk '{printf ("%d\n",$(NF-1)/'${cpu_number}'*6)}')
    load=$[load_1+load_5]
    load=$[load/2]
    if [[ load -ge $POLICE ]] || [[ s$load == 's' ]];then
        status=false
    else
        status=true
    fi

    PostToApi load_used=${load}% ${status}
}

CheckMem(){
    local status
    local mem
    local memMem
    local memTotal

    free -mht
    if [[ $(free -t | grep Total | awk '{print($2)}') -le 1500000 ]];then
       return
    fi
    memMem=$(free -t | grep Mem | awk '{printf ("%d\n",($4+$5+$6)/$2*100)}')
    memTotal=$(free -t | grep Total | awk '{printf ("%d\n",$3/$2*100)}')
    mem=$[(memMem+memTotal)/2]
    if [[ ${mem} -ge $POLICE ]] || [[ s$mem == 's' ]];then
        status=false
    else
        status=true
    fi

    PostToApi mem_used=${mem}% ${status}
}

CheckDf(){
    local status
    local df

    df -h
    df=$(df | grep ^/ | grep -v 'loop' | awk '{print $1}')
    for i in $df
    do
        df_mount=$(df | grep ^$i | awk '{print $6}')
        df_h=$(df | grep ^$i | awk '{print $5}' | awk -F% '{print $1}')
        df_inode=$(df -i | grep ^$i | awk '{print $5}' | awk -F% '{print $1}')
        if [[ ${df_h} -ge $POLICE ]];then
            status=false
        else
            status=true
        fi
        PostToApi storage_${i}_${df_mount}=${df_h}% ${status}

        if [[ ${df_inode} -ge $POLICE ]];then
            status=false
        else
            status=true
        fi
        PostToApi inode_${i}_${df_mount}=${df_h}% ${status}
    done

}

# 所有的基本信息监控
CheckAll(){
    CheckCpu &
    sleep 0.01

    CheckLoad &
    sleep 0.01

    CheckMem &
    sleep 0.01

    CheckDf &
    wait
}

CheckFalcon(){
    local status
    ps -aux | grep 'falcon' | grep -v 'grep'
    if [[ $? -ne 0 ]];then
        status=false
    else
        status=true
    fi

    PostToApi falcon-agent ${status}
}

CheckNginx(){

    local status
    ps -aux | grep 'nginx:' | grep -v 'orange' | grep -v 'grep'
    if [[ $? -ne 0 ]];then
        status=false
    else
        status=true
    fi

    PostToApi nginx ${status}
}

CheckPhp(){

    local status
    ps -aux | grep 'php-fpm: m' | grep -v 'grep'
    if [[ $? -ne 0 ]];then
        status=false
    else
        status=true
    fi

    PostToApi php_fpm ${status}
}

CheckMysql(){

    local status
    ps -aux | grep mysqld | grep -v 'grep'
    if [[ $? -ne 0 ]];then
        status=false
    else
        status=true
        if [[ $(ps -aux | grep mysqld | grep -v 'grep' | wc -l | awk '{print $1}') -ge 200 ]];then
            status=false
        fi
    fi

    PostToApi mysqld ${status}
}

CheckMysqlSlave(){
    local status=1
    local io_status=1
    local sql_status=0

    timeout 5 strace /work/servers/mysql/bin/mysql -u root -h 127.0.0.1 -pk8U@*hy4icomxz -e 'show slave status\G' &> ${TMPMYSQLLOGFILE}
    cat ${TMPMYSQLLOGFILE}
    grep 'Slave_IO_Running: Yes' ${TMPMYSQLLOGFILE}
    io_status=$?
    grep 'Slave_SQL_Running: Yes' ${TMPMYSQLLOGFILE}
    if [[ $? -ne 0 ]];then
        ps -aux | grep LogicBackupMysql | grep -v grep
        sql_status=$?
        if [[ $sql_status -ne 0 ]];then
            ps -aux | grep mysqldump | grep -v grep
            sql_status=$?
        fi
    fi
    echo $(date) "io = ${io_status}  sql = ${sql_status}" >> ${LOGFILE}
    if [[ ${io_status} -ne 0 ]];then
        if [[ ${sql_status} -ne 0 ]] ;then
            PostToApi slave=io+sql false
        else
            PostToApi slave=io false
        fi
    elif [[ ${sql_status} -ne 0 ]] ;then
        PostToApi slave=sql false
    else
        PostToApi slave true
    fi

}

# 进行监控报警
# 参数1：监控模块
Monitor(){
     case $1 in
         all)
             CheckAll
             ;;
         falcon)
             CheckFalcon
             ;;
         nginx)
             CheckNginx
             ;;
         mysql)
             CheckMysql
             ;;
         slave)
             CheckMysqlSlave
             ;;
         php)
             CheckPhp
             ;;

     esac
}

# 错误日志记录并警报退出
# 参数1：记录异常数
MyError(){
    local error_num
    error_num=$1
    PostToApi "monitor"  true
    echo $(date) ${error_num} >> ${LOGFILE}
    exit ${error_num}
}

# 简单测试环境并完成初始化
Init(){
    mkdir -p ${LOGDIR}
    if [[ ! -w ${LOGDIR} ]];then
        LOGFILE='/tmp/tmp-test.log'
        PostToApi ${LOGDIR}_permission_deny false
    fi
    id ${USER}
    if [[ $? -ne 0 ]];then
        echo you have to init enviroment!
        MyError 2
    fi
    touch ${LOGFILE}
    if [[ ! -w ${LOGFILE} ]];then
        LOGFILE='/tmp/tmp-test.log'
        PostToApi ${LOGFILE}_permission_deny false
    fi
}

# 服务检查前基本检查
MonitorInit(){
    local filenumber

    if [[ s${NATINFO} == "s" ]] || [[ s${IPINFO} == "s" ]];then
         MyError 1
    fi

    PostToApi survive true
    echo true > ${TMPSTATUSLOGFILE}

    filenumber=$(ps aux | grep $0 | grep -v grep | wc -l | awk '{print $1}')
    if [[ ${filenumber} -ge ${POLICE} ]];then
        PostToApi patrol-shell=${filenumber} false
    else
        PostToApi patrol-shell=${filenumber} true
    fi
}

# 处理日志
# 参数1：日志绝对路径    参数2：日志报警大小上限
CleanLogFile(){
    local logfile=$1
    local maxsize=$2
    if [[ -f ${logfile} ]] && [[ $(ls -l ${logfile} | awk '{print $5}') -ge ${maxsize} ]];then
        rm -f ${logfile}.bak
        mv ${logfile} ${logfile}.bak
    fi
}

Main(){
    echo $(date) $@
    local monitor_modle=""
    local short_opts="hn:i:t:faxmsp"
    local long_opts="help,nat:,ip:,time:,all,falcon,nginx,mysql,slave,php"
    local argsw
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
                NATINFO=$2
                shift 2
                ;;
            -t|--time)
                DATE=$(TZ=Asia/Shanghai date  -d "@$2" "+%Y-%m-%d %H:%M")
                shift 2
                ;;
            -i|--ip)
                IPINFO=$2
                shift 2
                ;;
            -a|--all)
                monitor_modle=${monitor_modle}"all "
                shift
                ;;
            -f|--falcon)
                monitor_modle=${monitor_modle}"falcon "
                shift
                ;;
            -x|--nginx)
                monitor_modle=${monitor_modle}"nginx "
                shift
                ;;
            -m|--mysql)
                monitor_modle=${monitor_modle}"mysql "
                shift
                ;;
            -s|--slave)
                monitor_modle=${monitor_modle}"slave "
                shift
                ;;
            -p|--php)
                monitor_modle=${monitor_modle}"php "
                shift
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

    MonitorInit

    # 清理日志
    CleanLogFile ${RETRYLOGFILE} 1048576
    CleanLogFile ${LOGFILE} ${MAXSIZE}
    CleanLogFile /work/logs/openfalcon/patrol-nat-monitor.log ${MAXSIZE}
    CleanLogFile /work/logs/openfalcon/monitor.log $[MAXSIZE*10]
    find ${LOGDIR} -name 'patrol-detail.*.log' -mmin +1000 -exec rm -f {} \;

    # 错开巡查机器上传数据的时间
    sleep $[RANDOM%10].$[RANDOM%999]
    # 获取基本参数后，开始监控模块
    echo "$(date) start to patrol monitor." >> ${LOGFILE}

    for i in ${monitor_modle}
    do
        sleep 0.01
        Monitor $i &
    done

    # 等待监控完毕，解锁
    wait
    if [[ $(cat ${TMPSTATUSLOGFILE}) == 'true' ]];then
        rm -f ${TMPLOGFILE}
    fi
    rm -f ${TMPSTATUSLOGFILE} ${TMPCPULOGFILE} ${TMPMYSQLLOGFILE}
}

# 显示脚本用法
Usage(){
    cat <<EOF

USAGE:$0 [OPTIONS] [work_password]

选择安装模式：
    -h | --help          查看帮助信息
    -n | --nat           指定nat机IP信息
    -i | --ip            指定本机IP信息
    -t | --time          指定当前时间信息
    -f | --falcon        监控openfalcon

EOF
}

Init
Main $@ &> ${TMPLOGFILE} &
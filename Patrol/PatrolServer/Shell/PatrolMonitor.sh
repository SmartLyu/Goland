#!/bin/bash
# ==================
# Description: 巡查使用监控程序
# Created By: 于志远
# Version: 0.1
# Last Modified: 2019-5-10
# ==================

# 全局变量
TIMEOUT='3'
USER=work

STEPLENTH='60'
URL="patrol.ijunhai.com:8686/monitor/collect"
HOST=${HOSTNAME}
LOGDIR="/work/logs/openfalcon"
LOGFILE="${LOGDIR}/patrol-monitor.log"
DATE=$(TZ=Asia/Shanghai date "+%Y-%m-%d %H:%M")
MAXSIZE=10485760
ERRORLOG="${LOGDIR}/error.log"

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

    data="{ \
        \"time\": \"${DATE}\", \
        \"IP\": \"${ip_info}\", \
        \"hostname\": \"${HOST}\", \
        \"info\": \"${tag_name}\", \
        \"status\": ${status} \
        }"
    echo "$(date) ${ip_info}:${tag_name}  ${status}  " >> ${LOGFILE}
    curl -s --connect-timeout ${TIMEOUT} -X POST -d "${data}" $URL &>> ${LOGFILE}
    if [[ $? -ne 0 ]];then
        sleep 1
        echo First $(date) ${data} >> /tmp/retry.log
        curl -s --connect-timeout ${TIMEOUT} -X POST -d "${data}" $URL 2>> ${LOGFILE}
        if [[ $? -ne 0 ]];then
            sleep 1
            echo Second $(date)  ${data} >> /tmp/retry.log
            curl -s --connect-timeout ${TIMEOUT} -X POST -d "${data}" $URL 2>> ${LOGFILE}
        fi
    fi
    echo $(date) Push data ${tag_name} successfully >> ${LOGFILE}
}

CheckCpu(){
    local status
    local cpu
    cpu=$(top -bn 2 |grep Cpu | tail -1 | cut -d "," -f 1 | cut -d ":" -f 2 | awk '{print $1}');

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
        top -bn 2 &>> ${ERRORLOG}
    else
        status=true
    fi

    PostToApi cpu_used=${cpu}% ${status}
}

CheckLoad(){
    local status
    local load
    local cpu_number
    cpu_number=$(cat /proc/cpuinfo| grep "processor"| wc -l)
    if [[ ${cpu_number} == 1 ]] || [[ ${cpu_number} == 2 ]];then
        return
    fi
    load=$(uptime | awk '{printf ("%d\n",$NF/'${cpu_number}'*100)}')
    if [[ load -ge $POLICE ]] || [[ s$load == 's' ]];then
        status=false
        top -bn 2 &>> ${ERRORLOG}
    else
        status=true
    fi

    PostToApi load_used=${load}% ${status}
}

CheckMem(){
    local status
    local mem
    mem=$(free -t | grep Total | awk '{printf ("%d\n",$3/$2*100)}')
    if [[ ${mem} -ge $POLICE ]] || [[ s$mem == 's' ]];then
        status=false
        top -bn 2 &>> ${ERRORLOG}
    else
        status=true
    fi

    PostToApi mem_used=${mem}% ${status}
}

CheckDf(){
    local status
    local df
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
    fi

    PostToApi mysqld ${status}
}

CheckMysqlSlave(){
    local status=1
    local io_status=1
    local sql_status=0
    /work/servers/mysql/bin/mysql -u root -h 127.0.0.1 -pk8U@*hy4icomxz -e 'show slave status\G' | grep 'Slave_IO_Running: Yes'
    io_status=$?
    /work/servers/mysql/bin/mysql -u root -h 127.0.0.1 -pk8U@*hy4icomxz -e 'show slave status\G' | grep 'Slave_SQL_Running: Yes'
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

    if [[ s${NATINFO} == "s" ]] || [[ s${IPINFO} == "s" ]];then
         MyError 1
    fi
    mkdir -p ${LOGDIR}
    if [[ ! -w ${LOGDIR} ]];then
        PostToApi ${LOGDIR}_permission false
        LOGFILE='/tmp/tmp-test.log'
    fi
    id ${USER}
    if [[ $? -ne 0 ]];then
        echo you have to init enviroment!
        MyError 2
    fi
    touch ${LOGFILE}

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
    local monitor_modle=""
    local short_opts="hn:i:faxmsp"
    local long_opts="help,nat:,ip:,all,falcon,nginx,mysql,slave,php"
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

    Init
    PostToApi survive true

    # 清理日志
    CleanLogFile /tmp/retry.log 1048576
    CleanLogFile ${LOGFILE} ${MAXSIZE}
    CleanLogFile /work/logs/openfalcon/patrol-nat-monitor.log ${MAXSIZE}
    CleanLogFile /work/logs/openfalcon/monitor.log $[MAXSIZE*10]
    CleanLogFile ${ERRORLOGE} ${MAXSIZE}

    echo "$(date) prepare to patrol monitor." >> ${LOGFILE}
    # 错开巡查机器上传数据的时间
    sleep $[RANDOM%10].$[RANDOM%999]
    # 获取基本参数后，开始监控模块
    echo "$(date) start to patrol monitor." >> ${LOGFILE}

    for i in ${monitor_modle}
    do
        sleep 0.01
        Monitor $i &
    done
}

# 显示脚本用法
Usage(){
    cat <<EOF

USAGE:$0 [OPTIONS] [work_password]

选择安装模式：
    -h | --help          查看帮助信息
    -n | --nat           指定nat机IP信息
    -i | --ip            指定本机IP信息
    -f | --falcon        监控openfalcon

EOF
}

Main $@ &
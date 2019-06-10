#!/bin/bash
# ==================
# Description: 巡查使用监控程序
# Created By: 于志远
# Version: 0.1
# Last Modified: 2019-5-10
# ==================

# 全局变量
TIMEOUT='10'
USER=work

STEPLENTH='60'
URL="134.175.50.184:8686/monitor/collect"
HOST=${HOSTNAME}
LOGDIR="/work/logs/openfalcon"
LOGFILE="${LOGDIR}/patrol-monitor.log"
DATE=$(date "+%Y-%m-%d %H:%M")

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

    ip_info=${NATINFO}" to "${IPINFO}
    tag_name=$1
    status=$2

    data="{ \
        \"time\": \"${DATE}\", \
        \"IP\": \"${ip_info}\", \
        \"hostname\": \"${HOST}\", \
        \"info\": \"${tag_name}\", \
        \"status\": ${status} \
        }"
    echo "${ip_info}:${tag_name}  ${status}  " >> ${LOGFILE}
    curl -s --connect-timeout ${TIMEOUT} -X POST -d "${data}" $URL
    echo
}

CheckCpu(){
    local status
    local cpu
    cpu=$(top -n 1 |grep Cpu | cut -d "," -f 1 | cut -d ":" -f 2 | awk '{print $2}');
    echo $cpu | grep '%'
    if [[ $? -eq 0 ]];then
        cpu=$(echo $cpu | awk -F% '{print $1}')
    else
        cpu=$(echo $cpu | awk '{print $1*100}')
    fi
    echo $cpu

    if [[ $cpu -ge $POLICE ]] || [[ s$cpu == 's' ]];then
        status=false
    else
        status=true
    fi
    echo '$(date) cpu is '${cpu} >> ${LOGFILE}

    PostToApi cpu=${cpu} ${status}
}

CheckLoad(){
    local status
    local load
    local cpu_number
    cpu_number=$(cat /proc/cpuinfo| grep "processor"| wc -l)
    load=$(uptime | awk '{printf ("%d\n",$NF/'${cpu_number}'*100)}')
    if [[ load -ge $POLICE ]] || [[ s$load == 's' ]];then
        status=false
    else
        status=true
    fi
    echo '$(date) load is '${load} >> ${LOGFILE}

    PostToApi load=${load} ${status}
}

CheckMem(){
    local status
    local mem
    mem=$(free -t | grep Total | awk '{printf ("%d\n",$3/$2*100)}')
    if [[ ${mem} -ge $POLICE ]] || [[ s$mem == 's' ]];then
        status=false
    else
        status=true
    fi
    echo '$(date) mem is '${mem} >> ${LOGFILE}

    PostToApi mem=${mem} ${status}
}

CheckDf(){
    local status
    local df
    df=$(df | grep ^/ | awk '{print $5}' | awk -F% '{print $1}')
    for i in $df
    do
        if [[ ${i} -ge $POLICE ]];then
            status=false
        else
            status=true
        fi
        echo '$(date) df is '${i} >> ${LOGFILE}

        PostToApi RemainingStorage=${i} ${status}
    done

}

# 所有的基本信息监控
CheckAll(){
    CheckCpu
    CheckLoad
    CheckMem
    CheckDf
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
    /work/bin/mysql -u root -h 127.0.0.1 -pk8U@*hy4icomxz -e 'show slave status\G' | grep 'Slave_IO_Running: Yes'
    io_status=$?
    /work/bin/mysql -u root -h 127.0.0.1 -pk8U@*hy4icomxz -e 'show slave status\G' | grep 'Slave_SQL_Running: Yes'
    if [[ $? -ne 0 ]];then
        ps -aux | grep LogicBackupMysql | grep -v grep
        sql_status=$?
    fi
    echo "io = ${io_status}  sql = ${sql_status}" >> ${LOGFILE}
    if [[ ${io_status} -ne 0 ]] || [[ ${sql_status} -ne 0 ]] ;then
        status=false
    else
        status=true
    fi

    PostToApi slave ${status}
}

# 进行监控报警
# 参数1：监控模块
Monitor(){
     case $1 in
         all)
             CheckAll
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
    id ${USER} && mkdir -p ${LOGDIR}
    if [[ $? -ne 0 ]];then
        echo you have to init enviroment!
        MyError 2
    fi
    touch ${LOGFILE}

}

Main(){
    local monitor_modle=""
    local short_opts="hn:i:faxms"
    local long_opts="help,nat:,ip:,all,falcon,nginx,mysql,slave"
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
                CheckFalcon
                return
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

    # 错开巡查机器上传数据的时间
    sleep $[RANDOM%3].$[RANDOM%100]

    # 获取基本参数后，开始监控模块
    Init

    for i in ${monitor_modle}
    do
        Monitor $i
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
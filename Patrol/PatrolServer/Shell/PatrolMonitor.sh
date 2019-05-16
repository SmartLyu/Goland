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
URL="134.175.50.184:8666/monitor/collect"
HOST=${HOSTNAME}
LOGDIR="/work/logs/openfalcon"
LOGFILE="${LOGDIR}/patrol-monitor.log"

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
        \"IP\": \"${ip_info}\", \
        \"hostname\": \"${HOST}\", \
        \"info\": \"${tag_name}\", \
        \"status\": ${status} \
        }"
    echo "${ip_info}:${tag_name}  ${status}  " >> ${LOGFILE}
    curl -s --connect-timeout ${TIMEOUT} -X POST -d "${data}" $URL
    echo
}

CheckFalcon(){
    local status=1
    ps -aux | grep 'falcon' | grep -v 'grep'
    if [[ $? -ne 0 ]];then
        status=false
    else
        status=true
    fi

    PostToApi falcon-agent ${status}
}

# 进行监控报警
# 参数1：监控模块
Monitor(){
     case $1 in
         falcon)
             CheckFalcon
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
    local short_opts="hn:i:f"
    local long_opts="help,nat:,ip:,falcon"
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
            -f|--falcon)
                monitor_modle=${monitor_modle}"falcon "
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
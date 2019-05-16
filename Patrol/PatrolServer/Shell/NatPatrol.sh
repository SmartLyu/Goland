#!/bin/bash
# ==================
# Description: 巡查调用nat使用的程序
# Created By: 于志远
# Version: 0.1
# Last Modified: 2019-5-14
# ==================

# 全局变量
TIMEOUT='10'
USER=work

STEPLENTH='60'
URL="134.175.50.184:8666/monitor/collect"
HOSTSURL="134.175.50.184:8666/monitor/nat"
HOST=${HOSTNAME}
LOGDIR="/work/logs/openfalcon"
LOGFILE="${LOGDIR}/patrol-monitor.log"
DOWNLOAD_URL="134.175.50.184:8666/shell/monitor"

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
        \"IP\": \"${ip_info}\", \
        \"hostname\": \"${HOST}\", \
        \"info\": \"${tag_name}\", \
        \"status\": ${status} \
        }"
    echo "${ip_info}:${tag_name}  ${status}  " >> ${LOGFILE}
    curl -s --connect-timeout ${TIMEOUT} -X POST -d "${data}" $URL
    echo
}

# 发送hosts数据给监控的API
# 参数1：后端服务器IP
PostHostsToApi(){
    local ip_info
    local data

    ip_info=${NATINFO}" to "${1}

    data="{ \
        \"IP\": \"${ip_info}\" \
        }"
    echo "post ${ip_info}  " >> ${LOGFILE}
    curl -s --connect-timeout ${TIMEOUT} -X POST -d "${data}" $HOSTSURL
    echo
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

    if [[ s${NATINFO} == "s" ]] ;then
         MyError 1
    fi
    id ${USER} && mkdir -p ${LOGDIR}
    if [[ $? -ne 0 ]];then
        echo you have to init enviroment!
        MyError 2
    fi
    touch ${LOGFILE}

}

# 遍历hosts进行监控
Monitor(){
    local all_ip
    all_ip=$(awk '{print $1}' /etc/hosts | grep -v ^# |  grep -v ^$ | grep -E ^\(10\)\\..*\|\(172\)\\..*\|\(192\)\\..*)

    for i in ${all_ip}
    do
        PostHostsToApi $i
        ssh $i "wget "${DOWNLOAD_URL}" --timeout 10 -O /tmp/patrol-tmp.sh; \
               /bin/bash /tmp/patrol-tmp.sh --nat "${NATINFO}" --ip "${i}" -f;\
               rm -f /tmp/patrol-tmp.sh" &
    done
}

Main(){
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
                NATINFO=$2
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

Main $@ &
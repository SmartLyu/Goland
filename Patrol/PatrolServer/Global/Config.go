package Global

import (
	"strconv"
	"time"
)

var (
	ApiPost          = "8666"                         // 程序端口
	DataFileDir      = "/work/data/patril/"           // 存放历史监控信息
	DataFileName     = ".monitor.log"                 // 监控后缀名
	LogFileDir       = "/work/logs/patril/"           // 日志存放目录(提前准备好)
	LogFileName      = ".patril.log"                  // 日志后缀名
	AcessLogFileName = ".access.log"                  // 记录日志信息
	MonitorShellFile = "/work/sh/PatrolMonitor.sh"    // 巡查脚本存放位置
	NatFileName      = "patrilNat.file"               // 记录nat机器发送的后端机器
	NatShellFile     = "/work/sh/NatPatrol.sh"        // Nat使用的巡查脚本存放位置
	ListenSig        = make(chan int)                 // 监听后台阻塞信号
	CocoUrl          = "http://10.4.0.4:8666/monitor" // coco的端口
)

// 自动分隔错误日志
func UpdateLog() (string, string) {
	return LogFileDir + time.Now().Format("2006-01") + "/",
		LogFileDir + time.Now().Format("2006-01") + "/" + strconv.Itoa(time.Now().Day()) + LogFileName
}

// 自动分隔错误日志
func UpdateAcessLog() (string, string) {
	return LogFileDir + time.Now().Format("2006-01") + "/",
		LogFileDir + time.Now().Format("2006-01") + "/" + strconv.Itoa(time.Now().Day()) + AcessLogFileName
}

// 自动分隔巡查信息
func UpdateFile() (string, string) {
	return DataFileDir + time.Now().Format("2006-01") + "/",
		DataFileDir + time.Now().Format("2006-01") + "/" + strconv.Itoa(time.Now().Day()) + DataFileName
}
package Global

import (
	"time"
)

var (
	ApiPost                = "8666"                         // 程序端口
	ApiPublicPost          = "8686"                         // 对外程序端口
	DataFileDir            = "/work/data/patril/"           // 存放历史监控信息
	DataFileName           = ".monitor.log"                 // 监控后缀名
	LogFileDir             = "/work/logs/patril/"           // 日志存放目录(提前准备好)
	LogFileName            = ".patril.log"                  // 日志后缀名
	AcessLogFileName       = ".access.log"                  // 记录日志信息
	MonitorShellFile       = "/work/sh/PatrolMonitor.sh"    // 巡查脚本存放位置
	NatShellFile           = "/work/sh/NatPatrol.sh"        // Nat使用的巡查脚本存放位置
	ErrorMap               = NewErrorMapType()              // 存储报警信息至内存
	NatHostsMap            = NewNatHostsMapType()           // 存储Nat机器中子服务器信息至内存
	ErrorMax               = 3                              // 最多报警次数
	MaxSearchLen     int64 = 100                            // 搜索文件最大次数
	MaxReturnLen     int64 = 1000000                        // 查询预估临界值
	ListenSig              = make(chan int)                 // 监听后台阻塞信号
	ListenPublicSig        = make(chan int)                 // 监听后台公共端口阻塞信号
	CocoUrl                = "http://10.4.0.4:8666/monitor" // coco的端口
	IsPolice               = true                           // 是否报警
)

// 自动分隔错误日志
func UpdateLog() (string, string) {
	return LogFileDir + time.Now().Format("2006-01") + "/",
		LogFileDir + time.Now().Format("2006-01/02") + LogFileName
}

// 自动分隔错误日志
func UpdateAcessLog() (string, string) {
	return LogFileDir + time.Now().Format("2006-01") + "/",
		LogFileDir + time.Now().Format("2006-01/02") + AcessLogFileName
}

// 自动分隔巡查信息
func UpdateFile() (string, string) {
	return DataFileDir + time.Now().Format("2006-01/02/15") + "/",
		DataFileDir + time.Now().Format("2006-01/02/15/04") + DataFileName
}

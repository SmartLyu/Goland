package Global

import (
	"time"
)

var (
	ApiPost          = "8888"                                                                   // 程序端口
	GetURL           = "http://139.159.217.194:8666/lyu-data/"                                  // 获取后端存储文件
	DataFileDir      = "/work/apps/lyu-data/"                                                   // 存放生成html文件
	LogFileDir       = "/work/logs/yu-editor/"                                                  // 日志存放目录(提前准备好)
	LogFileName      = ".error.log"                                                            // 日志后缀名
	AcessLogFileName = ".access.log"                                                            // 记录日志信息
	ShellFile        = "/work/apps/lyu-editor/lyuSed.sh"                                                              // html文件自动修正脚本
	ListenSig        = make(chan int)                                                           // 监听后台阻塞信号
	letterRunes      = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890") //随机字符池
	letterLenth      = 12                                                                       //随机字符长度
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

// 自动分隔每日数据信息
func UpdateFile() (string) {
	return DataFileDir + time.Now().Format("2006-01") + "/"
}

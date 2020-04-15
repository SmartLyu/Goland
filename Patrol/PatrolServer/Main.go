package main

import (
	"./Api"
	"./CallCoco"
	"./Global"
	"./Mysql"
	"time"
)

func main() {
	// 初始化日志
	Global.Log()
	// 初始化连接数据库
	Mysql.InitDB()

	// 监听查询端口
	go Api.StartApi(Global.ApiPost)
	// 监听输入输出端口
	go Api.StartPublicApi(Global.ApiPublicPost)

	// 开始巡查所有机器的计划任务
	CallCoco.StartAllCrontab()

	// 开始清空错误队列的计划任务
	CallCoco.CrontabToDelMap()

	// 开始检查存活状态的计划任务
	CallCoco.CrontabToCheckHosts()

	// 开始日志切割的计划任务
	CallCoco.CrontabToCutLog()

	// 刚刚开始任务暂时不启动报警
	Global.IsPolice = false
	go func() {
		time.Sleep(time.Minute * 2)
		if !Global.IsPolice {
			Global.IsPolice = true
		}
	}()

	// 持续监听端口
	<-Global.ListenSig
	<-Global.ListenPublicSig
}

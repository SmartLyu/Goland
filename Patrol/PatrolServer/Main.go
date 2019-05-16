package main

import (
	"./Api"
	"./CallCoco"
	"./Global"
	"./Mysql"
)

func main() {
	Mysql.InitDB()
	go Api.StartApi(Global.ApiPost)
	CallCoco.StartAllCrontab()
	CallCoco.CrontabToCheckHosts()
	<-Global.ListenSig
}

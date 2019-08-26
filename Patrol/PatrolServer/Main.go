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
	go Api.StartPublicApi(Global.ApiPublicPost)
	CallCoco.StartAllCrontab()
	CallCoco.CrontabToDelMap()
	CallCoco.CrontabToCheckHosts()
	<-Global.ListenSig
	<-Global.ListenPublicSig
}

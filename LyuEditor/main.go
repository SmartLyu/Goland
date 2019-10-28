package main

import (
	"./Global"
	"./Api"
)

func main(){
	go Api.StartApi(Global.ApiPost)
	<-Global.ListenSig
}

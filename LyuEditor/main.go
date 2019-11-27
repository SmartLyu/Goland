package main
import (
	"./Api"
	"./Global"
)
func main(){
	go Api.StartApi(Global.ApiPost)
	<-Global.ListenSig
}

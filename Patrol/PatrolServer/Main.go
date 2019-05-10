package main

import (
	"./Api"
	"./Global"
)

func main() {
	Api.StartApi(Global.ApiPost)
	//CallPolice.CallPolice("hello world")
}

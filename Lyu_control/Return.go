package main

import (
	"io/ioutil"
	"os"
	"time"
)

func TryReturn(msgContent MsgContent){

	var text string
	if msgContent.MsgType == "text" {
		text = msgContent.Content
	} else if msgContent.MsgType == "event" {
		text = msgContent.EventKey
	} else {
		text = "help"
	}
	err := myCmd("bash", ShellFile, text)
	if err != nil {
		WriteFile("Error: " + msgContent.Content + "do error")
	}

	justWaite := 1
	// 检查lock文件
	for {
		_, err = os.Stat(LockFile)
		if err != nil {
			break
		}
		time.Sleep(1 * time.Second)
		justWaite ++
		if justWaite > WaiteTime {
			SendMessage(msgContent.Content + " 执行异常！")
			return
		}
	}

	fileData, err := ioutil.ReadFile(ShowFile)
	SendMessage(string(fileData))
}

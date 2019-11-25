package CallPolice

import (
	"time"
)

func CallPolice(message string) {
	id := SecretId
	id.content = "巡查异常：\n   " + time.Now().Format("2006年01月02日 15时04分05秒") + "\n   " + message

	if err := SendWeiXinMessage(id); err != nil {
		WriteErrorLog(err.Error())
	} else {
		WriteInfoLog("Message:" + id.content + " successfully")
	}
}

func CallRestore(message string) {
	id := SecretId
	id.content = "巡查检测到恢复：\n   " + time.Now().Format("2006年01月02日 15时04分05秒") + "\n   " + message

	if err := ForceSendMessage(id); err != nil {
		WriteErrorLog(err.Error())
	} else {
		WriteInfoLog("Message:" + id.content + " successfully")
	}
}

func CallMessage(message string) {
	id := MessageId
	id.content = "巡查系统操作：\n   " + time.Now().Format("2006年01月02日 15时04分05秒") + "\n   " + message

	if err := ForceSendMessage(id); err != nil {
		WriteErrorLog(err.Error())
	} else {
		WriteInfoLog("Message:" + id.content + " successfully")
	}
}

package CallPolice

import (
	"../Global"
	"time"
)

func CallPolice(message string) {
	id := SecretId
	id.content = "巡查异常：" + time.Now().Format("2006年01月02日 15时04分05秒") + "\n" + message

	if err := SendWeiXinMessage(id); err != nil {
		(err.Error())
	} else {
		Global.InfoLog.Println("Message:" + id.content + " successfully")
	}
}

func CallRestore(message string) {
	id := SecretId
	id.content = "巡查检测到恢复：" + time.Now().Format("2006年01月02日 15时04分05秒") + "\n" + message

	if err := SendWeiXinMessage(id); err != nil {
		Global.ErrorLog.Println(err.Error())
	} else {
		Global.InfoLog.Println("Message:" + id.content + " successfully")
	}
}

func CallMessage(message string) {
	id := MessageId
	id.content = "巡查系统操作：" + time.Now().Format("2006年01月02日 15时04分05秒") + "\n" + message

	if err := ForceSendMessage(id); err != nil {
		Global.ErrorLog.Println(err.Error())
	} else {
		Global.InfoLog.Println("Message:" + id.content + " successfully")
	}
}

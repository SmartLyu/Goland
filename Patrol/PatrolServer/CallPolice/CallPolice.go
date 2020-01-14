package CallPolice

import (
	"../Global"
	"time"
)

func CallPolice(message ...string) {
	Global.PoliceLog.Println(message)
	id := SecretId
	id.content = "巡查发现异常：" + time.Now().Format("2006年01月02日 15时04分05秒")
	for _, i := range message {
		id.content = id.content + "\n  " + i
	}

	if err := SendWeiXinMessage(id); err != nil {
		(err.Error())
	} else {
		Global.InfoLog.Println("Message:" + id.content + " successfully")
	}
}

func CallRestore(message ...string) {
	id := SecretId
	id.content = "巡查检测恢复：" + time.Now().Format("2006年01月02日 15时04分05秒")
	for _, i := range message {
		id.content = id.content + "\n  " + i
	}

	if err := SendWeiXinMessage(id); err != nil {
		Global.ErrorLog.Println(err.Error())
	} else {
		Global.InfoLog.Println("Message:" + id.content + " successfully")
	}
}

func CallMessage(message ...string) {
	Global.PoliceLog.Println(message)
	id := MessageId
	id.content = "巡查系统异常：" + time.Now().Format("2006年01月02日 15时04分05秒")
	for _, i := range message {
		id.content = id.content + "\n  " + i
	}

	if err := ForceSendMessage(id); err != nil {
		Global.ErrorLog.Println(err.Error())
	} else {
		Global.InfoLog.Println("Message:" + id.content + " successfully")
	}
}

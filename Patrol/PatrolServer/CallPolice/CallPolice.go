package CallPolice

import (
	"../Global"
	"time"
)

func CallPolice(message ...string) {
	Global.PoliceLog.Println(Global.IsPolice, message)
	policeMessage := "巡查发现异常：" + time.Now().Format("2006年01月02日 15时04分05秒")
	for _, i := range message {
		policeMessage = policeMessage + "\n" + i
	}

	if err := SendPoliceMessage(policeMessage, "失败"); err != nil {
		Global.ErrorLog.Println(err.Error())
		CallMessage("报警请求异常：" + err.Error())
	} else {
		Global.InfoLog.Println("Message:" + policeMessage + " successfully")
	}
}

func CallRestore(message ...string) {
	Global.PoliceLog.Println(Global.IsPolice, message)
	policeMessage := "巡查发现恢复：" + time.Now().Format("2006年01月02日 15时04分05秒")
	for _, i := range message {
		policeMessage = policeMessage + "\n" + i
	}

	if err := SendPoliceMessage(policeMessage, "成功"); err != nil {
		Global.ErrorLog.Println(err.Error())
	} else {
		Global.InfoLog.Println("Message:" + policeMessage + " successfully")
	}
}

func CallMessage(message ...string) {
	Global.PoliceLog.Println(Global.IsPolice, message)
	id := MessageId
	id.content = "巡查系统日志：" + time.Now().Format("2006年01月02日 15时04分05秒")
	for _, i := range message {
		id.content = id.content + "\n" + i
	}

	if err := SendWeiXinMessage(id); err != nil {
		Global.ErrorLog.Println(err.Error())
	} else {
		Global.InfoLog.Println("Message:" + id.content + " successfully")
	}
}

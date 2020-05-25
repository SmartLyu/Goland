package CallPolice

import (
	"../Global"
	"time"
)

func CallPolice(id DingdingID, message ...string) {
	Global.PoliceLog.Println(Global.IsPolice, id.dingdingJson.At.AtMobiles, message)
	// 钉钉通知所有人
	id.memssage = "巡查发现异常：" + time.Now().Format("2006年01月02日 15时04分05秒")
	for _, i := range message {
		id.memssage = id.memssage + "\n" + i
	}
	SendPoliceMessage(id)
}

func CallRestore(id DingdingID, message ...string) {
	id.dingdingJson.At.AtMobiles = nil
	Global.PoliceLog.Println(Global.IsPolice, id.dingdingJson.At.AtMobiles, message)
	// 钉钉通知所有人
	id.memssage = "巡查发现恢复：" + time.Now().Format("2006年01月02日 15时04分05秒")
	for _, i := range message {
		id.memssage = id.memssage + "\n" + i
	}
	SendPoliceMessage(id)
}

func CallMessage(message ...string) {
	Global.PoliceLog.Println(Global.IsPolice, message)
	// 钉钉通知负责人
	id := messageDingdingID
	id.memssage = "巡查系统日志：" + time.Now().Format("2006年01月02日 15时04分05秒")
	for _, i := range message {
		id.memssage = id.memssage + "\n" + i
	}
	SendPoliceMessage(id)
}

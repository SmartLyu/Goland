package CallPolice

import (
	"../File"
)

func CallPolice(message string) {
	id := SecretId
	id.content = "巡查异常：\n   "+message

	if err := SendWeiXinMessage(id); err != nil {
		File.WriteErrorLog(err.Error())
	} else {
		File.WriteInfoLog("Message:" + id.content + "to " + string(id.agentid) + "successfully")
	}
}

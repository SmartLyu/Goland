package CallPolice

func CallPolice(message string) {
	id := SecretId
	id.content = "巡查异常：\n   "+message

	if err := SendWeiXinMessage(id); err != nil {
		WriteErrorLog(err.Error())
	} else {
		WriteInfoLog("Message:" + id.content + "to " + string(id.agentid) + "successfully")
	}
}
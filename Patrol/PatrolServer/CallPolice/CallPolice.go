package CallPolice

func CallPolice(message string) {
	id := SecretId
	id.content = "巡查异常：\n   " + message

	if err := SendWeiXinMessage(id); err != nil {
		WriteErrorLog(err.Error())
	} else {
		WriteInfoLog("Message:" + id.content + " successfully")
	}
}

func CallMessage(message string) {
	id := MessageId
	id.content = "巡查系统操作：\n   " + message

	if err := SendWeiXinMessage(id); err != nil {
		WriteErrorLog(err.Error())
	} else {
		WriteInfoLog("Message:" + id.content + " successfully")
	}
}
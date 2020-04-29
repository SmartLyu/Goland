package CallPolice

import "../Global"

func SendPoliceMessage(dingdingID DingdingID) {
	if !Global.IsPolice {
		return
	}

	// 微信通知负责人
	weixinId := MessageId
	weixinId.content = dingdingID.memssage
	if err := SendWeiXinMessage(weixinId); err != nil {
		Global.ErrorLog.Println("send weixin error: ", err.Error())
	} else {
		Global.InfoLog.Println("WeiXin Message:" + weixinId.content + " successfully")
	}

	// 钉钉通知指定人员
	if err := SendDingdingMessage(dingdingID); err != nil {
		Global.ErrorLog.Println("send dingding error: ", err.Error())
	} else {
		Global.InfoLog.Println("Dingding Message:" + dingdingID.memssage + " successfully")
	}
}

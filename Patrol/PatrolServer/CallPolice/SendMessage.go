package CallPolice

import (
	"../Global"
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

const (
	//发送消息使用导的url
	sendurl = `https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=`
	//获取token使用导的url
	getToken = `https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=`
)

type corpText struct {
	corpid     string
	corpsecret string
	touser     string
	toparty    string
	totag      string
	agentid    int
	content    string
	safe       string
}

type accessToken struct {
	Access_token string `json:"access_token"`
	Expires_in   int    `json:"expires_in"`
}

var requestError = errors.New("request error,check url or network")

//定义一个简单的文本消息格式
type sendMsg struct {
	Touser  string            `json:"touser"`
	Toparty string            `json:"toparty"`
	Totag   string            `json:"totag"`
	Msgtype string            `json:"msgtype"`
	Agentid int               `json:"agentid"`
	Text    map[string]string `json:"text"`
	Safe    int               `json:"safe"`
}

type sendMsgError struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

func ForceSendMessage(id corpText) error{
	corpid := id.corpid
	corpsecret := id.corpsecret

	var m = sendMsg{
		Touser:  id.touser,
		Toparty: id.toparty,
		Totag:   id.totag,
		Msgtype: "text",
		Agentid: id.agentid,
		Text:    map[string]string{"content": id.content},
	}

	token, err := GetToken(corpid, corpsecret)
	if err != nil {
		return err
	}

	buf, err := json.Marshal(m)
	if err != nil {
		return err
	}

	err = SendMsg(token.Access_token, buf)
	if err != nil {
		return err
	}

	return nil
}

func SendWeiXinMessage(id corpText) error {
	if ! Global.IsPolice {
		return nil
	}
	return ForceSendMessage(id)
}

//发送消息.msgbody 必须是 API支持的类型
func SendMsg(TimeAccessToken string, msgbody []byte) error {
	body := bytes.NewBuffer(msgbody)
	resp, err := http.Post(sendurl+TimeAccessToken, "application/json", body)
	if resp.StatusCode != 200 {
		return requestError
	}
	buf, _ := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()

	var e sendMsgError
	err = json.Unmarshal(buf, &e)

	if err != nil {
		return err
	}

	if e.Errcode != 0 && e.Errmsg != "ok" {
		return errors.New(string(buf))
	}
	return nil
}

//通过corpid 和 corpsecret 获取token
func GetToken(corpid, corpsecret string) (at accessToken, err error) {

	resp, err := http.Get(getToken + corpid + "&corpsecret=" + corpsecret)
	if err != nil {
		return
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		err = requestError
		return
	}

	buf, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(buf, &at)
	if err != nil {
		return
	}
	if at.Access_token == "" {
		err = errors.New("corpid or corpsecret is error")
	}
	return
}

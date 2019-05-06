package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	//发送消息使用导的url
	sendurl = `https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=`
	//获取token使用导的url
	getToken = `https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=`
)

type accessToken struct {
	TimeAccessToken string `json:"accessToken"`
	ExpiresIn       int    `json:"expires_in"`
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

func main() {
	touser := "1"
	agentid := 0
	content := "Test Api"
	corpid := "wwd0b64b830dd50312"
	corpsecret := "Mebw1BA1cZBY17FrGJcRJniaGWt_LCQe_DGyPLDznRI"

	var m = sendMsg{Touser: touser, Msgtype: "text", Agentid: agentid, Text: map[string]string{"content": content}}

	///-p "wx2468f5838693e123" -s "JbjkM1jYq8g3GaHjOTgj27y4n4_7Dsv4FV94I5BMRSrBsm_aTsMUVJMhGu_DFGDSF"
	token, err := GetToken(corpid, corpsecret)
	if err != nil {
		println(err.Error())
		return
	}

	buf, err := json.Marshal(m)
	if err != nil {
		return
	}

	err = SendMsg(token.TimeAccessToken, buf)
	if err != nil {
		println(err.Error())
	}
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
	if at.TimeAccessToken == "" {
		err = errors.New("corpid or corpsecret error.")
	}
	return
}

func Parse(jsonpath string) ([]byte, error) {
	var zs = []byte("//")
	File, err := os.Open(jsonpath)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = File.Close()
	}()

	var buf []byte
	b := bufio.NewReader(File)
	for {
		line, _, err := b.ReadLine()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, err
		}
		line = bytes.TrimSpace(line)
		if len(line) <= 0 {
			continue
		}
		index := bytes.Index(line, zs)
		if index == 0 {
			continue
		}
		if index > 0 {
			line = line[:index]
		}
		buf = append(buf, line...)
	}
	return buf, nil
}

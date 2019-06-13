package main

import (
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"weworkapi_golang-master/wxbizmsgcrypt"
)

type MsgContent struct {
	ToUsername   string `xml:"ToUserName"`
	FromUsername string `xml:"FromUserName"`
	CreateTime   uint32 `xml:"CreateTime"`
	MsgType      string `xml:"MsgType"`
	Content      string `xml:"Content"`
	MsgId        string `xml:"MsgId"`
	AgentID      uint32 `xml:"AgentID"`
	EventKey     string `xml:"EventKey"`
}

// 用于微信回调测试
func Test(w http.ResponseWriter, r *http.Request) {
	//解析参数，默认是不会解析的
	err := r.ParseForm()
	if err != nil {
		WriteFile("Error: Recv:" + r.RemoteAddr)
	}

	// 读取用户输入的参数信息
	verifyMsgsig := r.Form.Get("msg_signature")
	verifyTimestamp := r.Form.Get("timestamp")
	verifyNonce := r.Form.Get("nonce")
	verifyEchostr := r.Form.Get("echostr")

	wxcpt := wxbizmsgcrypt.NewWXBizMsgCrypt(token, encodingAeskey, receiverId, wxbizmsgcrypt.XmlType)

	echostr, cryptErr := wxcpt.VerifyURL(verifyMsgsig, verifyTimestamp, verifyNonce, verifyEchostr)
	if nil != cryptErr {
		WriteFile("Error: verifyurl fail")
	}
	WriteFile("verifyurl success echostr" + string(echostr))
	_, _ = w.Write(echostr)

	// 返回信息
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
}

// 接受消息
func Control(w http.ResponseWriter, r *http.Request) {

	var msgContent MsgContent

	//解析参数，默认是不会解析的
	err := r.ParseForm()
	if err != nil {
		WriteFile("Error: Recv:" + r.RemoteAddr)
	}

	// 读取用户输入的参数信息
	verifyMsgsig := r.Form.Get("msg_signature")
	verifyTimestamp := r.Form.Get("timestamp")
	verifyNonce := r.Form.Get("nonce")

	wxcpt := wxbizmsgcrypt.NewWXBizMsgCrypt(token, encodingAeskey, receiverId, wxbizmsgcrypt.XmlType)

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		WriteFile("Error: " + err.Error())
	}
	if err := r.Body.Close(); err != nil {
		WriteFile("Error: " + err.Error())
	}

	// 获取用户输入内容
	msg, cryptErr := wxcpt.DecryptMsg(verifyMsgsig, verifyTimestamp, verifyNonce, []byte(body))
	if nil != cryptErr {
		WriteFile("Error: verifyurl fail")
	}

	err = xml.Unmarshal(msg, &msgContent)
	if nil != err {
		WriteFile("Error: Unmarshal fail")
	} else {
		// 接收到请求后 异步处理数据
		go TryReturn(msgContent)
	}

	//fmt.Println("body: ", string(body))
	//fmt.Println("msg: ", string(msg))
	//fmt.Println("CreateTime: ", msgContent.CreateTime)
	//fmt.Println("Agentid: ", msgContent.AgentID)
	//fmt.Println("MsgType: ", msgContent.MsgType)
	//fmt.Println("Content: ", msgContent.Content)
	//fmt.Println("Msgid: ", msgContent.MsgId)
	//fmt.Println("FromUsername: ", msgContent.FromUsername)
	//fmt.Println("ToUsername: ", msgContent.ToUsername)
	//fmt.Println("EventKey: ", msgContent.EventKey)

	// 返回信息
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func ReturnNatShell(w http.ResponseWriter, r *http.Request) {
	//解析参数，默认是不会解析的
	err := r.ParseForm()
	if err != nil {
		WriteFile("Error: Recv:" + r.RemoteAddr)
	}

	des := NatShellFile
	desStat, err := os.Stat(des)
	if err != nil {
		WriteFile("Error: File Not Exit " + des)
		http.NotFoundHandler().ServeHTTP(w, r)
	} else if (desStat.IsDir()) {
		WriteFile("Error: File Is Dir" + des)
		http.NotFoundHandler().ServeHTTP(w, r)
	} else {
		data, err := ioutil.ReadFile(des)
		if err != nil {
			WriteFile("Error: Read File Err: " + err.Error())
		} else {
			WriteFile("Send File:" + des)
			_, err = w.Write([]byte(data))
			if err != nil {
				WriteFile("Error: http writed Err " + err.Error())
			}
		}
	}
}
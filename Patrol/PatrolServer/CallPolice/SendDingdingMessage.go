package CallPolice

import (
	"../Global"
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// 钉钉body请求的json结构
type DingdingJson struct {
	Msgtype string            `json:"msgtype"`
	Text    DingdingContent   `json:"text"`
	At      DingdingAtMobiles `json:"at"`
}

type DingdingContent struct {
	Content string `json:"content"`
}

type DingdingAtMobiles struct {
	AtMobiles []int `json:"atMobiles"`
}

type sendMsgError struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

// 存储钉钉机器的数据类型
type DingdingID struct {
	memssage     string
	secret       string
	accessToken  string
	dingdingJson DingdingJson
}

func SendDingdingMessage(id DingdingID) error {
	if !Global.IsPolice {
		return nil
	}

	// 获取sign认证
	value := url.Values{}
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	stringToSign := fmt.Sprintf("%d\n%s", timestamp, id.secret)
	dingdingHmac := hmac.New(sha256.New, []byte(id.secret))
	dingdingHmac.Write([]byte(stringToSign))
	signedByte := dingdingHmac.Sum(nil)
	signedBase64 := base64.StdEncoding.EncodeToString(signedByte)
	if signedBase64 == "" {
		Global.ErrorLog.Println("Get Dingding signedBase64 is nil")
	}
	value.Set("access_token", id.accessToken)
	value.Set("timestamp", fmt.Sprintf("%d", timestamp))
	value.Set("sign", signedBase64)

	// 生成body内容
	dj := id.dingdingJson
	dj.Text.Content = id.memssage
	msgbody, err := json.Marshal(dj)
	if err != nil {
		return err
	}
	body := bytes.NewBuffer(msgbody)

	// 发出请求
	request, err := http.NewRequest(http.MethodPost, policeDingdingUrl, body)
	if err != nil {
		return err
	}
	request.URL.RawQuery = value.Encode()
	request.Header.Add("Content-Type", "application/json;charset=utf-8")
	resp, err := (&http.Client{}).Do(request)
	if err != nil {
		return err
	}

	// 判断返回是否符合要求
	if resp.StatusCode != 200 {
		return requestError
	}
	buf, _ := ioutil.ReadAll(resp.Body)
	if err := resp.Body.Close(); err != nil {
		Global.ErrorLog.Println(err.Error())
	}

	// 解析返回json，判断请求是否成功
	var e sendMsgError
	err = json.Unmarshal(buf, &e)
	if err != nil {
		return err
	}

	if e.Errcode != 0 && e.Errmsg != "ok" {
		return errors.New(string(buf))
	}

	Global.InfoLog.Println(body, string(buf))
	return nil
}

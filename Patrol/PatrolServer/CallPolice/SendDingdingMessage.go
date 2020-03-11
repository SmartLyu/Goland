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

func SendPoliceMessage(memssage string, status string) error {
	if !Global.IsPolice {
		return nil
	}

	/*
		id := SecretId
		id.content = memssage
		if err := SendWeiXinMessage(id); err != nil {
			Global.ErrorLog.Println(err.Error())
		} else {
			Global.InfoLog.Println("Message:" + id.content + " successfully")
		}
	*/

	// 获取sign认证
	value := url.Values{}
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	stringToSign := fmt.Sprintf("%d\n%s", timestamp, policeDingdingSecret)
	dingdingHmac := hmac.New(sha256.New, []byte(policeDingdingSecret))
	dingdingHmac.Write([]byte(stringToSign))
	signedByte := dingdingHmac.Sum(nil)
	signedBase64 := base64.StdEncoding.EncodeToString(signedByte)
	if signedBase64 == "" {
		Global.ErrorLog.Println("Get Dingding signedBase64 is nil")
	}
	value.Set("access_token", policeDingdingAccessToken)
	value.Set("timestamp", fmt.Sprintf("%d", timestamp))
	value.Set("sign", signedBase64)

	// 生成body内容
	dj := dingdingJson
	dj.Text.Content = memssage
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

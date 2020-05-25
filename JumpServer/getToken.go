package main

import (
	"encoding/json"
	"fmt"
)

func ReturnToken() (string, error) {
	var getTokenJson = PostTokenJson{
		Username: username,
		Password: password,
	}
	jsonstr, err := json.Marshal(&getTokenJson)
	if err != nil {
		ErrorLog.Println(err)
		return "", err
	}
	body, err := httpPostJson(url+"api/v1/authentication/auth/", jsonstr)
	if err != nil {
		ErrorLog.Println(err)
		return "", err
	}

	// 断言
	var v interface{}
	if err := json.Unmarshal([]byte(body), &v); err != nil {
		ErrorLog.Println(err)
		return "", err
	}
	data := v.(map[string]interface{})
	return fmt.Sprintf("%v", data["token"]), nil
}

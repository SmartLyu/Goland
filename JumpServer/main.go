package main

import (
	"fmt"
	"os"
)

func main() {
	// 启动日志
	startLog()

	// 获取token
	token, err := ReturnToken()
	if err != nil || token == "" {
		ErrorLog.Println(err, "token may is nil")
		os.Exit(-1)
	}

	data := make(map[string]string)
	data["Authorization"] = "Bearer " + token
	fmt.Println(data)
	body, err := httpGetJson(url+"api/v1/users/users/", data)
	fmt.Print(body)
}

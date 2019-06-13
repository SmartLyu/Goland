package main

import "net/http"

// 启动服务api端口号
func StartApi(port string) {
	router := NewRouter()
	WriteFile("Error: " + http.ListenAndServeTLS(":"+port, CrtFile, KeyFile, router).Error())
	ListenSig <- 0
}

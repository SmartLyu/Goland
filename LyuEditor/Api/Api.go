package Api

import (
	"../File"
	"../Global"
	"net/http"
)

// 启动服务api端口号
func StartApi(port string) {
	router := NewRouter()
	File.WriteErrorLog(http.ListenAndServe(":"+port, router).Error())
	Global.ListenSig <- 0
}


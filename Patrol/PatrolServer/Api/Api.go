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

// 启动公共api接口，用于接受巡查结果
func StartPublicApi(public_port string) {
	router_public := NewPublicRouter()
	File.WriteErrorLog(http.ListenAndServe(":"+public_port, router_public).Error())
	Global.ListenPublicSig <- 0
}

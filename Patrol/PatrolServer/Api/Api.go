package Api

import (
	"../File"
	"../Global"
	"net/http"
)

func StartApi(port string) {
	router := NewRouter()
	File.WriteErrorLog(http.ListenAndServe(":"+port, router).Error())
	Global.ListenSig <- 0
}

package CallCoco

import (
	"../File"
	"../Mysql"
	"strconv"
)

func CallCoco(ip string, port string) {
	err := httpPostJson(ip, port)
	if err != nil {
		File.WriteErrorLog("Post coco to nat " + ip + " error " + err.Error())
		return
	}
	File.WriteInfoLog("call coco to connect " + ip)
}

func CallAllNatMonitor() {
	for _, i := range Mysql.SelectAllNatTable() {
		CallCoco(i.IP, strconv.Itoa(i.Port))
	}
}

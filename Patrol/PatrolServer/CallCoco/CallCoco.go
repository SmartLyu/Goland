package CallCoco

import (
	"../File"
	"../Mysql"
	"strconv"
)

func CallCoco(hostname string, ip string, port string) {
	err := httpPostJson(ip, port)
	if err != nil {
		File.WriteErrorLog("请求 coco 连接 nat： " + hostname + "-" + ip + " 失败！\n  具体异常为：error " + err.Error())
		return
	}
	File.WriteInfoLog("call coco to connect " + ip)
}

func CallAllNatMonitor() {
	for _, i := range Mysql.SelectAllNatTable() {
		CallCoco(i.HostName, i.IP, strconv.Itoa(i.Port))
	}
}

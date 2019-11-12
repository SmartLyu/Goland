package CallCoco

import (
	"../File"
	"../Global"
	"../Mysql"
	"strconv"
)

func CallCoco(hostname string, ip string, port string) {
	err := httpPostJson(ip, port)
	if err == nil {
		File.WriteInfoLog("call coco to connect " + ip)
		return
	} else {
		err := httpPostJson(ip, port)
		if err != nil {
			_, file := Global.UpdateLog()
			_, errf := File.FindWorkInFile(file, "请求 coco 连接 nat", hostname+"-"+ip, " 失败！")
			if errf == nil {
				File.WriteErrorLog("Error\t" + "请求 coco 连接 nat： " + hostname +
					"-" + ip + " 失败！具体异常为：error " + err.Error())
			}

			File.WriteInfoLog("Warning\t" + "请求 coco 连接 nat： " + hostname +
				"-" + ip + " 失败一次！具体异常为：error " + err.Error())
		}
	}
}

func CallAllNatMonitor() {
	for _, i := range Mysql.SelectAllNatTable() {
		CallCoco(i.HostName, i.IP, strconv.Itoa(i.Port))
	}
}

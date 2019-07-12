package CallCoco

import (
	"../CallPolice"
	"../File"
	"../Global"
	"../Mysql"
	"fmt"
	"strconv"
	"time"
)

func CallCoco(hostname string, ip string, port string) {
	err := httpPostJson(ip, port)
	if err != nil {
		dir, file := Global.UpdateLog()

		_, errf := File.FindWorkInFile(file, time.Now().Add(-time.Minute * 1).Format("2006-01-02 15:04"),
			"请求 coco 连接 nat： "+hostname+"-"+ip+" 失败！")
		if errf == nil {
			fmt.Println("police")
			CallPolice.CallPolice("巡查服务器异常：Error \t" + "请求 coco 连接 nat： " + hostname + "-" + ip +
				" 失败！\n  具体异常为：error " + err.Error())
		}
		fmt.Println(err)

		File.WriteLog("Error\t"+"请求 coco 连接 nat： "+hostname+"-"+ip+" 失败！\n  具体异常为：error "+err.Error(), dir, file)
		return
	}
	File.WriteInfoLog("call coco to connect " + ip)
}

func CallAllNatMonitor() {
	for _, i := range Mysql.SelectAllNatTable() {
		CallCoco(i.HostName, i.IP, strconv.Itoa(i.Port))
	}
}

package CallPolice

import (
	"../Global"
	"strconv"
	"strings"
)

func Judge(monitorjson Global.MonitorJson) {
	mapkey := monitorjson.Hostname + "-" + strings.Split(monitorjson.Info, "=")[0]
	if ! monitorjson.Status {
		if  _, isError := Global.ErrorMap[mapkey]; isError && Global.ErrorMap[mapkey] <= Global.ErrorMax{
			CallPolice(monitorjson.Hostname + " 的 " + monitorjson.Info + " 异常 \n具体服务器信息：\n   " + monitorjson.IP)
		}
		Global.ErrorMap[mapkey] ++
	} else {
		if _, isError := Global.ErrorMap[mapkey]; isError && Global.ErrorMap[mapkey] > 2{
			CallRestore(monitorjson.Hostname + " 的 " + monitorjson.Info + " 状态已经恢复")
			delete(Global.ErrorMap, mapkey)
		}
	}
	WriteInfoLog("获取到 " + monitorjson.IP + "的 " + monitorjson.Hostname +
		" 的数据：" + strconv.FormatBool(monitorjson.Status))
}
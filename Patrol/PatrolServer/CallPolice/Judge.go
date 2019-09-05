package CallPolice

import (
	"../Global"
	"strconv"
	"strings"
)

func Judge(monitorjson Global.MonitorJson) {
	mapkey := monitorjson.Hostname + "-" + strings.Split(monitorjson.Info, "=")[0]
	if ! monitorjson.Status {
		if Global.ErrorMap.Exist(mapkey) && Global.ErrorMap.Get(mapkey) <= Global.ErrorMax {
			CallPolice(monitorjson.Hostname + " 的 " + monitorjson.Info + " 异常 \n具体服务器信息：\n   " + monitorjson.IP)
		}
		Global.ErrorMap.Add(mapkey, 1)
	} else {
		if Global.ErrorMap.Exist(mapkey) {
			if Global.ErrorMap.Get(mapkey) >= 2 {
				CallRestore(monitorjson.Hostname + " 的 " + monitorjson.Info + " 状态已经恢复")
			}
			Global.ErrorMap.Delete(mapkey)
		}
	}
	WriteInfoLog("获取到 " + monitorjson.IP + "的 " + monitorjson.Hostname +
		" 的数据：" + strconv.FormatBool(monitorjson.Status))
}

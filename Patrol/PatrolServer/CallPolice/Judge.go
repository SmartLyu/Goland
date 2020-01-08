package CallPolice

import (
	"../Global"
	"strconv"
	"strings"
)

func Judge(monitorjson Global.MonitorJson) {
	mapkey := monitorjson.IP + ":-}" + strings.Split(monitorjson.Info, "=")[0]
	if !monitorjson.Status {
		if Global.ErrorMap.Exist(mapkey) &&
			Global.ErrorMap.Get(mapkey) <= Global.ErrorMax {
			if monitorjson.Hostname == "PatrolMessage" {
				CallMessage(monitorjson.Hostname+" 的 "+monitorjson.Info+" 异常 ", "具体服务器信息："+
					monitorjson.IP, "   异常发生时间："+monitorjson.Time)
			} else {
				CallPolice(monitorjson.Hostname+" 的 "+monitorjson.Info+" 异常 ", "具体服务器信息："+
					monitorjson.IP, "   异常发生时间："+monitorjson.Time)
			}

		}
		Global.ErrorMap.Add(mapkey, 1)
	} else {
		if Global.ErrorMap.Exist(mapkey) {
			if Global.ErrorMap.Get(mapkey) >= 2 {
				if monitorjson.Hostname == "PatrolMessage" {
					CallMessage(monitorjson.Hostname+" 的 "+monitorjson.Info+" 状态已经恢复",
						"恢复发生时间："+monitorjson.Time)
				} else {
					CallRestore(monitorjson.Hostname+" 的 "+monitorjson.Info+" 状态已经恢复",
						"恢复发生时间："+monitorjson.Time)
				}
			}
			Global.ErrorMap.Delete(mapkey)
		}
	}
	Global.InfoLog.Println(monitorjson.Time + " 获取到 " + monitorjson.IP + " 的 " + monitorjson.Hostname +
		"\t的数据：" + monitorjson.Info + " 的状态为：" + strconv.FormatBool(monitorjson.Status))
}

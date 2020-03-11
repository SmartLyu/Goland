package CallPolice

import (
	"../Global"
	"strings"
)

func Judge(monitorjson Global.MonitorJson) {
	mapkey := monitorjson.IP + ":-}" + strings.Split(monitorjson.Info, "=")[0]
	if !monitorjson.Status {
		// 判断异常状态是否报警
		Global.PoliceLock.Lock()
		defer Global.PoliceLock.Unlock()

		if Global.ErrorMap.Exist(mapkey) {
			// 报警发给所有负责人
			if Global.ErrorMap.Get(mapkey) <= Global.ErrorMax {
				if monitorjson.Hostname == "PatrolMessage" {
					CallMessage(monitorjson.Hostname+" 的 "+monitorjson.Info+" 异常 ", "具体服务器信息："+
						monitorjson.IP, "异常发生时间："+monitorjson.Time)
				} else {
					CallPolice(monitorjson.Hostname+" 的 "+monitorjson.Info+" 异常 ", "具体服务器信息："+
						monitorjson.IP, "异常发生时间："+monitorjson.Time)
				}
			}
		}
		Global.ErrorMap.Add(mapkey, 1)
	} else {
		// 判断是否是异常恢复
		if Global.ErrorMap.Exist(mapkey) {
			Global.PoliceLock.Lock()
			defer Global.PoliceLock.Unlock()
			if Global.ErrorMap.Exist(mapkey) && Global.ErrorMap.Get(mapkey) >= 2 {
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
	Global.InfoLog.Println(" 获取到 " + monitorjson.Time + " 的 " + monitorjson.IP + " 的数据：" + monitorjson.Info)
}

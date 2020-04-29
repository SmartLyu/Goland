package CallPolice

import (
	"../Global"
	"strings"
)

func Judge(monitorjson Global.MonitorJson) {
	id := policeDingdingID
	mapkey := monitorjson.IP + ":-}" + strings.Split(monitorjson.Info, "=")[0]
	if !monitorjson.Status {
		// 判断异常状态是否报警
		Global.PoliceLock.Lock()
		defer Global.PoliceLock.Unlock()

		// 判断是否符合ignore条件，如果符合，则放弃报警监控
		if CheckIgnoreTime(monitorjson) {
			Global.InfoLog.Println(" 获取到 " + monitorjson.Time + " 的 " + monitorjson.IP + " 的数据：" + monitorjson.Info)
			Global.PoliceLog.Println("ignore this error: ", monitorjson.Hostname+" 的 "+monitorjson.Info+" 异常 ", "具体服务器信息："+
				monitorjson.IP, "异常发生时间："+monitorjson.Time)
			return
		}

		// 判断是否已存在异常，避免误报
		if Global.ErrorMap.Exist(mapkey) {
			// 报警发给所有负责人
			if Global.ErrorMap.Get(mapkey) <= Global.ErrorMax {
				if monitorjson.Hostname == Global.PatrolMessageString {
					CallMessage(monitorjson.Hostname+" 的 "+monitorjson.Info+" 异常 ", "具体服务器信息："+
						monitorjson.IP, "异常发生时间："+monitorjson.Time)
				} else {
					id.dingdingJson.At.AtMobiles = ReadDingdingAtFile(monitorjson.Hostname)
					CallPolice(id, monitorjson.Hostname+" 的 "+monitorjson.Info+" 异常 ", "具体服务器信息："+
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
				if monitorjson.Hostname == Global.PatrolMessageString {
					CallMessage(monitorjson.Hostname+" 的 "+monitorjson.Info+" 状态已经恢复",
						"恢复发生时间："+monitorjson.Time)
				} else {
					id.dingdingJson.At.AtMobiles = ReadDingdingAtFile(monitorjson.Hostname)
					CallRestore(id, monitorjson.Hostname+" 的 "+monitorjson.Info+" 状态已经恢复",
						"恢复发生时间："+monitorjson.Time)
				}
			}
			Global.ErrorMap.Delete(mapkey)
		}
	}
	Global.InfoLog.Println(" 获取到 " + monitorjson.Time + " 的 " + monitorjson.IP + " 的数据：" + monitorjson.Info)
}

package CallPolice

import (
	"../Global"
)

func Judge(monitorjson Global.MonitorJson){
	if ! monitorjson.Status {
		CallPolice(monitorjson.Hostname+" 的 "+ monitorjson.Info + " 异常 \n具体服务器信息：\n   " +monitorjson.IP)
	}
}
package Api

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	return router
}

func NewPublicRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes_pulic {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	return router
}

// 对内全功能接口
var routes = Routes{
	// 获取指定日期的监控信息，格式   url?time=年.月.日&key=关键字
	Route{
		"GetInfo",
		"GET",
		"/monitor/info",
		ReturnMonitorInfo,
	},

	// 获取指定日期的指定监控信息值，格式   url?time=年.月.日&key=关键字
	Route{
		"GetInfo",
		"GET",
		"/monitor/infolist",
		ReturnMonitorInfoList,
	},

	// 收集监控数据并保存处理
	Route{
		"MonitorCollect",
		"POST",
		"/monitor/collect",
		PostMonitorInfo,
	},

	// 收集nat机器中记录的hosts信息
	Route{
		"MonitorNatInfo",
		"POST",
		"/monitor/nat",
		PostNatInfo,
	},

	// 新增nat机器
	Route{
		"AddNat",
		"POST",
		"/admin/add",
		AddNatMonitor,
	},

	// 修改nat机器
	Route{
		"UpdateNat",
		"POST",
		"/admin/updata",
		UpdataNatMonitor,
	},

	// 删除nat机器
	Route{
		"DeleteNat",
		"POST",
		"/admin/del",
		DeleteNatMonitor,
	},

	// 获取所有nat机器信息
	Route{
		"SelectNat",
		"GET",
		"/admin/select",
		SelectNatMonitor,
	},

	// 重新读取crontab信息
	Route{
		"RestartCrontab",
		"POST",
		"/crontab/reload",
		ReloadCrontabNat,
	},

	// 控制coco机器监控某nat机器，格式  url?nat=IP地址
	Route{
		"MonitorNat",
		"POST",
		"/monitor/check_nat",
		ReturnNatMonitor,
	},

	// 控制coco机器监控所有nat机器
	Route{
		"MonitorAllNat",
		"POST",
		"/monitor/all",
		ReturnAllNatMonitor,
	},

	// 获取监控脚本
	Route{
		"MonitorShell",
		"GET",
		"/shell/monitor",
		ReturnMonitorShell,
	},

	// 获取nat机器提权脚本
	Route{
		"NatShell",
		"GET",
		"/shell/nat",
		ReturnNatShell,
	},

	// 修改是否报警的状态
	Route{
		"NatShell",
		"POST",
		"/police/change",
		ChangPoliceStatus,
	},

	// 获取是否报警的状态
	Route{
		"NatShell",
		"GET",
		"/police/status",
		ReturnPoliceStatus,
	},

	// 获取异常列表
	Route{
		"ErrorMap",
		"GET",
		"/police/map",
		ReturnPoliceMap,
	},

	// 获取hosts列表队列的状态
	Route{
		"NatHostsMap",
		"GET",
		"/monitor/map",
		ReturnNatHostsMap,
	},
}

// 暴露对外的相关模块
var routes_pulic = Routes{
	// 收集监控数据并保存处理
	Route{
		"MonitorCollect",
		"POST",
		"/monitor/collect",
		PostMonitorInfo,
	},

	// 收集nat机器中记录的hosts信息
	Route{
		"MonitorNatInfo",
		"POST",
		"/monitor/nat",
		PostNatInfo,
	},

	// 获取监控脚本
	Route{
		"MonitorShell",
		"GET",
		"/shell/monitor",
		ReturnMonitorShell,
	},

	// 获取nat机器提权脚本
	Route{
		"NatShell",
		"GET",
		"/shell/nat",
		ReturnNatShell,
	},
}

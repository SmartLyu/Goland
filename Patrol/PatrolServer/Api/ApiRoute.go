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

var routes = Routes{
	// 获取指定日期的监控信息，格式   url?time=年.月.日&key=关键字
	Route{
		"GetInfo",
		"GET",
		"/monitor/info",
		ReturnMonitorInfo,
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
		"/monitor/nat",
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
}
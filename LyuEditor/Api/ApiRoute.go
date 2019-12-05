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

// 对内全功能接口
var routes = Routes{
	// 获取指定html文件信息
	Route{
		"GetInfo",
		"POST",
		"/gethtml/file",
		GetHtmlFile,
	},
	// 获取成品url
	Route{
		"GetUrl",
		"POST",
		"/gethtml/list",
		GetHtmlList,
	},
}

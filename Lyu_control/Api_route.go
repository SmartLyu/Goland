package main
import (
	"github.com/gorilla/mux"
	"net/http"
	"time"
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
	// 微信相关接口
	Route{
		"GetInfo",
		"GET",
		"/shell",
		Test,
	},
	Route{
		"PostInfo",
		"POST",
		"/shell",
		Control,
	},
	Route{
		"PostInfo",
		"POST",
		"/shell/nat",
		ReturnNatShell,
	},
}

// 记录访问记录
func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		inner.ServeHTTP(w, r)

		WriteFile(r.Method+"\t"+r.RequestURI+"\t"+name+"\t"+time.Since(start).String())
	})
}

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

// 记录访问记录
func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		inner.ServeHTTP(w, r)

		WriteLog(r.Method + "\t" + r.RequestURI + "\t" + name + "\t" + time.Since(start).String())
	})
}

func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	route := routes
	var handler http.Handler
	handler = route.HandlerFunc
	handler = Logger(handler, route.Name)
	router.
		Methods(route.Method).
		Path(route.Pattern).
		Name(route.Name).
		Handler(handler)

	return router
}

var routes = Route{
	"StartMonitor",
	"Post",
	"/monitor",
	SshToNat,
}

func Api(port string) {
	router := NewRouter()
	WriteLog(http.ListenAndServe(":"+port, router).Error())
}

package Api

import (
	"../File"
	"net/http"
	"time"
)

// 记录访问记录
func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		inner.ServeHTTP(w, r)

		File.WriteAccessLog(r.Method + "\t" + r.RequestURI + "\t" + name + "\t" + time.Since(start).String())

	})
}

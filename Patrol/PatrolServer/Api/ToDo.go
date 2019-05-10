package Api

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

// 收集监控信息
func PostMonitorInfo(w http.ResponseWriter, r *http.Request) {
	var jsonfile MonitorJson

	// 读取用户post的信息
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	// 解析json格式信息
	if err := json.Unmarshal(body, &jsonfile); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	// 添加数据
	Jsons = append(Jsons, jsonfile)

	// 返回信息
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(jsonfile); err != nil {
		panic(err)
	}
}

// 展示监控信息
func ShowMonitorInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(Jsons); err != nil {
		panic(err)
	}
}
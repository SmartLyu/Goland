package Api

import (
	"../CallPolice"
	"../File"
	"../Global"
	"../Mysql"

	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

// 收集监控信息
func PostMonitorInfo(w http.ResponseWriter, r *http.Request) {
	var jsonfile Global.MonitorJson
	var hostjson Global.HostsTable

	// 读取用户post的信息
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		Global.ErrorLog.Println(err.Error())
	}
	if err := r.Body.Close(); err != nil {
		Global.ErrorLog.Println(err.Error())
	}

	// 解析json格式信息
	if err := json.Unmarshal(body, &jsonfile); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			Global.ErrorLog.Println(err.Error(), ", data is ", body)
		}
		if !jsonfile.Exist() {
			Global.InfoLog.Println("Error: has got empty json data")
		}
		return
	}

	// 添加数据
	go func() {
		if err := File.WriteFile(Global.ReadJson(jsonfile), jsonfile.Time); err != nil {
			Global.ErrorLog.Println("write info " + err.Error())
		}
		CallPolice.Judge(jsonfile)

		hostjson.IP = jsonfile.IP

		if jsonfile.Info == "survive" {
			Mysql.DeleteHosts(hostjson)
		}
	}()

	// 返回信息
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(jsonfile); err != nil {
		Global.ErrorLog.Println(err.Error())
	}
}

// 收集Nat后端机器信息
func PostNatInfo(w http.ResponseWriter, r *http.Request) {
	var jsonfile Global.HostsTable

	// 读取用户post的信息
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		Global.ErrorLog.Println(err.Error())
	}
	if err := r.Body.Close(); err != nil {
		Global.ErrorLog.Println(err.Error())
	}

	// 解析json格式信息
	if err := json.Unmarshal(body, &jsonfile); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			Global.ErrorLog.Println(err.Error(), ", data is ", body)
		}
	}

	// 添加数据
	Mysql.InsertHosts(jsonfile)

	// 返回信息
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(jsonfile); err != nil {
		Global.ErrorLog.Println(err.Error())
	}
}

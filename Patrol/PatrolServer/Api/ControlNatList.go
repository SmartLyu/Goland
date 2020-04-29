package Api

import (
	"../CallCoco"
	"../Global"
	"../Mysql"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

// 添加nat信息进入数据库
func AddNatMonitor(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json

	var jsonfile Global.NatTable

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
	if !Mysql.InsertNat(jsonfile) {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
	}

	// 启动监控巡查
	go CallCoco.CrontabToCallCoco(jsonfile)

	// 返回信息
	if err := json.NewEncoder(w).Encode(jsonfile); err != nil {
		Global.ErrorLog.Println(err.Error())
	}
}

// 删除nat信息进入数据库
func DeleteNatMonitor(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json

	var jsonfile Global.NatTable

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

	// 删除数据
	if !Mysql.DeleteNat(jsonfile) {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
	}

	// 停止监控
	CallCoco.StopCrontab(jsonfile)

	// 返回信息
	if err := json.NewEncoder(w).Encode(jsonfile); err != nil {
		Global.ErrorLog.Println(err.Error())
	}
}

// 修改nat信息进入数据库
func UpdataNatMonitor(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json

	var jsonfile Global.NatTable

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

	// 重新录入数据
	if !Mysql.DeleteNat(jsonfile) {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		if !Mysql.InsertNat(jsonfile) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusCreated)
		}
	}

	// 重启监控
	CallCoco.StopCrontab(jsonfile)
	go CallCoco.CrontabToCallCoco(jsonfile)

	// 返回信息
	if err := json.NewEncoder(w).Encode(jsonfile); err != nil {
		Global.ErrorLog.Println(err.Error())
	}
}

// 显示nat信息
func SelectNatMonitor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json

	jsonfile := Mysql.SelectAllNatTable()

	// 返回信息
	if err := json.NewEncoder(w).Encode(jsonfile); err != nil {
		w.WriteHeader(http.StatusNotFound)
		Global.ErrorLog.Println(err.Error())
	}
	w.WriteHeader(http.StatusOK)
}

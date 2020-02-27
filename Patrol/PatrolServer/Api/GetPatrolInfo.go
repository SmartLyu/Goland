package Api

import (
	"../CallCoco"
	"../File"
	"../Global"
	"../Mysql"
	"encoding/json"
	"net/http"
	"time"
)

// 重新启动crontab项目
func ReloadCrontabNat(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json

	jsonfile := Mysql.SelectAllNatTable()
	CallCoco.ReStartAllCrontab()

	// 返回信息
	if err := json.NewEncoder(w).Encode(jsonfile); err != nil {
		w.WriteHeader(http.StatusNotFound)
		Global.ErrorLog.Println(err.Error())
	}
	w.WriteHeader(http.StatusOK)
}

// 监控所有nat机器
func ReturnAllNatMonitor(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json

	//解析参数，默认是不会解析的
	err := r.ParseForm()
	if err != nil {
		Global.ErrorLog.Println("Recv: " + r.RemoteAddr)
	}

	nowTime := time.Now()
	CallCoco.CallAllNatMonitor()
	w.WriteHeader(http.StatusOK)

	time.Sleep(1 * time.Second)
	_, _ = w.Write([]byte("[{ \"readme\":\"！返回信息并不全面\"},"))
	// 获取当前日期时间信息
	pwd := Global.DataFileDir
	des := pwd + time.Now().Format("2006-01/02/15/04") + Global.DataFileName

	// 根据当前时间作为关键字检索文件
	data, err := File.FindWorkInFile(des)

	// 检查最近三秒的数据并返回
	if err != nil {
		Global.ErrorLog.Println("Find key Err: " + err.Error())
		http.NotFoundHandler().ServeHTTP(w, r)
	} else {
		Global.InfoLog.Println("Send File:" + des)
		_, err = w.Write([]byte(data))
		_, _ = w.Write([]byte("\n{\"readme\": \"已向所有后端机器发出监控请求，查看监控信息请 Get URL: " +
			" http://134.175.50.184:8666/monitor/info?" +
			"time=" + time.Now().Format("2006.01.02") + "&key=" +
			nowTime.Format("2006-01-02 15:04") + "\"}]"))
		if err != nil {
			Global.ErrorLog.Println("http writed Err " + err.Error())
		}
	}
	w.WriteHeader(http.StatusOK)
}

// 主动控制监控某个nat
func ReturnNatMonitor(w http.ResponseWriter, r *http.Request) {

	//解析参数，默认是不会解析的
	err := r.ParseForm()
	if err != nil {
		Global.ErrorLog.Println("Recv: " + r.RemoteAddr)
	}

	// 读取用户输入的参数信息
	getNat := r.Form.Get("nat")
	getHostname := r.Form.Get("host")
	getPort := r.Form.Get("port")

	nowTime := time.Now()
	CallCoco.CallCoco(getHostname, getNat, getPort)
	w.WriteHeader(http.StatusOK)

	time.Sleep(2 * time.Second)
	// 获取当前时间信息
	pwd := Global.DataFileDir
	des := pwd + time.Now().Format("2006-01/02/15/04") + Global.DataFileName

	// 根据当前时间作为关键字检索文件
	key := nowTime.Format("2006-01-02 15:04")
	data, err := File.FindWorkInFile(des, key)

	// 检查最近三秒的数据并返回
	if err != nil {
		Global.ErrorLog.Println("Find key Err: " + key + err.Error())
		http.NotFoundHandler().ServeHTTP(w, r)
	} else {
		Global.InfoLog.Println("Send File:" + des)
		_, err = w.Write([]byte(data))
		_, _ = w.Write([]byte("\n{\"readme\": \"已向所有后端机器发出监控请求，查看监控信息请 Get URL: " +
			" http://134.175.50.184:8666/monitor/info?" +
			"time=" + time.Now().Format("2006.01.02") + "&key=" +
			nowTime.Format("2006-01-02 15:04") + "\"}]"))
		if err != nil {
			Global.ErrorLog.Println("http writed Err " + err.Error())
		}
	}
}

// 显示主机队列信息
func ReturnNatHostsMap(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json

	var jsonfiles []Global.HostsTable
	var ht Global.HostsTable

	for key, value := range Global.NatHostsMap.Data {
		ht.HostName = value.Hostname
		ht.IP = key
		jsonfiles = append(jsonfiles, ht)
	}

	// 返回信息
	if err := json.NewEncoder(w).Encode(jsonfiles); err != nil {
		w.WriteHeader(http.StatusNotFound)
		Global.ErrorLog.Println(err.Error())
	}
	w.WriteHeader(http.StatusOK)
}

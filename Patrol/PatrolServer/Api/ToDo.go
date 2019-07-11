package Api

import (
	"../CallCoco"
	"../CallPolice"
	"../File"
	"../Global"
	"../Mysql"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// 收集监控信息
func PostMonitorInfo(w http.ResponseWriter, r *http.Request) {
	var jsonfile Global.MonitorJson
	var hostjson Global.HostsTable

	// 读取用户post的信息
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		File.WriteErrorLog(err.Error())
	}
	if err := r.Body.Close(); err != nil {
		File.WriteErrorLog(err.Error())
	}

	// 解析json格式信息
	if err := json.Unmarshal(body, &jsonfile); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			File.WriteErrorLog(err.Error())
		}
	}

	// 添加数据
	if err := File.WriteFile(Global.ReadJson(jsonfile)); err != nil {
		File.WriteErrorLog(err.Error())
	}
	CallPolice.Judge(jsonfile)

	hostjson.IP = jsonfile.IP
	hostjson.Time = time.Now().Format("2006-01-02 15:04")
	Mysql.DeleteHosts(hostjson)

	// 返回信息
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(jsonfile); err != nil {
		File.WriteErrorLog(err.Error())
	}
}

// 收集Nat后端机器信息
func PostNatInfo(w http.ResponseWriter, r *http.Request) {
	var jsonfile Global.HostsTable

	// 读取用户post的信息
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		File.WriteErrorLog(err.Error())
	}
	if err := r.Body.Close(); err != nil {
		File.WriteErrorLog(err.Error())
	}

	// 解析json格式信息
	if err := json.Unmarshal(body, &jsonfile); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			File.WriteErrorLog(err.Error())
		}
	}

	// 添加数据
	jsonfile.Time = time.Now().Format("2006-01-02 15:04")
	Mysql.InsertHosts(jsonfile)

	// 返回信息
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(jsonfile); err != nil {
		File.WriteErrorLog(err.Error())
	}
}

// 展示监控信息
func ReturnMonitorInfo(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json

	//解析参数，默认是不会解析的
	err := r.ParseForm()
	if err != nil {
		File.WriteErrorLog("Recv:" + r.RemoteAddr)
	}

	pwd := Global.DataFileDir

	// 读取用户输入的参数信息
	getTime := r.Form.Get("time")
	getTimeTmp := strings.Split(getTime, ".")
	if len(getTimeTmp) != 3 {
		File.WriteErrorLog("Get Monitor Info ,But Time stye Error")
		http.NotFoundHandler().ServeHTTP(w, r)
	}
	des := pwd + getTimeTmp[0] + "-" + getTimeTmp[1] + string(os.PathSeparator) + getTimeTmp[2] + Global.DataFileName

	getKey1 := r.Form.Get("key1")
	getKey2 := r.Form.Get("key2")
	getKey3 := r.Form.Get("key3")

	// 读取监控存储的文件
	desStat, err := os.Stat(des)
	if err != nil {
		File.WriteErrorLog("File Not Exit" + des)
		http.NotFoundHandler().ServeHTTP(w, r)
	} else if (desStat.IsDir()) {
		File.WriteErrorLog("File Is Dir" + des)
		http.NotFoundHandler().ServeHTTP(w, r)
	} else {
		_, _ = w.Write([]byte("["))
		// 根据用户是否有检索需求不同处理
		if getKey1 == "" {
			fileData, err := ioutil.ReadFile(des)
			if err != nil {
				File.WriteErrorLog("Read File Err: " + err.Error())
				http.NotFoundHandler().ServeHTTP(w, r)
			} else {
				File.WriteInfoLog("Send File:" + des)
				_, err = w.Write(fileData)
				if err != nil {
					File.WriteErrorLog("http writed Err " + err.Error())
				}
			}
		} else {
			data, err := File.FindWorkInFile(des, getKey1, getKey2, getKey3)
			if err != nil {
				File.WriteInfoLog("Find key Err " + err.Error())
				http.NotFoundHandler().ServeHTTP(w, r)
			} else {
				File.WriteInfoLog("Send File:" + des)
				_, err = w.Write([]byte(data))
				if err != nil {
					File.WriteErrorLog("http writed Err" + err.Error())
				}
			}
		}
		_, _ = w.Write([]byte("\n{\"Time\":\"" + time.Now().Format("2006-01-02 15:04") +
			"\",\"IP\":\"127.0.0.1\", \"Hostname\":\"JH-Api-QCloudGZ3-Patrol\"," +
			" \"Info\":\"patrol\", \"Status\":true}]"))
	}
}

// 添加nat信息进入数据库
func AddNatMonitor(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json

	var jsonfile Global.NatTable

	// 读取用户post的信息
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		File.WriteErrorLog(err.Error())
	}
	if err := r.Body.Close(); err != nil {
		File.WriteErrorLog(err.Error())
	}

	// 解析json格式信息
	if err := json.Unmarshal(body, &jsonfile); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			File.WriteErrorLog(err.Error())
		}
	}

	// 添加数据
	if ! Mysql.InsertNat(jsonfile) {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
	}

	// 启动监控巡查
	go CallCoco.CrontabToCallCoco(jsonfile)

	// 返回信息
	if err := json.NewEncoder(w).Encode(jsonfile); err != nil {
		File.WriteErrorLog(err.Error())
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
		File.WriteErrorLog(err.Error())
	}
	if err := r.Body.Close(); err != nil {
		File.WriteErrorLog(err.Error())
	}

	// 解析json格式信息
	if err := json.Unmarshal(body, &jsonfile); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			File.WriteErrorLog(err.Error())
		}
	}

	// 删除数据
	if ! Mysql.DeleteNat(jsonfile) {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
	}

	// 停止监控
	CallCoco.StopCrontab(jsonfile)

	// 返回信息
	if err := json.NewEncoder(w).Encode(jsonfile); err != nil {
		File.WriteErrorLog(err.Error())
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
		File.WriteErrorLog(err.Error())
	}
	if err := r.Body.Close(); err != nil {
		File.WriteErrorLog(err.Error())
	}

	// 解析json格式信息
	if err := json.Unmarshal(body, &jsonfile); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			File.WriteErrorLog(err.Error())
		}
	}

	// 重新录入数据
	if ! Mysql.DeleteNat(jsonfile) {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		if ! Mysql.InsertNat(jsonfile) {
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
		File.WriteErrorLog(err.Error())
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
		File.WriteErrorLog(err.Error())
	}
	w.WriteHeader(http.StatusOK)
}

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
		File.WriteErrorLog(err.Error())
	}
	w.WriteHeader(http.StatusOK)
}

// 返回监控脚本
func ReturnMonitorShell(w http.ResponseWriter, r *http.Request) {
	//解析参数，默认是不会解析的
	err := r.ParseForm()
	if err != nil {
		File.WriteErrorLog("Recv:" + r.RemoteAddr)
	}

	des := Global.MonitorShellFile
	desStat, err := os.Stat(des)
	if err != nil {
		File.WriteErrorLog("File Not Exit " + des)
		http.NotFoundHandler().ServeHTTP(w, r)
	} else if (desStat.IsDir()) {
		File.WriteErrorLog("File Is Dir" + des)
		http.NotFoundHandler().ServeHTTP(w, r)
	} else {
		data, err := ioutil.ReadFile(des)
		if err != nil {
			File.WriteErrorLog("Read File Err: " + err.Error())
		} else {
			File.WriteInfoLog("Send File:" + des)
			_, err = w.Write([]byte(data))
			if err != nil {
				File.WriteErrorLog("http writed Err " + err.Error())
			}
		}
	}
}

// 返回nat机器提权脚本
func ReturnNatShell(w http.ResponseWriter, r *http.Request) {
	//解析参数，默认是不会解析的
	err := r.ParseForm()
	if err != nil {
		File.WriteErrorLog("Recv:" + r.RemoteAddr)
	}

	des := Global.NatShellFile
	desStat, err := os.Stat(des)
	if err != nil {
		File.WriteErrorLog("File Not Exit " + des)
		http.NotFoundHandler().ServeHTTP(w, r)
	} else if (desStat.IsDir()) {
		File.WriteErrorLog("File Is Dir" + des)
		http.NotFoundHandler().ServeHTTP(w, r)
	} else {
		data, err := ioutil.ReadFile(des)
		if err != nil {
			File.WriteErrorLog("Read File Err: " + err.Error())
		} else {
			File.WriteInfoLog("Send File: " + des)
			_, err = w.Write([]byte(data))
			if err != nil {
				File.WriteErrorLog("http writed Err " + err.Error())
			}
		}
	}
}

// 监控所有nat机器
func ReturnAllNatMonitor(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json

	//解析参数，默认是不会解析的
	err := r.ParseForm()
	if err != nil {
		File.WriteErrorLog("Recv: " + r.RemoteAddr)
	}

	nowTime := time.Now()
	CallCoco.CallAllNatMonitor()
	w.WriteHeader(http.StatusOK)

	time.Sleep(1 * time.Second)
	_, _ = w.Write([]byte("[{ \"readme\":\"！返回信息并不全面\"},"))
	// 获取当前日期时间信息
	pwd := Global.DataFileDir
	des := pwd + time.Now().Format("2006-01/02") + Global.DataFileName

	// 根据当前时间作为关键字检索文件
	key := nowTime.Format("2006-01-02 15:04")
	data, err := File.FindWorkInFile(des, key)

	// 检查最近三秒的数据并返回
	if err != nil {
		File.WriteErrorLog("Find key Err: " + key + err.Error())
		http.NotFoundHandler().ServeHTTP(w, r)
	} else {
		File.WriteInfoLog("Send File:" + des)
		_, err = w.Write([]byte(data))
		_, _ = w.Write([]byte("\n{\"readme\": \"已向所有后端机器发出监控请求，查看监控信息请 Get URL: " +
			" http://134.175.50.184:8666/monitor/info?" +
			"time=" + time.Now().Format("2006.01.02") + "&key=" +
			nowTime.Format("2006-01-02 15:04") + "\"}]"))
		if err != nil {
			File.WriteErrorLog("http writed Err " + err.Error())
		}
	}
	w.WriteHeader(http.StatusOK)
}

// 主动控制监控某个nat
func ReturnNatMonitor(w http.ResponseWriter, r *http.Request) {

	//解析参数，默认是不会解析的
	err := r.ParseForm()
	if err != nil {
		File.WriteErrorLog("Recv: " + r.RemoteAddr)
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
	des := pwd + time.Now().Format("2006-01/02") + Global.DataFileName

	// 根据当前时间作为关键字检索文件
	key := nowTime.Format("2006-01-02 15:04")
	data, err := File.FindWorkInFile(des, key)

	// 检查最近三秒的数据并返回
	if err != nil {
		File.WriteErrorLog("Find key Err: " + key + err.Error())
		http.NotFoundHandler().ServeHTTP(w, r)
	} else {
		File.WriteInfoLog("Send File:" + des)
		_, err = w.Write([]byte(data))
		_, _ = w.Write([]byte("\n{\"readme\": \"已向所有后端机器发出监控请求，查看监控信息请 Get URL: " +
			" http://134.175.50.184:8666/monitor/info?" +
			"time=" + time.Now().Format("2006.01.02") + "&key=" +
			nowTime.Format("2006-01-02 15:04") + "\"}]"))
		if err != nil {
			File.WriteErrorLog("http writed Err " + err.Error())
		}
	}
}

// 修改是否报警
func ChangPoliceStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json

	if Global.IsPolice {
		CallPolice.CallPolice(time.Now().Format("2006.01.02 15:04") + "\n  现在进入维护模式，暂停报警功能")
		Global.IsPolice = false
		go func() {
			time.Sleep(time.Hour * 6)
			Global.IsPolice = true
			CallPolice.CallPolice(time.Now().Format("2006.01.02 15:04") + "\n  现在结束维护功能，开启报警功能")
		}()

	} else {
		Global.IsPolice = true
		CallPolice.CallPolice(time.Now().Format("2006.01.02 15:04") + "\n  现在结束维护功能，开启报警功能")
	}

	// 返回信息
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
}

// 显示是否报警
func ReturnPoliceStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json

	jsonfile := "{\"status\": \"" + strconv.FormatBool(Global.IsPolice) + "\"}"

	// 返回信息
	if err := json.NewEncoder(w).Encode(jsonfile); err != nil {
		w.WriteHeader(http.StatusNotFound)
		File.WriteErrorLog(err.Error())
	}
	w.WriteHeader(http.StatusOK)
}

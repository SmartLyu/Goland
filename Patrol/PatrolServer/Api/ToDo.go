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
	"sort"
	"strconv"
	"strings"
	"sync"
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
		if !jsonfile.Exist() {
			File.WriteInfoLog("Error: has got empty json data")
		}
		return
	}

	// 添加数据
	go func() {
		if err := File.WriteFile(Global.ReadJson(jsonfile), jsonfile.Time); err != nil {
			File.WriteErrorLog("write info " + err.Error())
		}
		CallPolice.Judge(jsonfile)

		hostjson.IP = jsonfile.IP
		hostjson.Time = time.Now().Format("2006-01-02 15:04")

		if jsonfile.Info == "survive" {
			Mysql.DeleteHosts(hostjson)
		}
	}()

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
	// 存储查询数据
	var monitorJsons []Global.MonitorJson
	// 总计查询数据个数
	var SearchNumber int64 = 0
	// 记录查询多少空值
	var emptyNumber = 0
	// 记录是否是预估值
	var approximate = "="
	// 处理时区问题
	loc, _ := time.LoadLocation("Asia/Shanghai")
	// 线程等待
	var goSync sync.WaitGroup
	// 修改锁定
	var Numberlock sync.Mutex
	var Datalock sync.Mutex
	jsonString := "["

	// 读取用户输入的时间参数信息
	getTimeStart, err := time.Parse(time.RFC3339, r.Form.Get("time1"))
	if err != nil {
		File.WriteErrorLog("Get Monitor Start Day Info: " + r.Form.Get("time1") + " , But Time stye Error")
		http.NotFoundHandler().ServeHTTP(w, r)
		return
	}
	getTimeStart = getTimeStart.In(loc)

	getTimeEnd, err := time.Parse(time.RFC3339, r.Form.Get("time2"))
	if err != nil {
		File.WriteErrorLog("Get Monitor End Day Info: " + r.Form.Get("time2") + " , But Time stye Error")
		http.NotFoundHandler().ServeHTTP(w, r)
		return
	}
	getTimeEnd = getTimeEnd.In(loc)

	// 解析获取指定范围内所有日期节点
	var getTime = make([]Global.DateTimeStyle, 2, 2)

	if getTimeStart.Unix() > getTimeEnd.Unix() {
		getTimeStart, getTimeEnd = getTimeEnd, getTimeStart
	}

	getTime = Global.GetAllTime(getTimeStart.Unix(), getTimeEnd.Unix())

	// 读取用户输入的关键字参数信息
	var getKey = make([]string, 5, 5)
	var keyTmpNumber = 1
	for {
		getKey = append(getKey, r.Form.Get("key"+strconv.Itoa(keyTmpNumber)))
		if r.Form.Get("key"+strconv.Itoa(keyTmpNumber+1)) == "" {
			break
		}
		keyTmpNumber++
		if keyTmpNumber >= 10000 {
			break
		}
	}

	// 读取用户正则匹配表达式
	getGexp := r.Form.Get("gexp")

	// 读取监控存储的文件
	// 按天循环获取监控信息
	for tmpNumber, tmpTime := range getTime {
		if !tmpTime.Exist() {
			emptyNumber++
			continue
		}
		Global.MaxSearchSigLen <- 1
		if SearchNumber >= Global.MaxReturnLen {
			SearchNumber = SearchNumber / int64(tmpNumber) * int64(len(getTime)-emptyNumber)
			approximate = "~"
			break
		}
		goSync.Add(1)
		goFunc := func(tmpTime Global.DateTimeStyle) {
			var data string
			des := pwd + tmpTime.Year + "-" + tmpTime.Month + string(os.PathSeparator) + tmpTime.Day +
				string(os.PathSeparator) + tmpTime.Hour + string(os.PathSeparator) + tmpTime.Minute + Global.DataFileName
			_, err := os.Stat(des)
			if err == nil {
				var tmpSearchNumber int64 = 0
				tmpSearchNumber, data, err = File.SearchWordInFile(des, getGexp, getKey...)
				if err == nil {
					if SearchNumber <= Global.MaxSearchLen {
						Datalock.Lock()
						jsonString = jsonString + data
						Datalock.Unlock()
					}
					Numberlock.Lock()
					SearchNumber += tmpSearchNumber
					Numberlock.Unlock()
				}
			}
			goSync.Done()
			<-Global.MaxSearchSigLen
		}
		go goFunc(tmpTime)
	}
	goSync.Wait()

	jsonString = jsonString + "\n{\"Time\":\"" + time.Now().Format("2006-01-02 15:04") +
		"\",\"IP\":\"0.0.0.0\", \"Hostname\":\"All-JH-Api-QCloudGZ3-Patrol\"," +
		"\"Info\":\"AllHitSearch_Patrol_Info" + approximate +
		strconv.FormatInt(SearchNumber, 10) + "\", \"Status\":true}]"

	monitorJsons = Global.ReadMonitorJson(jsonString)
	if err := json.NewEncoder(w).Encode(monitorJsons); err != nil {
		w.WriteHeader(http.StatusNotFound)
		File.WriteErrorLog(err.Error())
	}
	w.WriteHeader(http.StatusOK)
	File.WriteInfoLog("Get Patrol info all hit search is " + strconv.FormatInt(SearchNumber, 10) +
		", from " + getTimeStart.Format("2006-01-02:15:04") +
		" to " + getTimeEnd.Format("2006-01-02:15:04"))
}

// 展示监控信息指定信息值
func ReturnMonitorInfoList(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json

	//解析参数，默认是不会解析的
	err := r.ParseForm()
	if err != nil {
		File.WriteErrorLog("Recv:" + r.RemoteAddr)
	}

	pwd := Global.DataFileDir
	// 存储查询数据
	var monitorJsons []Global.MonitorJson
	// 总计查询数据个数
	var SearchNumber int64 = 0
	// 记录查询多少空值
	var emptyNumber = 0
	// 存储list使用的数据处理
	type infoListJson struct {
		value map[string]int
	}

	var returnJson = make(map[string]infoListJson)
	var sorted_keys = make([]string, 0)

	// 处理时区问题
	loc, _ := time.LoadLocation("Asia/Shanghai")
	var data string
	jsonString := "["

	// 读取用户输入的时间参数信息
	getTimeStart, err := time.Parse(time.RFC3339, r.Form.Get("time1"))
	if err != nil {
		File.WriteErrorLog("Get Monitor Start Day Info: " + r.Form.Get("time1") + " , But Time stye Error")
		http.NotFoundHandler().ServeHTTP(w, r)
		return
	}
	getTimeStart = getTimeStart.In(loc)

	getTimeEnd, err := time.Parse(time.RFC3339, r.Form.Get("time2"))
	if err != nil {
		File.WriteErrorLog("Get Monitor End Day Info: " + r.Form.Get("time2") + " , But Time stye Error")
		http.NotFoundHandler().ServeHTTP(w, r)
		return
	}
	getTimeEnd = getTimeEnd.In(loc)

	// 解析获取指定范围内所有日期节点
	var getTime = make([]Global.DateTimeStyle, 2, 2)

	if getTimeStart.Unix() > getTimeEnd.Unix() {
		getTimeStart, getTimeEnd = getTimeEnd, getTimeStart
	}

	getTime = Global.GetAllTime(getTimeStart.Unix(), getTimeEnd.Unix())

	// 读取用户输入的关键字参数信息
	var getKey = make([]string, 5, 5)
	var keyTmpNumber = 1
	for {
		getKey = append(getKey, r.Form.Get("key"+strconv.Itoa(keyTmpNumber)))
		if r.Form.Get("key"+strconv.Itoa(keyTmpNumber+1)) == "" {
			break
		}
		keyTmpNumber++
		if keyTmpNumber >= 10000 {
			break
		}
	}

	// 读取用户正则匹配表达式
	getGexp := r.Form.Get("gexp")

	// 读取监控存储的文件
	// 按天循环获取监控信息
	for _, tmpTime := range getTime {
		if !tmpTime.Exist() {
			emptyNumber++
			continue
		}
		des := pwd + tmpTime.Year + "-" + tmpTime.Month + string(os.PathSeparator) + tmpTime.Day +
			string(os.PathSeparator) + tmpTime.Hour + string(os.PathSeparator) + tmpTime.Minute + Global.DataFileName
		_, err := os.Stat(des)
		if err != nil {
			continue
		} else {
			var tmpSearchNumber int64 = 0
			tmpSearchNumber, data, err = File.SearchWordInFile(des, getGexp, getKey...)
			if err != nil {
				File.WriteInfoLog("Find key Err " + err.Error())
			} else {
				File.WriteInfoLog("Send File:" + des + ", daytime： " + tmpTime.Print())
				if SearchNumber <= Global.MaxSearchLen*10 {
					jsonString = jsonString + data
				}
				if tmpSearchNumber > Global.MaxSearchLen {
					SearchNumber += Global.MaxSearchLen
				} else {
					SearchNumber += tmpSearchNumber
				}

				if SearchNumber > Global.MaxSearchLen*10 {
					break
				}
				if err != nil {
					File.WriteErrorLog("http writed Err" + err.Error())
				}
			}
		}
	}

	jsonString = jsonString + "\n{\"Time\":\"" + time.Now().Format("2006-01-02 15:04") +
		"\",\"IP\":\"0.0.0.0\", \"Hostname\":\"All-JH-Api-QCloudGZ3-Patrol\"," +
		"\"Info\":\"AllHitSearch_Patrol_Info" + "\", \"Status\":true}]"

	// 解析查出的json数据
	monitorJsons = Global.ReadMonitorJson(jsonString)
	for _, j := range monitorJsons {
		if len(returnJson[j.Time].value) < 1 {
			returnJson[j.Time] = infoListJson{
				value: make(map[string]int),
			}
		}
		tmpString := strings.Split(j.Info, "=")
		if len(tmpString) == 1 {
			continue
		}
		tmpInt, err := strconv.Atoi(strings.Split(tmpString[1], "%")[0])
		if err != nil {
			continue
		}
		returnJson[j.Time].value[j.Hostname+"="+strings.Split(j.Info, "=")[0]] = tmpInt
	}

	// 判断是否检错出有用数据
	if len(returnJson) < 2 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// 按照时间顺序排序
	for k, _ := range returnJson {
		sorted_keys = append(sorted_keys, k)
	}
	sort.Strings(sorted_keys)

	// 遍历书写不定json数据
	var returnFistJson = make([]string, 2, 2)
	var returnJsonString = ""
	for _, i := range sorted_keys {
		j := returnJson[i]
		var isGetValue = false
		for key, value := range j.value {
			if len(returnFistJson) < 10 {
				returnFistJson = append(returnFistJson, key)
			}
			if !isGetValue {
				returnJsonString = returnJsonString + ", { \"Time\": \"" + i + "\", "
				isGetValue = true
			}
			returnJsonString = returnJsonString + " \"" + key + "\": " + strconv.Itoa(value) + ","
		}
		if isGetValue {
			returnJsonString = strings.TrimRight(returnJsonString, ",") + " }"
		}
	}
	returnFistJson = Global.RemoveRepeatedElement(returnFistJson)

	// 按照前端需求，组合所有key信息
	returntring := "[ { \"Value\": [ \"Time\","
	for _, i := range returnFistJson {
		if i == "" {
			continue
		}
		returntring = returntring + " \"" + i + "\","
	}

	returntring = strings.TrimRight(returntring, ",") + " ]} " + returnJsonString + "]"

	// 发送json数据
	_, err = w.Write([]byte(returntring))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		File.WriteErrorLog("http writed Err" + err.Error())
	}
	w.WriteHeader(http.StatusOK)
	File.WriteInfoLog("Get Patrol info all hit search is " + strconv.FormatInt(SearchNumber, 10) +
		", from " + getTimeStart.Format("2006-01-02:15:04") +
		" to " + getTimeEnd.Format("2006-01-02:15:04"))
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
	} else if desStat.IsDir() {
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
	} else if desStat.IsDir() {
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
	des := pwd + time.Now().Format("2006-01/02/15/04") + Global.DataFileName

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
	des := pwd + time.Now().Format("2006-01/02/15/04") + Global.DataFileName

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
		CallPolice.CallMessage("  现在进入维护模式，暂停报警功能")
		Global.IsPolice = false
		go func() {
			time.Sleep(time.Hour * 6)
			if !Global.IsPolice {
				Global.IsPolice = true
				CallPolice.CallMessage("  现在结束维护功能，开启报警功能")
			}
		}()

	} else {
		Global.IsPolice = true
		CallPolice.CallMessage("  现在结束维护功能，开启报警功能")
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

// 显示报警队列信息
func ReturnPoliceMap(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json

	var jsonfiles []Global.ErrorJson
	allErrorNum := 0

	// map数据导出为json数据
	for key, value := range Global.ErrorMap.Data {
		var errJson Global.ErrorJson
		errJson.Key = key
		errJson.Value = value
		allErrorNum += value
		jsonfiles = append(jsonfiles, errJson)
	}

	// 最后统计未处理异常次数
	var errJson Global.ErrorJson
	errJson.Key = "Get All Patrol Error map Finished, Total Error Number"
	errJson.Value = allErrorNum
	jsonfiles = append(jsonfiles, errJson)

	// 返回信息
	if err := json.NewEncoder(w).Encode(jsonfiles); err != nil {
		w.WriteHeader(http.StatusNotFound)
		File.WriteErrorLog(err.Error())
	}
	w.WriteHeader(http.StatusOK)
}

// 显示主机队列信息
func ReturnNatHostsMap(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json

	var jsonfiles []Global.HostsTable

	for key, _ := range Global.NatHostsMap.Data {
		jsonfiles = append(jsonfiles, key)
	}

	// 返回信息
	if err := json.NewEncoder(w).Encode(jsonfiles); err != nil {
		w.WriteHeader(http.StatusNotFound)
		File.WriteErrorLog(err.Error())
	}
	w.WriteHeader(http.StatusOK)
}

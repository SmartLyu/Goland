package Api

import (
	"../File"
	"../Global"
	"encoding/json"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// 展示监控信息
func ReturnMonitorInfo(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json

	//解析参数，默认是不会解析的
	err := r.ParseForm()
	if err != nil {
		Global.ErrorLog.Println("Recv:" + r.RemoteAddr)
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
		Global.ErrorLog.Println("Get Monitor Start Day Info: " + r.Form.Get("time1") + " , But Time stye Error")
		http.NotFoundHandler().ServeHTTP(w, r)
		return
	}
	getTimeStart = getTimeStart.In(loc)

	getTimeEnd, err := time.Parse(time.RFC3339, r.Form.Get("time2"))
	if err != nil {
		Global.ErrorLog.Println("Get Monitor End Day Info: " + r.Form.Get("time2") + " , But Time stye Error")
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
		Global.ErrorLog.Println(err.Error())
	}
	w.WriteHeader(http.StatusOK)
	Global.InfoLog.Println("Get Patrol info all hit search is " + strconv.FormatInt(SearchNumber, 10) +
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
		Global.ErrorLog.Println("Recv:" + r.RemoteAddr)
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
		Global.ErrorLog.Println("Get Monitor Start Day Info: " + r.Form.Get("time1") + " , But Time stye Error")
		http.NotFoundHandler().ServeHTTP(w, r)
		return
	}
	getTimeStart = getTimeStart.In(loc)

	getTimeEnd, err := time.Parse(time.RFC3339, r.Form.Get("time2"))
	if err != nil {
		Global.ErrorLog.Println("Get Monitor End Day Info: " + r.Form.Get("time2") + " , But Time stye Error")
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
				Global.InfoLog.Println("Find key Err " + err.Error())
			} else {
				Global.InfoLog.Println("Send File:" + des + ", daytime： " + tmpTime.Print())
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
					Global.ErrorLog.Println("http writed Err" + err.Error())
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
		Global.ErrorLog.Println("http writed Err" + err.Error())
	}
	w.WriteHeader(http.StatusOK)
	Global.InfoLog.Println("Get Patrol info all hit search is " + strconv.FormatInt(SearchNumber, 10) +
		", from " + getTimeStart.Format("2006-01-02:15:04") +
		" to " + getTimeEnd.Format("2006-01-02:15:04"))
}

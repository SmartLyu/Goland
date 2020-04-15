package Api

import (
	"../CallPolice"
	"../Global"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

// 修改是否报警
func ChangPoliceStatus(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json

	if Global.IsPolice {
		CallPolice.CallMessage("现在进入维护模式，暂停报警功能")
		Global.IsPolice = false
		go func() {
			time.Sleep(time.Hour * 6)
			if !Global.IsPolice {
				Global.IsPolice = true
				CallPolice.CallMessage("现在结束维护功能，开启报警功能")
			}
		}()

	} else {
		Global.IsPolice = true
		CallPolice.CallMessage("现在结束维护功能，开启报警功能")
	}

	// 返回信息
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
}

// 显示是否报警
func ReturnPoliceStatus(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json

	jsonfile := "{\"status\": \"" + strconv.FormatBool(Global.IsPolice) + "\"}"

	// 返回信息
	if err := json.NewEncoder(w).Encode(jsonfile); err != nil {
		w.WriteHeader(http.StatusNotFound)
		Global.ErrorLog.Println(err.Error())
	}
	w.WriteHeader(http.StatusOK)
}

// 显示报警队列信息
func ReturnPoliceMap(w http.ResponseWriter) {
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
		Global.ErrorLog.Println(err.Error())
	}
	w.WriteHeader(http.StatusOK)
}

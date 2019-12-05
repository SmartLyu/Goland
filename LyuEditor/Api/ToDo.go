package Api

import (
	"../File"
	"../Global"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// 获取post文件信息
func GetHtmlFile(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "text/plain")                   //返回数据格式是txt

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		File.WriteErrorLog(err.Error())
	}

	randomFile := Global.RandStringRunes() + ".html"
	AbsoluteFile := time.Now().Format("2006-01") + "/" + randomFile

	err = File.WriteFile(string(body), randomFile)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		File.WriteErrorLog(err.Error())
		return
	}
	err = File.MyCmd("bash", Global.ShellFile, Global.DataFileDir+AbsoluteFile)
	if err != nil {
		w.WriteHeader(http.StatusFailedDependency)
		File.WriteErrorLog(err.Error())
		return
	}

	// 返回信息
	if err := json.NewEncoder(w).Encode(Global.GetURL + AbsoluteFile); err != nil {
		w.WriteHeader(http.StatusNotFound)
		File.WriteErrorLog(err.Error())
	}
	w.WriteHeader(http.StatusOK)
}

// 获取post文件信息
func GetHtmlList(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "text/plain")                   //返回数据格式是txt

	// 读取用户post的信息
	var jsonfile Global.ListJson
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
		return
	}
	if ! jsonfile.Exist() {
		File.WriteInfoLog("Error: has got empty json data")
		return
	}

	err = File.MyCmd("bash", Global.ShellList, jsonfile.YuUrl, jsonfile.CnName, jsonfile.EnName)
	if err != nil {
		w.WriteHeader(http.StatusFailedDependency)
		File.WriteErrorLog(err.Error())
		return
	}
	err = File.SendMail(jsonfile.EnName+".html")
	if err != nil {
		w.WriteHeader(http.StatusFailedDependency)
		File.WriteErrorLog(err.Error())
		return
	}

	// 返回信息
	if err := json.NewEncoder(w).Encode(jsonfile); err != nil {
		w.WriteHeader(http.StatusNotFound)
		File.WriteErrorLog(err.Error())
	}
	w.WriteHeader(http.StatusOK)
}
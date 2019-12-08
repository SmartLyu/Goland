package MonitorApi

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type AddProjectJson struct {
	Name    string `json:"name"`
	Project string `json:"project"`
}

type AddStrategyJson struct {
	Smetric          string  `json:"metric"`
	Sagent           string  `json:"agent"`
	Sfunc            string  `json:"func"`
	Sop              string  `json:"op"`
	Sright_value     float64 `json:"right_value"`
	Snote            string  `json:"note"`
	Snodata          bool    `json:"nodata"`
	Snodata_value    int     `json:"nodata_value"`
	Snodata_interval int     `json:"nodata_interval"`
	Sproject         string  `json:"project"`
}

type PostDataJson struct {
	Dagent     string  `json:"agent"`
	Dmetric    string  `json:"metric"`
	Dvalue     float32 `json:"value"`
	Dtimestamp int64   `json:"timestamp"`
}

type AddRoleJson struct {
	Ruletpl   string   `json:"ruletpl"`
	Rule_name string   `json:"rule_name"`
	Project   []string `json:"projects"`
}

type DeleteRoleToStrategy struct {
	Uuid_list []string `json:"uuid_list"`
}

type ReturnDeleteRole struct {
	Code    int    `json:"code"`
	Data    string `json:"data"`
	Message string `json:"message"`
}

func ReadUUidJsonString(jsonData string) (code string, dataString []string, errReturn error) {
	fmt.Println(time.Now().Format("2006-01-02 15:04:05 +0800 CST")+" ", "Return Json:", jsonData)
	var v interface{}
	errReturn = nil
	if err := json.Unmarshal([]byte(jsonData), &v); err != nil {
		errReturn = err
	}
	data_1, ok := v.(map[string]interface{})
	if !ok {
		errReturn = errors.New("返回数据无法解析")
	}
	code = fmt.Sprintf("%v", data_1["code"])
	data_2, ok := data_1["data"].(map[string]interface{})
	if !ok {
		errReturn = errors.New("返回数据无法解析")
	}
	data_3, ok := data_2["data"].([]interface{})
	if !ok {
		errReturn = errors.New("返回数据无法解析")
	}
	for _, v := range data_3 {
		data, ok := v.(map[string]interface{})
		if !ok {
			errReturn = errors.New("返回数据无法解析")
		}
		dataString = append(dataString, fmt.Sprintf("%v", data["uuid"]))
	}
	return
}

func ReadAddRuleJsonString(jsonData string) (dataString string, errReturn error) {
	fmt.Println(time.Now().Format("2006-01-02 15:04:05 +0800 CST")+" ", "Return Json:", jsonData)
	var v interface{}
	errReturn = nil
	if err := json.Unmarshal([]byte(jsonData), &v); err != nil {
		errReturn = err
	}
	data_1, ok := v.(map[string]interface{})
	if !ok {
		return "", errors.New("返回数据无法解析")
	}
	if fmt.Sprintf("%v", data_1["code"]) != "200" {
		errReturn = errors.New("add rules error: " + jsonData)
	}
	data_2, ok := data_1["data"].([]interface{})
	if !ok {
		return "", errors.New("返回数据无法解析")
	}
	for _, v := range data_2 {
		data := v.(map[string]interface{})
		dataString = fmt.Sprintf("%v", data["uuid"])
	}
	return
}

func ReadAddAgentsJsonString(jsonData string) (code string, dataString string, errReturn error) {
	fmt.Println(time.Now().Format("2006-01-02 15:04:05 +0800 CST")+" ", "Return Json:", jsonData)
	var v interface{}
	errReturn = nil
	if err := json.Unmarshal([]byte(jsonData), &v); err != nil {
		errReturn = err
	}
	data_1, ok := v.(map[string]interface{})
	if !ok {
		return "", "", errors.New("返回数据无法解析")
	}
	code = fmt.Sprintf("%v", data_1["code"])
	data_2, ok := data_1["data"].(map[string]interface{})
	if !ok {
		return "", "", errors.New("返回数据无法解析")
	}
	dataString = fmt.Sprintf("%v", data_2["uuid"])
	return
}

func ReadAddDataJsonString(jsonData string) bool {
	fmt.Println(time.Now().Format("2006-01-02 15:04:05 +0800 CST")+" ", "Return Json:", jsonData)
	var v interface{}
	if err := json.Unmarshal([]byte(jsonData), &v); err != nil {
		return false
	}
	data_1, ok := v.(map[string]interface{})
	if !ok {
		return false
	}
	_, ok = data_1["code"]
	if !ok {
		return false
	}
	if fmt.Sprintf("%v", data_1["code"]) == "200" {
		return true
	}
	return false
}

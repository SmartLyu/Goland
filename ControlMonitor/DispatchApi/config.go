package DispatchApi

import (
	"../Log"
	"encoding/json"
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
)

var (
	ctaskJson   createTaskJson
	jobJson     createJobJson
	changeJob   changeJobJson
	taskAndJson bindTaskAndJson
	sTaskJson   startTaskJson
	taskName    = ""
)

func ReturnTaskName() string {
	return "云游戏行为测试" + uuid.Must(uuid.NewV1()).String()
}

type createTaskJson struct {
	Name    string `json:"name"`
	Type    int    `json:"type"`
	Project string `json:"project"`
}

type createJobJson struct {
	Name     string `json:"name"`
	Command  string `json:"command"`
	Args     string `json:"args"`
	Project  string `json:"project"`
	IsRandom bool   `json:"is_random"`
	Quantity int    `json:"quantity"`
	TimeOut  int    `json:"time_out"`
	CphRes   string `json:"cph_res"`
}

type changeJobJson struct {
	Uuid     string `json:"uuid"`
	Command  string `json:"command"`
	Args     string `json:"args"`
	Project  string `json:"project"`
	IsRandom bool   `json:"is_random"`
	Quantity int    `json:"quantity"`
	TimeOut  int    `json:"time_out"`
	CphRes   string `json:"cph_res"`
}

type bindTaskAndJson struct {
	TaskId string   `json:"task_id"`
	JobId  []string `json:"job_ids"`
}

type startTaskJson struct {
	Uuid   string `json:"uuid"`
	Status int    `json:"status"`
}

func readReturnCreateTaskJsonString(jsonData string) (uuid string, errReturn error) {
	Log.DebugLog.Println("Return Json:" + jsonData)
	var v interface{}
	errReturn = nil
	if err := json.Unmarshal([]byte(jsonData), &v); err != nil {
		errReturn = err
	}
	data1, ok := v.(map[string]interface{})
	if !ok {
		return "", errors.New("创建Task返回数据无法解析")
	}
	data2, ok := data1["task"].(map[string]interface{})
	if !ok {
		return "", errors.New("创建Task返回数据无法解析")
	}
	returnStr, ok := data2["uuid"]
	if !ok {
		return "", errors.New("创建Task返回数据无法解析")
	}
	return fmt.Sprintf("%v", returnStr), nil
}

func readReturnCreateJobJsonString(jsonData string) (uuid string, errReturn error) {
	Log.DebugLog.Println("Return Json:" + jsonData)
	var v interface{}
	errReturn = nil
	if err := json.Unmarshal([]byte(jsonData), &v); err != nil {
		errReturn = err
	}
	data1, ok := v.(map[string]interface{})
	if !ok {
		return "", errors.New("创建Task返回数据无法解析")
	}
	data2, ok := data1["job"].(map[string]interface{})
	if !ok {
		return "", errors.New("创建Task返回数据无法解析")
	}
	returnStr, ok := data2["uuid"]
	if !ok {
		return "", errors.New("创建Task返回数据无法解析")
	}
	return fmt.Sprintf("%v", returnStr), nil
}

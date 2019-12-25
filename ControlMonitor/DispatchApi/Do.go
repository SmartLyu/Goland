package DispatchApi

import (
	"../Log"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
)

func StartTask(number int, timeOut int, project string, commod string, args string, url string) (string, error) {
	taskName = ReturnTaskName()
	ctaskJson = createTaskJson{
		Name:    taskName,
		Type:    2,
		Project: project,
	}

	js, _ := json.Marshal(&ctaskJson)
	taskJsonstr, _, err := httpJson(url+"/v1/task", "POST", js)
	if err != nil {
		taskName = ReturnTaskName()
		ctaskJson = createTaskJson{
			Name:    taskName,
			Type:    2,
			Project: project,
		}

		js, _ := json.Marshal(&ctaskJson)
		taskJsonstr, _, err = httpJson(url+"/v1/task", "POST", js)
		if err != nil {
			return "", err
		}
	}
	TaskUuid, err := readReturnCreateTaskJsonString(taskJsonstr)
	if err != nil {
		return "", err
	}

	// 计算 命令参数 的md5值并查看是否存在job
	Md5Inst := md5.New()
	Md5Inst.Write([]byte(commod + args))
	Result := fmt.Sprintf("%x", Md5Inst.Sum([]byte("")))
	getJobJsonstr, _, err := httpJson(url+"/v1/job?job.name="+Result, "GET", []byte(""))
	jobUuid, err := readReturnCreateJobJsonString(getJobJsonstr)
	if err != nil {
		// 创建job信息
		jobJson = createJobJson{
			Name:     Result,
			Command:  commod,
			Args:     args + TaskUuid,
			Project:  project,
			IsRandom: true,
			Quantity: number,
			TimeOut:  timeOut,
			CphRes:   "*",
		}
		js, _ = json.Marshal(&jobJson)
		jobJsonstr, _, err := httpJson(url+"/v1/job", "POST", js)
		if err != nil {
			return "", err
		}
		jobUuid, err = readReturnCreateJobJsonString(jobJsonstr)
		if err != nil {
			return "", err
		}
	} else {
		// 修改job信息
		changeJob = changeJobJson{
			Uuid:     jobUuid,
			Command:  commod,
			Args:     args + TaskUuid,
			Project:  project,
			IsRandom: true,
			Quantity: number,
			TimeOut:  timeOut,
			CphRes:   "*",
		}
		js, _ = json.Marshal(&changeJob)
		jobJsonstr, _, err := httpJson(url+"/v1/job", "PUT", js)
		if err != nil {
			return "", err
		}
		Log.DebugLog.Println(jobJsonstr)
		_, err = readReturnCreateJobJsonString(jobJsonstr)
		if err != nil {
			return "", err
		}
	}

	// 绑定task和job
	var jobMap = []string{
		jobUuid,
	}
	taskAndJson = bindTaskAndJson{
		TaskId: TaskUuid,
		JobId:  jobMap,
	}
	js, _ = json.Marshal(&taskAndJson)
	bindJsonstr, _, err := httpJson(url+"/v1/task/jobs", "POST", js)
	if err != nil {
		return "", err
	}
	if bindJsonstr != "{}" {
		return "", errors.New("绑定task和job异常" + TaskUuid + "和" + jobUuid)
	}

	// 开启task任务
	sTaskJson = startTaskJson{
		Uuid:   TaskUuid,
		Status: 1,
	}
	js, _ = json.Marshal(&sTaskJson)
	startJsonstr, _, err := httpJson(url+"/v1/task", "PUT", js)
	if err != nil {
		return "", err
	}
	_, err = readReturnCreateTaskJsonString(startJsonstr)
	if err != nil {
		return "", err
	}

	Log.InfoLog.Println("Success Get task UUid: " + TaskUuid + ", job UUid: " + jobUuid)
	return TaskUuid, nil
}

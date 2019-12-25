package main

import (
	"./CheckEsInfo"
	"./DispatchApi"
	"./Log"
	"./MonitorApi"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"runtime"
	"strconv"
	"time"
)

func main() {
	flag.StringVar(&configfile, "config", "error", "输入配置文件位置")
	flag.Parse()

	if configfile == "error" {
		flag.Usage()
		log.Fatal("请输入配置文件位置")
	}

	Load(configfile)
	err := Log.Log(config.LogDir)
	if err != nil {
		Log.ErrorLog.Println("LogCreate Error: " + err.Error())
		PoilceByPatrol(configfile, "false")
		return
	}
	switch config.Type {
	case "test":
		StartTest()
	case "monitor":
		StartMonitor()
	default:
		flag.Usage()
		log.Fatal("配置文件格式错误: Type is error (test , monitor)")
	}

}

func StartTest() {
	projectUuid := "e905c3a0335c11e9af90fa163ef8d597"
	var taskId = make([]string, 0)
	Log.InfoLog.Println("now start for test task, now project's goroutine: " +
		strconv.Itoa(runtime.NumGoroutine()))

	for i := 1; i <= config.TestCycleTimes; i++ {
		// 开始任务
		taskTmpId, err := DispatchApi.StartTask(config.Task1, config.TimeOut, projectUuid,
			config.Shell, config.Args, config.Url)
		if err != nil {
			Log.ErrorLog.Println("StartTask Error: " + err.Error())
			PoilceByPatrol(configfile+"="+taskTmpId, "false")
			return
		}
		taskId = append(taskId, taskTmpId)
		// 开始等待任务完成，记录任务结束时间
		Log.InfoLog.Println("Successfully start waiting for the " + strconv.Itoa(i) + " task: " + taskTmpId +
			" to start next in about " + strconv.Itoa(config.TestCycleWait) + "s")
		time.Sleep(time.Duration(config.TestCycleWait) * time.Second)
	}

	// 开始等待任务完成，记录任务结束时间
	Log.InfoLog.Println("Successfully start waiting for all " + strconv.Itoa(config.TestCycleTimes) +
		" tasks to finish in about " + strconv.Itoa(config.TimeOut+waiteTime) + "s")

	taskJson := taskConfigJSON{
		Type:          config.Type,
		AgentUuid:     "",
		Tasks:         config.Tasks,
		Success:       config.SuccessTip,
		TaskId:        taskId,
		IsMonitor:     false,
		AllNumber:     config.Task1,
		SuccessNumber: config.TaskSuccessNumber1,
		TimeOut:       config.TimeOut + waiteTime,
		CheckTime:     config.CheckTime,
	}
	js, _ := json.Marshal(&taskJson)
	status, err := CheckEsInfo.GetTaskInfo(js)

	if err != nil {
		Log.ErrorLog.Println("CheckEsInfo Error: " + err.Error() + ", wait " +
			strconv.Itoa(config.TimeOut) + "s for finishing task")
		PoilceByPatrol(configfile, "false")
		time.Sleep(time.Duration(config.TimeOut) * time.Second)
		return
	}
	Log.InfoLog.Println(projectUuid + " check status: " + strconv.FormatBool(status) + " is finished")
}

func StartMonitor() {
	for {
		Log.InfoLog.Println("now start for monitor task, now project's goroutine: " +
			strconv.Itoa(runtime.NumGoroutine()))
		Load(configfile)
		StartProject()
		sleepTime := config.CheckTime + 5
		Log.InfoLog.Println("wait " + strconv.Itoa(sleepTime) + "s to start next time, now project's goroutine: " +
			strconv.Itoa(runtime.NumGoroutine()))
		time.Sleep(time.Duration(sleepTime) * time.Second)
		fmt.Println()
	}
}

func StartProject() {
	projectUuid, agentsUuid, err := MonitorApi.DoDefault(config.Project, config.Place, config.IsMaintain)
	if err != nil {
		Log.ErrorLog.Println("MonitorApi DoDefault Error: " + err.Error())
		PoilceByPatrol(configfile, "false")
		return
	}
	fmt.Println()
	if config.IsMaintain {
		Log.InfoLog.Println("Now is maintain time, wait " + strconv.Itoa(config.TimeOut) + "s to read config ...")
		time.Sleep(time.Duration(config.TimeOut) * time.Second)
		return
	}

	Log.InfoLog.Println("Get ProjectUUid: " + projectUuid + ", and get agentsUUid: " + agentsUuid)
	status := StartTask(projectUuid, agentsUuid, config.Task1, config.TaskSuccessNumber1, false)
	if !status {
		sleepTime := config.CheckTime + 5
		Log.InfoLog.Println("First check status is False , wait " +
			strconv.Itoa(sleepTime) + "s to start the second times")
		time.Sleep(time.Duration(sleepTime) * time.Second)
		status = StartTask(projectUuid, agentsUuid, config.Task2, config.TaskSuccessNumber2, true)
	} else {
		Log.InfoLog.Println("First check status is Success")
	}
	Log.InfoLog.Println(projectUuid + ":" + agentsUuid + " check status: " + strconv.FormatBool(status) + " is finished")
	PoilceByPatrol(configfile, "true")
}

func StartTask(projectUuid string, agents string, taskNumber int, taskSuccessNumber int, isMonitor bool) bool {
	var taskId = make([]string, 0)

	// 记录任务始末时间
	taskTmpId, err := DispatchApi.StartTask(taskNumber, config.TimeOut, projectUuid, config.Shell, config.Args, config.Url)
	if err != nil {
		Log.ErrorLog.Println("StartTask Error: " + err.Error())
		PoilceByPatrol(configfile+"="+taskTmpId, "false")
		return false
	}
	taskId = append(taskId, taskTmpId)

	// 开始等待任务完成，记录任务结束时间
	Log.InfoLog.Println("Successfully start waiting for the task: " + taskId[0] +
		" to finish in about " + strconv.Itoa(config.TimeOut+waiteTime) + "s")

	taskJson := taskConfigJSON{
		Type:          config.Type,
		AgentUuid:     agents,
		Tasks:         config.Tasks,
		Success:       config.SuccessTip,
		TaskId:        taskId,
		IsMonitor:     isMonitor,
		AllNumber:     taskNumber,
		SuccessNumber: taskSuccessNumber,
		TimeOut:       config.TimeOut + waiteTime,
		CheckTime:     config.CheckTime,
	}
	js, _ := json.Marshal(&taskJson)
	status, err := CheckEsInfo.GetTaskInfo(js)
	if err != nil {
		Log.ErrorLog.Println("CheckEsInfo Error: " + err.Error())
		PoilceByPatrol(configfile+"="+taskTmpId, "false")
		return false
	}
	PoilceByPatrol(configfile+"="+taskTmpId, "true")
	return status
}

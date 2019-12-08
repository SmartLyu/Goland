package main

import (
	"./CheckEsInfo"
	"./DispatchApi"
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
	for {
		Load(configfile)
		StartProject()
		sleepTime := config.CheckTime + 5
		SuccessFulOut("wait " + strconv.Itoa(sleepTime) + "s to start next time, now project's goroutine: " + strconv.Itoa(runtime.NumGoroutine()))
		time.Sleep(time.Duration(sleepTime) * time.Second)
		fmt.Println()
	}
}

func StartProject() {
	policeStatus = "true"
	//defer PoilceByPatrol(configfile, policeStatus)

	projectUuid, agentsUuid, err := MonitorApi.DoDefault(config.Project, config.Place)
	if err != nil {
		ErrorOut("MonitorApi DoDefault Error: " + err.Error())
		policeStatus = "false"
		return
	}
	fmt.Println()
	SuccessFulOut("Get ProjectUUid: " + projectUuid + ", and get agentsUUid: " + agentsUuid)
	status := StartTask(projectUuid, agentsUuid, config.Task1, config.TaskSuccessNumber1, false)
	if !status {
		sleepTime := config.CheckTime + 5
		SuccessFulOut("First check status is False , wait " +
			strconv.Itoa(sleepTime) + "s to start the second times")
		time.Sleep(time.Duration(sleepTime) * time.Second)
		status = StartTask(projectUuid, agentsUuid, config.Task2, config.TaskSuccessNumber2, true)
	} else {
		SuccessFulOut("First check status is Success")
	}
	SuccessFulOut(projectUuid + ":" + agentsUuid + " check status: " + strconv.FormatBool(status) + " is finished")
}

func StartTask(projectUuid string, agents string, taskNumber int, taskSuccessNumber int, isMonitor bool) bool {
	policeStatus = "true"
	//defer PoilceByPatrol(configfile, policeStatus)

	// 记录任务始末时间
	taskId, err := DispatchApi.StartTask(taskNumber, config.TimeOut, projectUuid, config.Shell, config.Args)
	if err != nil {
		ErrorOut("StartTask Error: " + err.Error())
		policeStatus = "false"
		return false
	}

	// 开始等待任务完成，记录任务结束时间
	SuccessFulOut("Successfully start waiting for the task: " + taskId +
		" to finish in about " + strconv.Itoa(config.TimeOut+waiteTime) + "s")

	taskJson := taskConfigJSON{
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
		ErrorOut("CheckEsInfo Error: " + err.Error())
		policeStatus = "false"
		return false
	}
	PoilceByPatrol("task-"+configfile+"="+taskId, strconv.FormatBool(status))
	return status
}

func SuccessFulOut(msg string) {
	fmt.Printf("%c[%d;%d;%dm%s%c[0m\n", 0x1B, 1, 0, 32, time.Now().Format("2006-01-02 15:04:05 +0800 CST")+" Info "+msg, 0x1B)
}

func ErrorOut(msg string) {
	fmt.Printf("%c[%d;%d;%dm%s%c[0m\n", 0x1B, 1, 0, 31, time.Now().Format("2006-01-02 15:04:05 +0800 CST")+" ERROR "+msg, 0x1B)
}

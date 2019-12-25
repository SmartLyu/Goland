package CheckEsInfo

import (
	"../Log"
	"../MonitorApi"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"
)

// GetTaskInfo 开始获取数据
func GetTaskInfo(taskJson []byte) (bool, error) {
	//读取的数据为json格式，需要进行解码
	allInfo = make(map[string]taskInfo)
	isMonitor = make(map[string]bool)
	err := json.Unmarshal(taskJson, &Taskconfig)
	if err != nil {
		return false, errors.New("配置文件格式错误: " + err.Error())
	}
	switch Taskconfig.Type {
	case "test":
		return DoTestGetTaskInfo()
	case "monitor":
		return DoMonitorGetTaskInfo()
	}
	return false, errors.New("Type配置信息导入异常")
}

func DoTestGetTaskInfo() (bool, error) {
	var getEsInfo GetESInfoJSON
	var sleepTimes = Taskconfig.TimeOut/Taskconfig.CheckTime + 1

	for i := 1; i <= sleepTimes; i++ {
		time.Sleep(time.Duration(Taskconfig.CheckTime) * time.Second)
		Log.DebugLog.Println("Has wait " + strconv.Itoa(Taskconfig.CheckTime*i) + "s, and start " + strconv.Itoa(i) + " times to get log info")
		for _, i := range Taskconfig.TaskId {
			err := GetInfo(getEsInfo.BuildJSON(i))
			if err != nil {
				return false, errors.New("获取Es数据错误: " + err.Error())
			}
		}

		ReturnAllInfo()
	}

	var returnBool = true
	for _, i := range isMonitor {
		returnBool = returnBool && i
	}
	return returnBool, nil
}

func DoMonitorGetTaskInfo() (bool, error) {
	var getEsInfo GetESInfoJSON
	var sleepTimes = Taskconfig.TimeOut/Taskconfig.CheckTime + 1
	var status = true

	for i := 1; i <= sleepTimes; i++ {
		time.Sleep(time.Duration(Taskconfig.CheckTime) * time.Second)
		Log.DebugLog.Println("Has wait " + strconv.Itoa(Taskconfig.CheckTime*i) +
			"s, and start " + strconv.Itoa(i) + " times to get log info")
		err := GetInfo(getEsInfo.BuildJSON(Taskconfig.TaskId[0]))
		if err != nil {
			return false, errors.New("获取Es数据错误: " + err.Error())
		}
		if i == sleepTimes {
			status, err = CheckSuccessNumber(true)
		} else {
			status, err = CheckSuccessNumber(false)
		}

		if status == false {
			return false, nil
		}

		if err != nil {
			return false, errors.New("读取Es数据错误: " + err.Error())
		}

		if finishNumber == len(Taskconfig.Tasks) {
			break
		}
	}

	for _, i := range isMonitor {
		status = status && i
	}

	if !status {
		Log.ErrorLog.Println(Taskconfig.TaskId[0] + " is error")
		for i, j := range allInfo {
			var mapstr = i + ": [ "
			for _, i := range Taskconfig.Tasks {
				mapstr = mapstr + i + ":" + j.task[i] + " "
			}
			mapstr = mapstr + " ]"
			Log.ErrorLog.Println(Taskconfig.TaskId[0] + ": " + mapstr)
		}
	}
	return status, nil
}

// CheckSuccessNumber 返回状态信息，ture代表超时算不正常，false代表超时算正常
func CheckSuccessNumber(status bool) (bool, error) {
	var Successcheck = make(map[string]int)
	var Errorcheck = make(map[string]int)
	finishNumber = 0
	for _, i := range Taskconfig.Tasks {
		Successcheck[i] = 0
		Errorcheck[i] = 0
	}
	for i, j := range allInfo {
		var mapstr = i + ": [ "
		for _, i := range Taskconfig.Tasks {
			mapstr = mapstr + i + ":" + j.task[i] + " "
		}
		mapstr = mapstr + " ]"
		for k, v := range j.task {
			if !strings.Contains(mapstr, k) {
				mapstr = mapstr + " ? " + k + ":" + v
				continue
			}
			if v == Taskconfig.Success {
				Successcheck[k]++
				continue
			}
			if v != emptyString {
				Errorcheck[k]++
				continue
			}
		}
		Log.DebugLog.Println(Taskconfig.TaskId[0] + ": " + mapstr)
	}

	for _, i := range Taskconfig.Tasks {
		_, ok := isMonitor[i]
		if ok {
			finishNumber++
			continue
		}
		if Successcheck[i] >= Taskconfig.SuccessNumber {
			err := MonitorApi.Monitor(Taskconfig.AgentUuid, i, 0)
			if err != nil {
				return true, err
			}
			isMonitor[i] = true
			finishNumber++
			continue
		}
		if Errorcheck[i] > Taskconfig.AllNumber-Taskconfig.SuccessNumber {
			if Taskconfig.IsMonitor {
				err := MonitorApi.Monitor(Taskconfig.AgentUuid, i, 1)
				if err != nil {
					return true, err
				}
				return false, nil
			}
			isMonitor[i] = false
			finishNumber++
			continue
		}
		if status && Successcheck[i] < Taskconfig.SuccessNumber {
			if Taskconfig.IsMonitor {
				err := MonitorApi.Monitor(Taskconfig.AgentUuid, i, 1)
				if err != nil {
					return true, err
				}
				return false, nil
			}
			isMonitor[i] = false
			finishNumber++
			continue
		}
	}
	return true, nil
}

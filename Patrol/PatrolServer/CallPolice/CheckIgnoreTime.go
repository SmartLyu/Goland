package CallPolice

import (
	"../Global"
	"encoding/json"
	"io/ioutil"
	"regexp"
	"time"
)

// 判断是否属于日常维护时间
func CheckIgnoreTime(monitorjson Global.MonitorJson) (status bool) {
	var jsons []Global.IgnoreTimeJson
	status = false

	data, err := ioutil.ReadFile(Global.IgnoreTimeFile)
	if err != nil {
		CallMessage(Global.IgnoreTimeFile + " json file can not read , error is " + err.Error())
		return
	}

	if err := json.Unmarshal(data, &jsons); err != nil {
		CallMessage(Global.IgnoreTimeFile + " json is error , error is " + err.Error())
		return
	}

	for _, key := range jsons {
		if check, _ := regexp.Match(key.IP, []byte(monitorjson.IP)); check {
			if check, _ := regexp.Match(key.Hostname, []byte(monitorjson.Hostname)); check {
				if check, _ := regexp.Match(key.Info, []byte(monitorjson.Info)); check {
					startTime, err := time.Parse("15:04:05", key.StartTime)
					if err != nil {
						CallMessage(Global.IgnoreTimeFile + " start time is error , error is " + err.Error())
						return
					}
					endTime, err := time.Parse("15:04:05", key.EndTime)
					if err != nil {
						CallMessage(Global.IgnoreTimeFile + " end time is error , error is " + err.Error())
						return
					}
					nowTime, _ := time.Parse("15:04:05", time.Now().Format("15:04:05"))
					if nowTime.After(startTime) && nowTime.Before(endTime) {
						status = true
					}
				}
			}
		}
	}
	return
}

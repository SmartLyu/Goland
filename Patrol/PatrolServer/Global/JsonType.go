package Global

import (
	"encoding/json"
	"log"
	"sort"
	"time"
)

type MonitorJson struct {
	Time     string `json:"Time"`
	IP       string `json:"IP"`
	Hostname string `json:"Hostname"`
	Info     string `json:"Info"`
	Status   bool   `json:"Status"`
}

type NatTable struct {
	IP       string `json:"IP"`
	HostName string `json:"HOSTNAME"`
	Port     int    `json:"PORT"`
	Time     int    `json:"TIME"`
}

type ErrorJson struct {
	Key   string `json:"KEY"`
	Value int    `json:"VALUE"`
}

type HostsTable struct {
	IP       string `json:"IP"`
	HostName string `json:"hostname"`
}

type DateTimeStyle struct {
	Year   string
	Month  string
	Day    string
	Hour   string
	Minute string
}

type DingdingAtJson struct {
	Hostname string `json:"hostname"`
	Mobiles  []int  `json:"mobiles"`
}

func ReadJson(mj MonitorJson) string {
	js, _ := json.Marshal(&mj)
	return string(js)
}

func (m *MonitorJson) Exist() bool {

	if m.IP == "" || m.Info == "" || m.Time == "" {
		return false
	}
	return true
}

// sort 需要的接口配置
type NameSorter []MonitorJson

func (a NameSorter) Len() int      { return len(a) }
func (a NameSorter) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a NameSorter) Less(i, j int) bool {
	if a[i].IP < a[j].IP {
		return true
	} else if a[i].IP > a[j].IP {
		return false
	} else if a[i].Time == a[j].Time {
		if a[i].Info < a[j].Info {
			return true
		} else {
			return false
		}
	} else if a[i].Time < a[j].Time {
		return true
	} else {
		return false
	}

}

// 解析json格式的string信息
func ReadMonitorJson(body string) []MonitorJson {
	var jsonfile []MonitorJson
	if err := json.Unmarshal([]byte(body), &jsonfile); err != nil {
		log.Fatal(err.Error())
	}
	sort.Sort(NameSorter(jsonfile))
	return jsonfile
}

// 时间范围获取各分钟信息
func GetAllTime(getTimeStart int64, getTimeEnd int64) (getTime []DateTimeStyle) {
	var tmpDate DateTimeStyle

	for i := getTimeStart; i <= getTimeEnd; i += 60 {
		tmpDate = DateTimeStyle{
			time.Unix(i, 0).Format("2006"),
			time.Unix(i, 0).Format("01"),
			time.Unix(i, 0).Format("02"),
			time.Unix(i, 0).Format("15"),
			time.Unix(i, 0).Format("04"),
		}
		getTime = append(getTime, tmpDate)
	}
	return
}

func (m *DateTimeStyle) Exist() bool {
	if m.Month == "" || m.Minute == "" || m.Day == "" || m.Year == "" || m.Hour == "" {
		return false
	}
	return true
}

func (m *DateTimeStyle) Print() string {
	if m.Exist() {
		return m.Year + "." + m.Month + "." + m.Day + "-" + m.Hour + ":" + m.Minute
	}
	return ""
}

func RemoveRepeatedElement(arr []string) (newArr []string) {
	newArr = make([]string, 0)
	sort.Strings(arr)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
}

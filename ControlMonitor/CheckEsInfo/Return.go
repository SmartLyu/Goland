package CheckEsInfo

import (
	"../Log"
	"sort"
	"strconv"
	"strings"
)

// sort 需要的接口配置
type NameSorter []taskInfo

func (a NameSorter) Len() int      { return len(a) }
func (a NameSorter) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a NameSorter) Less(i, j int) bool {
	if a[i].serverip < a[j].serverip {
		return true
	} else if a[i].serverip > a[j].serverip {
		return false
	} else {
		if a[i].serverport < a[j].serverport {
			return true
		} else {
			return false
		}
	}
}

func ReturnAllInfo() {
	finishNumber = 0
	var newMp = make([]taskInfo, 0)
	for _, v := range allInfo {
		newMp = append(newMp, v)
	}
	sort.Sort(NameSorter(newMp))
	var ip = make([]string, 0)
	var totalInfo = make(map[string]totalMap)
	for _, v := range newMp {
		Log.DebugLog.Println(v.serverip + ":" + v.serverport + " - " + ReturnStatusMap(v))
		var isAdd = true
		if len(ip) != 0 {
			for _, i := range ip {
				if i == v.serverip {
					isAdd = false
					continue
				}
			}
		}
		if isAdd {
			ip = append(ip, v.serverip)
			totalInfo[v.serverip] = totalMap{
				make(map[string]int),
				make(map[string]int),
				make(map[string]int),
			}
		}
		for i, j := range v.task {
			if j == Taskconfig.Success {
				SuccessAdd(totalInfo[v.serverip], i)

			} else if j == emptyString {
				NilAdd(totalInfo[v.serverip], i)
			} else {
				Log.ErrorLog.Println(v.serverip + ":" + v.serverport + " - " + ReturnStatusMap(v))
				FalseAdd(totalInfo[v.serverip], i)
			}
		}
	}

	for _, i := range Taskconfig.Tasks {
		var all = make(map[string]int)
		for _, j := range ip {
			all["success"] += totalInfo[j].Success[i]
			all["false"] += totalInfo[j].False[i]
			all["nil"] += totalInfo[j].Nil[i]
		}
		Log.InfoLog.Println(i + ": [ Success: " + strconv.Itoa(all["success"]) + ", False: " +
			strconv.Itoa(all["false"]) + ", Nil: " + strconv.Itoa(all["nil"]) + " ]")
	}
}

func ReturnStatusMap(t taskInfo) (mapstr string) {
	mapstr = "[ "
	for _, i := range Taskconfig.Tasks {
		mapstr = mapstr + i + ":" + t.task[i] + " "
	}
	mapstr = mapstr + " ]"
	for i, j := range t.task {
		if !strings.Contains(mapstr, i) {
			mapstr = mapstr + " ? " + i + ":" + j
		}
	}
	return
}

type totalMap struct {
	Success map[string]int
	False   map[string]int
	Nil     map[string]int
}

func SuccessAdd(m totalMap, s string) {
	_, ok := m.Success[s]
	if ok {
		m.Success[s]++
		return
	}
	m.Success[s] = 1
}

func FalseAdd(m totalMap, s string) {
	_, ok := m.False[s]
	if ok {
		m.False[s]++
		return
	}
	m.False[s] = 1

}

func NilAdd(m totalMap, s string) {
	_, ok := m.Nil[s]
	if ok {
		m.Nil[s]++
		return
	}
	m.Nil[s] = 1
}

package CheckEsInfo

import "encoding/json"

var (
	// Taskconfig task数据查询配置信息
	Taskconfig  taskConfigJSON
	allInfo     map[string]taskInfo
	isMonitor   map[string]bool
	emptyString = "nil"

	// 记录已经获取完整的task个数
	finishNumber int
)

type taskConfigJSON struct {
	Type          string   `json:"type"`
	AgentUuid     string   `json:"agentUuid"`
	Tasks         []string `json:"tasks"`
	Success       string   `json:"success"`
	TaskId        []string `json:"taskId"`
	AllNumber     int      `json:"allNumber"`
	IsMonitor     bool     `json:"isMonitor"`
	SuccessNumber int      `json:"successNumber"`
	TimeOut       int      `json:"timeout"`
	CheckTime     int      `json:"checkTime"`
}

type taskInfo struct {
	serverip   string
	serverport string
	task       map[string]string
}

func (m *taskInfo) SetIp(ip string) {
	m.serverip = ip
}

func (m *taskInfo) SetPort(port string) {
	m.serverport = port
}

func InitTask(ip string, port string) {
	var info = taskInfo{
		ip,
		port,
		make(map[string]string),
	}

	for _, i := range Taskconfig.Tasks {
		info.task[i] = emptyString
	}

	allInfo[ip+":"+port] = info
}

// Settask 将es查到的日志信息存入数组
func SetTask(ipAndPort string, task string, status string) {
	allInfo[ipAndPort].task[task] = status
	if status != Taskconfig.Success {
		var getLocation = false
		for _, i := range Taskconfig.Tasks {
			if getLocation {
				allInfo[ipAndPort].task[i] = status
			}
			if i == task {
				getLocation = true
			}
		}
	}
}

// GetESInfoJSON 配置ES查询body数据
type GetESInfoJSON struct {
	Query GetESInfoJSONQuery `json:"query"`
	Size  int                `json:"size"`
}

// GetESInfoJSONQuery 配置ES查询body数据
type GetESInfoJSONQuery struct {
	Bool GetESInfoJSONBool `json:"bool"`
}

// GetESInfoJSONBool 配置ES查询body数据
type GetESInfoJSONBool struct {
	Must GetESInfoJSONMust `json:"must"`
}

// GetESInfoJSONMust 配置ES查询body数据
type GetESInfoJSONMust struct {
	Match GetESInfoJSONMatch `json:"match"`
}

// GetESInfoJSONMatch 配置ES查询body数据
type GetESInfoJSONMatch struct {
	Taskip string `json:"taskid"`
}

// BuildJSON 配置ES查询body数据
func (m *GetESInfoJSON) BuildJSON(taskId string) (js []byte) {
	m.Query.Bool.Must.Match.Taskip = taskId
	m.Size = 10000

	js, _ = json.Marshal(m)
	return
}

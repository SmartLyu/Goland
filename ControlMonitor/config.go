package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

var (
	config     configJSON
	configfile string
	// 调度默认启动延迟时间
	waiteTime    = 30
	policeStatus string
)

type configJSON struct {
	// Place 国内、海外
	Place              string   `json:"place"`
	Tasks              []string `json:"tasks"`
	SuccessTip         string   `json:"successTip"`
	Project            string   `json:"project"`
	Shell              string   `json:"shell"`
	Args               string   `json:"args"`
	Task1              int      `json:"task1"`
	TaskSuccessNumber1 int      `json:"taskSuccessNumber1"`
	Task2              int      `json:"task2"`
	TaskSuccessNumber2 int      `json:"taskSuccessNumber2"`
	TimeOut            int      `json:"timeOut"`
	CheckTime          int      `json:"checkTime"`
}

type taskConfigJSON struct {
	AgentUuid     string   `json:"agentUuid"`
	Tasks         []string `json:"tasks"`
	Success       string   `json:"success"`
	TaskId        string   `json:"taskId"`
	AllNumber     int      `json:"allNumber"`
	SuccessNumber int      `json:"successNumber"`
	TimeOut       int      `json:"timeout"`
	IsMonitor     bool     `json:"isMonitor"`
	CheckTime     int      `json:"checkTime"`
}

// Load 读配置文件
func Load(filename string) {
	//ReadFile函数会读取文件的全部内容，并将结果以[]byte类型返回
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("配置文件不存在: " + err.Error())
	}

	//读取的数据为json格式，需要进行解码
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatal("配置文件格式错误: " + err.Error())
	}

	if config.Place == "" {
		log.Fatal("配置文件格式错误: Place is empty")
	}
	if config.SuccessTip == "" {
		log.Fatal("配置文件格式错误: SuccessTip is empty")
	}
	if config.Project == "" {
		log.Fatal("配置文件格式错误: Project is empty")
	}
	if config.Shell == "" {
		log.Fatal("配置文件格式错误: Shell is empty")
	}
	if config.Args == "" {
		log.Fatal("配置文件格式错误: Args is empty")
	}
	if config.Task1 == 0 {
		log.Fatal("配置文件格式错误: Task1 is empty")
	}
	if config.Task2 == 0 {
		log.Fatal("配置文件格式错误: Task2 is empty")
	}
	if config.TaskSuccessNumber1 == 0 {
		log.Fatal("配置文件格式错误: TaskSuccessNumber1 is empty")
	}
	if config.TaskSuccessNumber2 == 0 {
		log.Fatal("配置文件格式错误: TaskSuccessNumber2 is empty")
	}
	if config.TimeOut == 0 {
		log.Fatal("配置文件格式错误: TimeOut is empty")
	}
	if config.CheckTime == 0 {
		log.Fatal("配置文件格式错误: CheckTime is empty")
	}
	if len(config.Tasks) == 0 {
		log.Fatal("配置文件格式错误: Tasks is empty")
	}
}

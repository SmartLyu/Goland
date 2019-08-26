package Global

import "encoding/json"

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
	Value int `json:"VALUE"`
}

type HostsTable struct {
	IP   string `json:"IP"`
	Time string `json:"TIME"`
}

func ReadJson(mj MonitorJson) (string) {
	js, _ := json.Marshal(&mj)
	return string(js)
}

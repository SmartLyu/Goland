package Api

import "encoding/json"

type MonitorJson struct {
	IP       string `json:"IP"`
	Hostname string `json:"Hostname"`
	Info     string `json:"Info"`
	Status   bool   `json:"Status"`
}

func ReadJson(mj MonitorJson) (string) {
	js, _ := json.Marshal(&mj)
	return string(js)
}

var Jsons []MonitorJson

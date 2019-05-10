package Api

type MonitorJson struct {
	IP       string `json:"IP"`
	Hostname string `json:"Hostname"`
	Info     string `json:"Info"`
	Status   bool   `json:"Status"`
}

var Jsons []MonitorJson

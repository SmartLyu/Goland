package main

import "time"

var (
	MonitorUrl          = "http://10.4.0.17:8686/monitor/collect"
	LogDir              = "/work/logs/"
	LogFile             = "patrol.log"
	ApiPost             = "8666"
	SshUser             = "work"
	Sshkey              = "/work/apps/.secret"
	NatShellDownloadUrl = "134.175.50.184:8686/shell/nat"
	listenSig           = make(chan int)
)

func PostJson(message string, status string) string {
	jsonstr := "{" +
		"\"time\": \"" + time.Now().Format("2006-01-02 15:04") + "\"," +
		"\"IP\": \"123.207.233.139-JH-Bak-QCloudGZ-Nat=}10.4.0.4\"," +
		"\"hostname\": \"JH-Api-QCloudGZ3-Jumpserver\"," +
		"\"info\": \"" + message + "\"," +
		"\"status\": " + status + " }"
	return jsonstr
}

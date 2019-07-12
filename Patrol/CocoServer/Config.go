package main

var (
	MonitorUrl              = "134.175.50.184:8666/monitor/collect"
	LogDir                  = "/work/logs/"
	LogFile                 = "patrol.log"
	ApiPost                 = "8666"
	SshUser                 = "work"
	Sshkey                  = "/work/apps/.secret"
	NatShellDownloadUrl     = "134.175.50.184:8686/shell/nat"
)

func PostJson(message string, status string) string {
	jsonstr := "{" +
		"\"IP\":  \"193.112.24.232\"" +
		"\"hostname\": \"JH-Api-QCloudGZ3-Jumpserver\"" +
		"\"info\": \"" + message + "\"" +
		"\"staus\": " + status + ""
	return jsonstr
}

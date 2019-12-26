package CallCoco

import (
	"../CallPolice"
	"../File"
	"../Global"
	"../Mysql"
	"strconv"
	"time"
)

func CallCoco(hostname string, ip string, port string) {
	err := httpPostJson(ip, port)
	if err == nil {
		jsonfile := Global.MonitorJson{
			Time:     time.Now().Format("2006-01-02 15:04"),
			IP:       "123.207.233.139-JH-Bak-QCloudGZ3-Nat=}10.4.0.17",
			Hostname: "PatrolMessage",
			Info:     "callCoco-" + hostname + "-" + ip,
			Status:   true,
		}
		CallPolice.Judge(jsonfile)
		if err := File.WriteFile(Global.ReadJson(jsonfile), jsonfile.Time); err != nil {
			Global.ErrorLog.Println("write info " + err.Error())
		}
		return
	} else {
		err := httpPostJson(ip, port)
		if err != nil {
			jsonfile := Global.MonitorJson{
				Time:     time.Now().Format("2006-01-02 15:04"),
				IP:       "123.207.233.139-JH-Bak-QCloudGZ3-Nat=}10.4.0.17",
				Hostname: "PatrolMessage",
				Info:     "callCoco-" + hostname + "-" + ip,
				Status:   false,
			}
			CallPolice.Judge(jsonfile)
			if err := File.WriteFile(Global.ReadJson(jsonfile), jsonfile.Time); err != nil {
				Global.ErrorLog.Println("write info " + err.Error())
			}
		} else {
			jsonfile := Global.MonitorJson{
				Time:     time.Now().Format("2006-01-02 15:04"),
				IP:       "123.207.233.139-JH-Bak-QCloudGZ3-Nat=}10.4.0.17",
				Hostname: "PatrolMessage",
				Info:     "callCoco-" + hostname + "-" + ip,
				Status:   true,
			}
			CallPolice.Judge(jsonfile)
			if err := File.WriteFile(Global.ReadJson(jsonfile), jsonfile.Time); err != nil {
				Global.ErrorLog.Println("write info " + err.Error())
			}
		}
	}
}

func CallAllNatMonitor() {
	for _, i := range Mysql.SelectAllNatTable() {
		CallCoco(i.HostName, i.IP, strconv.Itoa(i.Port))
	}
}

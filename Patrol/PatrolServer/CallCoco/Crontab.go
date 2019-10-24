package CallCoco

import (
	"../CallPolice"
	"../File"
	"../Global"
	"../Mysql"
	"github.com/robfig/cron"
	"strconv"
	"time"
)

var cs = make(map[string]*cron.Cron)

// 计划任务
func CrontabToCallCoco(nt Global.NatTable) {
	cs[nt.IP] = cron.New()
	spec := "0 */" + strconv.Itoa(nt.Time) + " * * * ?"

	err := cs[nt.IP].AddFunc(spec, func() {
		CallCoco(nt.HostName, nt.IP, strconv.Itoa(nt.Port))
	})

	if err != nil {
		File.WriteErrorLog("crontab is error: " + err.Error())
	}
	cs[nt.IP].Start()
}

// 计划任务
func CrontabToDelMap() {
	spec := "0 0 10 * * ?"
	c := cron.New()

	err := c.AddFunc(spec, func() {
		for key, _ := range Global.ErrorMap.Data {
			File.WriteInfoLog("delete error map: " + key)
			Global.ErrorMap.Delete(key)
		}
	})

	if err != nil {
		File.WriteErrorLog("crontab is error: " + err.Error())
	}
	c.Start()
}

// 开始所有库中nat机器的计划任务
func StartAllCrontab() {
	for _, i := range Mysql.SelectAllNatTable() {
		go CrontabToCallCoco(i)
	}
}

// 重新读取数据库中nat机器的计划任务
func ReStartAllCrontab() {
	for _, i := range cs {
		i.Stop()
	}
	for _, i := range Mysql.SelectAllNatTable() {
		go CrontabToCallCoco(i)
	}
}

func StopCrontab(nt Global.NatTable) {
	cs[nt.IP].Stop()
}

func CrontabToCheckHosts() {
	c := cron.New()
	spec := "50 * * * * ?"

	err := c.AddFunc(spec, func() {
		ht := Mysql.SelectHostsTable()
		for _, i := range ht {
			Mysql.DeleteHosts(i)
			pwd := Global.DataFileDir
			des := pwd + time.Now().Format("2006-01/02") + Global.DataFileName

			_, err := File.FindWorkInFile(des, time.Now().Add(-time.Minute * 1).Format("2006-01-02 15:04"),
				i.IP, "true")

			if err == nil {
				continue
			}

			_, err = File.FindWorkInFile(des, time.Now().Add(-time.Minute * 2).Format("2006-01-02 15:04"),
				i.IP, "true")

			if err == nil {
				continue
			}

			jsonfile := Global.MonitorJson{
				Time: time.Now().Format("2006-01-02 15:04"),
				IP: i.IP,
				Hostname: "Unknown-Hostname",
				Info: "survive",
				Status: false,
			}

			CallPolice.Judge(jsonfile)

			// 添加数据
			if err := File.WriteFile(Global.ReadJson(jsonfile)); err != nil {
				File.WriteErrorLog(err.Error())
			}
		}
	})

	if err != nil {
		File.WriteErrorLog("crontab is error: " + err.Error())
	}
	c.Start()
}

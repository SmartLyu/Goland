package CallCoco

import (
	"../CallPolice"
	"../File"
	"../Global"
	"../Mysql"
	"github.com/robfig/cron"
	"strconv"
	"sync"
	"time"
)

var cs = make(map[string]*cron.Cron)
var csLock sync.Mutex

// 计划任务
func CrontabToCallCoco(nt Global.NatTable) {
	csLock.Lock()
	defer csLock.Unlock()
	cs[nt.IP] = cron.New()
	spec := "0 */" + strconv.Itoa(nt.Time) + " * * * ?"

	err := cs[nt.IP].AddFunc(spec, func() {
		CallCoco(nt.HostName, nt.IP, strconv.Itoa(nt.Port))
	})

	if err != nil {
		Global.ErrorLog.Println("crontab is error: " + err.Error())
	}
	cs[nt.IP].Start()
}

// 计划任务
func CrontabToDelMap() {
	spec := "0 0 10 * * ?"
	c := cron.New()

	err := c.AddFunc(spec, func() {
		for key, _ := range Global.ErrorMap.Data {
			Global.InfoLog.Println("delete error map: " + key)
			Global.ErrorMap.Delete(key)
		}
	})

	if err != nil {
		Global.ErrorLog.Println("crontab is error: " + err.Error())
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
	csLock.Lock()
	defer csLock.Unlock()
	for _, i := range cs {
		i.Stop()
	}
	for _, i := range Mysql.SelectAllNatTable() {
		go CrontabToCallCoco(i)
	}
}

func StopCrontab(nt Global.NatTable) {
	csLock.Lock()
	defer csLock.Unlock()
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

			des := pwd + time.Now().Format("2006-01/02/15/04") + Global.DataFileName
			_, err := File.FindWorkInFile(des, i.IP, "survive", "true")
			if err == nil {
				continue
			}

			jsonfile := Global.MonitorJson{
				Time:     time.Now().Format("2006-01-02 15:04"),
				IP:       i.IP,
				Hostname: "Unknown-Hostname",
				Info:     "survive",
				Status:   false,
			}

			// 添加数据
			if err := File.WriteFile(Global.ReadJson(jsonfile), jsonfile.Time); err != nil {
				Global.ErrorLog.Println(err.Error())
			}

			des = pwd + time.Now().Add(-time.Minute*3).Format("2006-01/02/15/04") + Global.DataFileName
			_, err = File.FindWorkInFile(des, i.IP, "survive", "true")
			if err == nil {
				continue
			}

			des = pwd + time.Now().Add(-time.Minute*2).Format("2006-01/02/15/04") + Global.DataFileName
			_, err = File.FindWorkInFile(des, i.IP, "survive", "true")
			if err == nil {
				continue
			}

			des = pwd + time.Now().Add(-time.Minute*1).Format("2006-01/02/15/04") + Global.DataFileName
			_, err = File.FindWorkInFile(des, i.IP, "survive", "true")
			if err == nil {
				continue
			}

			CallPolice.Judge(jsonfile)
			if err := File.WriteFile(Global.ReadJson(jsonfile), jsonfile.Time); err != nil {
				Global.ErrorLog.Println("write info " + err.Error())
			}
		}

		// 查看是否有异常日志产生
		des := Global.UpdateLog(Global.ErrorFileName)
		message, err := File.FindWorkInFile(des, time.Now().Add(-1*time.Minute).Format("2006/01/02 15:04"))
		if err == nil {
			CallPolice.CallMessage(message)
		}
	})

	if err != nil {
		Global.ErrorLog.Println("crontab is error: " + err.Error())
	}
	c.Start()
}

func CrontabToCutLog() {
	c := cron.New()
	spec := "0 0 0 * * ?"

	err := c.AddFunc(spec, func() {
		Global.CutLog()
	})

	if err != nil {
		Global.ErrorLog.Println("crontab is error: " + err.Error())
	}
	c.Start()
}

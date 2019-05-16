package CallCoco

import (
	"../CallPolice"
	"../File"
	"../Global"
	"../Mysql"
	"github.com/robfig/cron"
	"strconv"
)

var cs = make(map[string]*cron.Cron)

func CrontabToCallCoco(nt Global.NatTable) {
	cs[nt.IP] = cron.New()
	spec := "0 */" + strconv.Itoa(nt.Time) + " * * * ?"

	err := cs[nt.IP].AddFunc(spec, func() {
		CallAllNatMonitor()
	})

	if err != nil {
		File.WriteErrorLog("crontab is error: " + err.Error())
	}
	cs[nt.IP].Start()
	select {}
}

func StartAllCrontab() {
	for _, i := range Mysql.SelectAllNatTable() {
		go CrontabToCallCoco(i)
	}
}

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
			CallPolice.CallPolice(i.IP + "\n has not return status")
		}
	})

	if err != nil {
		File.WriteErrorLog("crontab is error: " + err.Error())
	}
	c.Start()
	select {}

}

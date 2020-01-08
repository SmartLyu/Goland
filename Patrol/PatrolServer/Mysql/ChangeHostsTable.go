package Mysql

import (
	"../Global"
)

func InsertHosts(ht Global.HostsTable) bool {
	Global.NatHostsMap.Change(ht, true)
	Global.InfoLog.Println("insert " + ht.IP + " - " + ht.HostName + " successfully")
	return true
}

func DeleteHosts(ht Global.HostsTable) bool {
	if Global.NatHostsMap.Exist(ht) {
		Global.NatHostsMap.Delete(ht)
		Global.InfoLog.Println("delete " + ht.IP + " successfully")
	}
	return true
}

func SelectHostsTable() []Global.HostsTable {
	var hts []Global.HostsTable
	var ht Global.HostsTable
	for key, value := range Global.NatHostsMap.Data {
		ht.HostName = value.Hostname
		ht.IP = key
		hts = append(hts, ht)
	}
	return hts
}

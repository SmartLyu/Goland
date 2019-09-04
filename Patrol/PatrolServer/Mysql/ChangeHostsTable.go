package Mysql

import (
	"../File"
	"../Global"
)

func InsertHosts(ht Global.HostsTable) (bool) {
	Global.NatHostsMap[ht]=true
	File.WriteInfoLog("insert " + ht.IP + " - " + ht.Time + " successfully")
	return true
}

func DeleteHosts(ht Global.HostsTable) (bool) {
	if _, isError := Global.NatHostsMap[ht]; isError {
		delete(Global.NatHostsMap,ht)
		File.WriteInfoLog("delete " + ht.IP + " successfully")
	}
	return true
}

func SelectHostsTable() ([]Global.HostsTable) {
	var hts []Global.HostsTable
	for key, _ := range Global.NatHostsMap {
		hts = append(hts, key)
	}
	return hts
}
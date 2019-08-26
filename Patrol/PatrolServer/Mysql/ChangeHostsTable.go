package Mysql

import (
	"../File"
	"../Global"
)

func InsertHosts(ht Global.HostsTable) (bool) {
	//开启事务
	tx, err := DB.Begin()
	if err != nil {
		File.WriteErrorLog("Mysql DB Hosts insert start fail")
		return false
	}
	//准备sql语句
	stmt, err := tx.Prepare("INSERT INTO nat_hosts (`IP`, `time`) VALUES (?, ?)")
	if err != nil {
		File.WriteErrorLog("Mysql DB sql insert Prepare fail")
		return false
	}
	//将参数传递到sql语句中并且执行
	_, err = stmt.Exec(ht.IP, ht.Time)
	if err != nil {
		File.WriteErrorLog("Mysql DB insert Exec fail")
		return false
	}
	//将事务提交
	_ = tx.Commit()
	//获得上一个插入自增的id
	File.WriteInfoLog("insert " + ht.IP + " - " + ht.Time + " successfully")
	return true
}

func DeleteHosts(ht Global.HostsTable) (bool) {
	//开启事务
	tx, err := DB.Begin()
	if err != nil {
		File.WriteErrorLog("Mysql DB Hosts delete start fail")
	}
	//准备sql语句
	_ , err = tx.Exec("DELETE FROM patrol.nat_hosts WHERE IP = ?",ht.IP)
	if err != nil {
		File.WriteErrorLog("Mysql DB Nat delete Exec fail")
		return false
	}

	//提交事务
	_ = tx.Commit()
	//获得上一个insert的id
	File.WriteInfoLog("delete " + ht.IP + " successfully")
	return true
}

func SelectHostsTable() ([]Global.HostsTable) {
	//执行查询语句
	rows, err := DB.Query("SELECT * from patrol.nat_hosts")
	if err != nil {
		File.WriteErrorLog("查询出错了")
	}
	var hts []Global.HostsTable
	//循环读取结果
	for rows.Next() {
		var ht Global.HostsTable
		//将每一行的结果都赋值到一个user对象中
		err := rows.Scan(&ht.IP, &ht.Time)
		if err != nil {
			File.WriteErrorLog("Mysql DB Nat select rows fail")
		}
		//将user追加到users的这个数组中
		hts = append(hts, ht)
	}
	return hts
}

package Mysql

import (
	"../Global"
	"strconv"
)

func InsertNat(nt Global.NatTable) bool {
	//开启事务
	tx, err := DB.Begin()
	if err != nil {
		Global.ErrorLog.Println("Mysql DB Nat insert start fail" + err.Error())
		return false
	}
	//准备sql语句
	stmt, err := tx.Prepare("INSERT INTO nat_info (`IP`,`hostname`,`port`, `time`) VALUES (?, ?, ?, ?)")
	if err != nil {
		Global.ErrorLog.Println("Mysql DB Nat insert Prepare fail: " + err.Error())
		return false
	}
	//将参数传递到sql语句中并且执行
	_, err = stmt.Exec(nt.IP, nt.HostName, nt.Port, nt.Time)
	if err != nil {
		Global.ErrorLog.Println("Mysql DB Nat insert Exec fail: " + err.Error())
		return false
	}
	//将事务提交
	_ = tx.Commit()
	//获得上一个插入自增的id
	Global.InfoLog.Println("insert " + nt.IP + " - " + strconv.Itoa(nt.Port) +
		" - " + strconv.Itoa(nt.Time) + " successfully")
	return true
}

func DeleteNat(nt Global.NatTable) bool {
	//开启事务
	tx, err := DB.Begin()
	if err != nil {
		Global.ErrorLog.Println("Mysql DB Nat delete start  fail: " + err.Error())
		return false
	}
	//准备sql语句
	stmt, err := tx.Prepare("DELETE FROM nat_info WHERE ip = ?")
	if err != nil {
		Global.ErrorLog.Println("Mysql DB Nat delete Prepare fail: " + err.Error())
		return false
	}
	//设置参数以及执行sql语句
	_, err = stmt.Exec(nt.IP)
	if err != nil {
		Global.ErrorLog.Println("Mysql DB Nat delete Exec fail: " + err.Error())
		return false
	}
	//提交事务
	_ = tx.Commit()
	//获得上一个insert的id
	Global.InfoLog.Println("delete " + nt.IP + " - " + strconv.Itoa(nt.Port) +
		" - " + strconv.Itoa(nt.Time) + " successfully")
	return true
}

func SelectAllNatTable() []Global.NatTable {
	//执行查询语句
	rows, err := DB.Query("SELECT * from patrol.nat_info")
	if err != nil {
		Global.ErrorLog.Println("查询出错了")
	}
	var nts []Global.NatTable
	//循环读取结果
	for rows.Next() {
		var nt Global.NatTable
		//将每一行的结果都赋值到一个user对象中
		err := rows.Scan(&nt.IP, &nt.HostName, &nt.Port, &nt.Time)
		if err != nil {
			Global.ErrorLog.Println("Mysql DB Nat select rows fail: " + err.Error())
		}
		//将user追加到users的这个数组中
		nts = append(nts, nt)
	}
	return nts
}

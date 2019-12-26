package Mysql

import (
	"../Global"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

//数据库配置
const (
	userName = "lyu"
	password = "Patrol@123"
	ip       = "127.0.0.1"
	port     = "3306"
	dbName   = "patrol"
)

//Db数据库连接池
var DB *sql.DB

// 连接数据库
func InitDB() {
	//构建连接："用户名:密码@tcp(IP:端口)/数据库?charset=utf8"
	path := strings.Join([]string{userName, ":", password, "@tcp(", ip, ":", port, ")/", dbName, "?charset=utf8"}, "")

	//打开数据库,前者是驱动名，所以要导入： _ "github.com/go-sql-driver/mysql"
	DB, _ = sql.Open("mysql", path)
	//设置数据库最大连接数
	DB.SetConnMaxLifetime(50000)
	//设置上数据库最大闲置连接数
	DB.SetMaxIdleConns(10000)
	//验证连接
	if err := DB.Ping(); err != nil {
		Global.ErrorLog.Println("opon database fail")
		return
	}
	Global.InfoLog.Println("connnect success")
}

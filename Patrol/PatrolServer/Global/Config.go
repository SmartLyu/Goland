package Global

import (
	"log"
	"os"
	"sync"
	"time"
)

var (
	ApiPost             = "8666"                                   // 程序端口
	ApiPublicPost       = "8686"                                   // 对外程序端口
	DataFileDir         = "/work/data/patrol/"                     // 存放历史监控信息
	DataFileName        = ".monitor.log"                           // 监控后缀名
	LogFileDir          = "/work/logs/patrol/"                     // 日志存放目录(提前准备好)
	LogFileName         = ".patrol.log"                            // 日志后缀名
	ErrorFileName       = ".error.log"                             // 异常日志后缀名
	AcessLogFileName    = ".access.log"                            // 记录日志信息
	PoliceLogFileName   = ".police.log"                            // 报警日志信息
	MonitorShellFile    = "/work/sh/PatrolMonitor.sh"              // 巡查脚本存放位置
	NatShellFile        = "/work/sh/NatPatrol.sh"                  // Nat使用的巡查脚本存放位置
	DingdingDefaultAt   = "运维"                                     // 钉钉报警通知人员默认报警对象
	DingdingAtFile      = "/work/apps/patrol/config/dingding.json" // 存储钉钉报警通知人员匹配配置
	DingdingMobilesFile = "/work/apps/patrol/config/mobiles.json"  // 存储钉钉报警中人员电话的配置文件
	IgnoreTimeFile      = "/work/apps/patrol/config/ignore.json"   // 存储监控不关心时间段的服务器信息

	ErrorMax                  = 2                              // 最多报警次数
	MaxSearchLen        int64 = 500                            // 搜索文件最大次数
	MaxReturnLen        int64 = 1000000                        // 查询预估临界值
	MaxSearchSigLen           = make(chan int, 30)             // 查询并发线程数
	ListenSig                 = make(chan int)                 // 监听后台阻塞信号
	ListenPublicSig           = make(chan int)                 // 监听后台公共端口阻塞信号
	CocoUrl                   = "http://10.4.0.4:8666/monitor" // coco的端口
	ErrorMap                  = NewErrorMapType()              // 存储报警信息至内存
	NatHostsMap               = NewNatHostsMapType()           // 存储Nat机器中子服务器信息至内存
	PatrolMessageString       = "PatrolMessage"                // 报警为系统日志的hostname
	IsPolice                  = true                           // 是否报警

	PoliceLock    sync.Mutex // 报警发锁
	FileWriteLock sync.Mutex // 文件书写锁
)

// 自动分隔巡查信息
func UpdateFile(infoTime string) string {
	getTime, err := time.Parse("2006-01-02 15:04", infoTime)
	if err != nil {
		getTime = time.Now()
	}
	datadir := DataFileDir + getTime.Format("2006-01/02/15") + "/"
	datafile := DataFileDir + getTime.Format("2006-01/02/15/04") + DataFileName
	// 判断目录是否存在，不存在需要创建
	_, err = os.Stat(datadir)
	if err != nil && os.IsNotExist(err) {
		err := os.MkdirAll(datadir, os.ModePerm)
		if err != nil {
			log.Fatalln("日志目录创建失败")
		}
	}

	// 判断文件是否存在，不存在需要创建
	if _, err := os.Stat(datafile); err != nil {
		if !os.IsExist(err) {
			newFile, err := os.Create(datafile)
			if err != nil {
				log.Fatalln("日志文件创建失败")
			}
			_ = newFile.Close()
		}
	}
	return datafile
}

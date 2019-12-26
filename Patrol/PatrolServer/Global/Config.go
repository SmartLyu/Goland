package Global

import (
	"log"
	"os"
	"sync"
	"time"
)

var (
	ApiPost          = "8666"                      // 程序端口
	ApiPublicPost    = "8686"                      // 对外程序端口
	DataFileDir      = "/work/data/patrol/"        // 存放历史监控信息
	DataFileName     = ".monitor.log"              // 监控后缀名
	LogFileDir       = "/work/logs/patrol/"        // 日志存放目录(提前准备好)
	LogFileName      = ".patrol.log"               // 日志后缀名
	ErrorFileName    = ".error.log"                // 异常日志后缀名
	AcessLogFileName = ".access.log"               // 记录日志信息
	MonitorShellFile = "/work/sh/PatrolMonitor.sh" // 巡查脚本存放位置
	NatShellFile     = "/work/sh/NatPatrol.sh"     // Nat使用的巡查脚本存放位置
	ErrorMap         = NewErrorMapType()           // 存储报警信息至内存
	NatHostsMap      = NewNatHostsMapType()        // 存储Nat机器中子服务器信息至内存
	FileWriteLock    sync.Mutex                    // 文件书写锁

	ErrorMax              = 2                              // 最多报警次数
	MaxSearchLen    int64 = 300                            // 搜索文件最大次数
	MaxReturnLen    int64 = 1000000                        // 查询预估临界值
	MaxSearchSigLen       = make(chan int, 30)             // 查询并发线程数
	ListenSig             = make(chan int)                 // 监听后台阻塞信号
	ListenPublicSig       = make(chan int)                 // 监听后台公共端口阻塞信号
	CocoUrl               = "http://10.4.0.4:8666/monitor" // coco的端口
	IsPolice              = true                           // 是否报警
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

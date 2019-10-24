package File

import (
	"../Global"
	"../CallPolice"
	"errors"
	"log"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

var logLock sync.Mutex

// 写入日志信息
func WriteLog(message string, datadir string, logfile string) {

	logLock.Lock()
	defer logLock.Unlock()
	// 判断目录是否存在，不存在需要创建
	_, err := os.Stat(datadir)
	if err != nil && os.IsNotExist(err) {
		_ = os.MkdirAll(datadir, os.ModePerm)
	}

	// 判断文件是否存在，不存在需要创建
	if _, err := os.Stat(logfile); err != nil {
		if !os.IsExist(err) {
			newFile, err := os.Create(logfile)
			if err != nil {
				log.Panic(err)
			}
			_ = newFile.Close()
		}
	}

	f, err := os.OpenFile(logfile, os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(errors.New("cacheFileList.yml file create failed. err: " + err.Error()))
	} else {
		// 查找文件末尾的偏移量
		n, _ := f.Seek(0, os.SEEK_END)
		// 从末尾的偏移量开始写入内容
		_, err = f.WriteAt([]byte(time.Now().Format("2006-01-02 15:04:05")+"\t"+
			message+"\n"), n)
	}
	_ = f.Close()

	if err != nil {
		log.Fatal(errors.New("cacheFileList.yml file writed failed. err: " + err.Error()))
	}
}

func WriteErrorLog(message string) {
	dir, file := Global.UpdateLog()
	CallPolice.CallMessage("巡查服务器出现异常：Error \t"+message)
	WriteLog("Error\t"+message, dir, file)
}

func WriteInfoLog(message string) {
	dir, file := Global.UpdateLog()
	WriteLog("Info\t"+message+"\tgoroutine: " + strconv.Itoa(runtime.NumGoroutine()), dir, file)
}

func WriteAccessLog(message string) {
	dir, file := Global.UpdateAcessLog()
	WriteLog(message+"\tgoroutine: " + strconv.Itoa(runtime.NumGoroutine()), dir, file)
}

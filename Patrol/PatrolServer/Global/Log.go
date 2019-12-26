package Global

import (
	"log"
	"os"
	"time"
)

var (
	AccessLog *log.Logger
	InfoLog   *log.Logger
	ErrorLog  *log.Logger
)

func Log() {
	file, err := os.OpenFile(UpdateLog(LogFileName), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(UpdateLog(LogFileName) + " 文件无法打开")
	}
	InfoLog = log.New(file,
		"",
		log.Ldate|log.Ltime|log.Lshortfile)

	file, err = os.OpenFile(UpdateLog(ErrorFileName), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(UpdateLog(ErrorFileName) + " 文件无法打开")
	}
	ErrorLog = log.New(file,
		"",
		log.Ldate|log.Ltime|log.Lshortfile)

	file, err = os.OpenFile(UpdateLog(AcessLogFileName), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(UpdateLog(AcessLogFileName) + " 文件无法打开")
	}
	AccessLog = log.New(file,
		"",
		log.LstdFlags)
}

func CutLog() {
	file, err := os.OpenFile(UpdateLog(LogFileName), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		ErrorLog.Println(UpdateLog(LogFileName) + " 文件无法打开")
	}
	InfoLog.SetOutput(file)

	file, err = os.OpenFile(UpdateLog(AcessLogFileName), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		ErrorLog.Println(UpdateLog(AcessLogFileName) + " 文件无法打开")
	}
	AccessLog.SetOutput(file)

	file, err = os.OpenFile(UpdateLog(ErrorFileName), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		ErrorLog.Println(UpdateLog(ErrorFileName) + " 文件无法打开")
		return
	}
	ErrorLog.SetOutput(file)
}

// 自动分隔日志文件
func UpdateLog(fileName string) string {
	datadir := LogFileDir + time.Now().Format("2006-01") + "/"
	datafile := LogFileDir + time.Now().Format("2006-01/02") + fileName
	// 判断目录是否存在，不存在需要创建
	_, err := os.Stat(datadir)
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

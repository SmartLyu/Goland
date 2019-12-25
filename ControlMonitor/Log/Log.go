package Log

import (
	"log"
	"os"
)

var (
	DebugLog *log.Logger
	InfoLog  *log.Logger
	ErrorLog *log.Logger
)

func Log(fileDir string) error {
	// 判断目录是否存在，不存在需要创建
	_, err := os.Stat(fileDir)
	if err != nil && os.IsNotExist(err) {
		err := os.MkdirAll(fileDir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	// 判断文件是否存在，不存在需要创建
	if _, err := os.Stat(fileDir + "/info.log"); err != nil {
		if !os.IsExist(err) {
			newFile, err := os.Create(fileDir + "/info.log")
			if err != nil {
				return err
			}
			_ = newFile.Close()
		}
	}

	if _, err := os.Stat(fileDir + "/error.log"); err != nil {
		if !os.IsExist(err) {
			newFile, err := os.Create(fileDir + "/info.log")
			if err != nil {
				return err
			}
			_ = newFile.Close()
		}
	}

	if _, err := os.Stat(fileDir + "/detail.log"); err != nil {
		if !os.IsExist(err) {
			newFile, err := os.Create(fileDir + "/info.log")
			if err != nil {
				return err
			}
			_ = newFile.Close()
		}
	}

	file, _ := os.OpenFile(fileDir+"/info.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	InfoLog = log.New(file,
		"",
		log.LstdFlags)

	file, _ = os.OpenFile(fileDir+"/error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	ErrorLog = log.New(file,
		"",
		log.LstdFlags)

	file, _ = os.OpenFile(fileDir+"/detail.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	DebugLog = log.New(file,
		"",
		log.Ldate|log.Ltime|log.Lshortfile)
	return nil
}

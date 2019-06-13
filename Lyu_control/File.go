package main

import (
	"os"
	"time"
)

func UpdateFile() (string, string) {
	return DataFileDir + time.Now().Format("2006-01") + "/",
		DataFileDir + time.Now().Format("2006-01/02") + DataFileName
}

// 记录监控信息
func WriteFile(message string) {

	datadir, datafile := UpdateFile()

	// 判断目录是否存在，不存在需要创建
	_, err := os.Stat(datadir)
	if err != nil && os.IsNotExist(err) {
		err := os.MkdirAll(datadir, os.ModePerm)
		if err != nil {
			return
		}
	}

	// 判断文件是否存在，不存在需要创建
	if _, err := os.Stat(datafile); err != nil {
		if !os.IsExist(err) {
			newFile, err := os.Create(datafile)
			if err != nil {
				return
			}
			_ = newFile.Close()
		}
	}

	f, err := os.OpenFile(datafile, os.O_WRONLY, 0644)
	if err != nil {
		return
	} else {
		// 查找文件末尾的偏移量
		n, _ := f.Seek(0, os.SEEK_END)
		// 从末尾的偏移量开始写入内容
		_, err = f.WriteAt([]byte(time.Now().Format("2006/01/02 15:04:05")+"  "+message+",\n"), n)
	}
	_ = f.Close()

	if err != nil {
		return
	}
	return
}

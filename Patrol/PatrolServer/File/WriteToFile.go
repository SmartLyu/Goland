package File

import (
	"../Global"
	"errors"
	"os"
)

// 记录监控信息
func WriteFile(message string, logTime string) error {

	datafile := Global.UpdateFile(logTime)
	Global.FileWriteLock.Lock()
	defer Global.FileWriteLock.Unlock()

	f, err := os.OpenFile(datafile, os.O_WRONLY, 0644)
	if err != nil {
		return errors.New("cacheFileList.yml file create failed. err: " + err.Error())
	} else {
		// 查找文件末尾的偏移量
		n, _ := f.Seek(0, os.SEEK_END)
		// 从末尾的偏移量开始写入内容
		_, err = f.WriteAt([]byte(message+",\n"), n)
	}
	_ = f.Close()

	if err != nil {
		return errors.New("cacheFileList.yml file writed failed. err: " + err.Error())
	}
	return nil
}

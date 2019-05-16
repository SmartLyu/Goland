package main

import (
	"errors"
	"fmt"
	"github.com/qiniu/api.v7/storage"
	"os"
	"strings"
	"time"
)

func GetFileModTime(path string) (int64, error) {
	f, err := os.Open(path)
	if err != nil {
		return time.Now().Unix(), errors.New("open file error:" + err.Error())
	}
	defer func() {
		_ = f.Close()
	}()

	fi, err := f.Stat()
	if err != nil {
		return time.Now().Unix(), errors.New("stat fileinfo error:" + err.Error())
	}

	return fi.ModTime().UnixNano(), nil
}

// 判断文件是否上传
func CheckFile(uf UploadFile, file string, bucketManager *storage.BucketManager) (bool, error) {

	bucket := uf.Bucket
	key := uf.KeyName
	fileInfo, err := bucketManager.Stat(bucket, key)
	if err != nil {
		return true, nil
	}
	fileTime, err := GetFileModTime(file)
	if err != nil {
		return false, err
	}

	if fileInfo.PutTime*100 <= fileTime {
		return true, nil
	}

	return false, nil
}

// 判断是否上传后，上传文件
func ImportFile(uf UploadFile, file string, dir string, upToken string, bucketManager *storage.BucketManager,addname string) error {

	tmpString := strings.Split(file, "")
	if dir != "" {
		tmpString = strings.Split(file, dir+"/")
	}

	uf.LocalFile = file
	uf.KeyName = tmpString[1]
	uf.KeyName = addname + uf.KeyName
	if uf.KeyName == "" {
		return errors.New("the dir sets error")
	}
	isUpload, err := CheckFile(uf, file, bucketManager)
	if err != nil {
		return err
	}

	if isUpload == true {
		if err := UpDataFile(uf, upToken); err != nil {
			return err
		}
	} else {
		fmt.Println(uf.LocalFile + " is latest")
	}

	return nil
}

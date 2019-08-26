package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"strings"
)

// 设置key的默认值
var (
	accessKeyDefault = ""
	secretKeyDefault = ""
)

func main() {

	var uf UploadFile
	var dir string
	var name string
	var JustCheck string

	// 所有参数外部设定，保持程序无状态
	flag.StringVar(&uf.Bucket, "bucket", "error", "输入桶的名字")
	flag.StringVar(&dir, "dir", "", "输入本地的上传目录(上传过滤该目录名)，不能以/结尾")
	flag.StringVar(&uf.LocalFile, "file", "", "输入本地的上传文件位置")
	flag.StringVar(&accessKey, "accessKey", "***", "输入传输账户accesskey")
	flag.StringVar(&secretKey, "secretKey", "***", "输入传输账户secretKey")
	flag.StringVar(&name, "addName", "", "输入传输后在文件头")
	flag.StringVar(&JustCheck, "justCheck", "0", "是否只是判断文件在远端是否存在，1 代表只判断")
	flag.Parse()

	// 如果需要设定默认值可以再下面设定
	if secretKey == "***" {
		secretKey = secretKeyDefault
	}

	if accessKey == "***" {
		accessKey = accessKeyDefault
	}

	if uf.Bucket == "error" || (dir == "" && uf.LocalFile == "") {
		flag.Usage()
		return
	}

	// 获取七牛云的相关凭证
	upToken := UpDataGetToken(uf)
	bucketManager := GetFileData(uf)

	// 判断用户是否只是想判断是否存在
	if JustCheck == "1" {
		if uf.LocalFile != "" {
			tmpString := uf.LocalFile
			if dir != "" {
				tmpString = strings.TrimPrefix(tmpString, dir+"/")
			}
			uf.KeyName = tmpString
			if uf.KeyName == "" {
				log.Fatal(errors.New("the dir sets error"))
			}
			uf.KeyName = name + uf.KeyName

			IsRight, err := CheckFile(uf, uf.LocalFile, bucketManager)
			if ! IsRight {
				fmt.Println("OK")
			} else {
				fmt.Println("Error: ", err)
			}
			return
		}

		if dir != "" {
			files, err := GetAllFiles(dir)
			if err != nil {
				log.Fatal(errors.New("Get dir error:" + err.Error()))
			}

			for _, i := range files {
				IsRight, err := CheckDir(uf, i, dir, upToken, bucketManager, name)
				tmpString := strings.TrimPrefix(i, dir+"/")
				if ! IsRight {
					fmt.Println(i + " to " + name + tmpString + " is OK")
				} else {
					fmt.Println(i+" to "+name+tmpString+" is Error: ", err)
				}
			}
		}
		return
	}

	// 上传单个文件
	if uf.LocalFile != "" {
		upToken = UpDataGetToken(uf)
		bucketManager = GetFileData(uf)
		fmt.Println("prepare to upload " + uf.LocalFile)
		if err := ImportFile(uf, uf.LocalFile, dir, upToken, bucketManager, name); err != nil {
			log.Fatal(err)
		}
		return
	}

	// 上传整个目录
	if dir != "" {
		files, err := GetAllFiles(dir)
		if err != nil {
			log.Fatal(errors.New("Get dir error:" + err.Error()))
		}

		for _, i := range files {
			upToken = UpDataGetToken(uf)
			bucketManager = GetFileData(uf)
			fmt.Println("prepare to upload " + i)
			if err := ImportFile(uf, i, dir, upToken, bucketManager, name); err != nil {
				log.Fatal(err)
			}
		}
	}
}

package main

import (
	"errors"
	"flag"
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

	// 所有参数外部设定，保持程序无状态
	flag.StringVar(&uf.Bucket, "bucket", "error", "输入桶的名字")
	flag.StringVar(&dir, "dir", "", "输入本地的上传目录(上传过滤该目录名)，不能以/结尾")
	flag.StringVar(&uf.LocalFile, "file", "", "输入本地的上传文件位置")
	flag.StringVar(&accessKey, "accessKey", "***", "输入传输账户accesskey")
	flag.StringVar(&secretKey, "secretKey", "***", "输入传输账户secretKey")
	flag.Parse()

	// 如果需要设定默认值可以再下面设定
	if secretKey == "***" {
		secretKey = secretKeyDefault
	}

	if accessKey == "***" {
		accessKey = accessKeyDefault
	}

	if uf.Bucket == "error" || ( dir == "" && uf.LocalFile ==""){
		flag.Usage()
		return
	}

	// 获取七牛云的相关凭证
	upToken := UpDataGetToken(uf)
	bucketManager := GetFileData(uf)

	// 上传单个文件
	if uf.LocalFile != "" {
		tmpString := strings.Split(uf.LocalFile, "")
		if dir != "" {
			tmpString = strings.Split(uf.LocalFile, dir+"/")
		}
		uf.KeyName = tmpString[1]
		if uf.KeyName == "" {
			log.Fatal(errors.New("the dir sets error"))
		}

		if err := UpDataFile(uf, upToken); err != nil {
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
			if err := ImportFile(uf, i, dir, upToken, bucketManager); err != nil {
				log.Fatal(err)
			}
		}
	}
}

package main

import (
	"flag"
	"fmt"
	"github.com/qiniu/api.v7/auth"
	"github.com/qiniu/api.v7/storage"
	"time"
)

type QiNiu struct {
	accessKey string
	secretKey string
	url    string
	keyfile   string
}

func main() {
	var qn QiNiu
	var style string
	flag.StringVar(&qn.url, "bucket", "error", "输入cdn的url")
	flag.StringVar(&qn.accessKey, "accessKey", "***", "输入传输账户accesskey")
	flag.StringVar(&qn.secretKey, "secretKey", "***", "输入传输账户secretKey")
	flag.StringVar(&qn.keyfile, "keyfile", "error", "输入传输后在文件头")
	flag.StringVar(&style, "style", "private", "公开或者私有空间 (public | private)")
	flag.Parse()

	if style == "public" {
		// 公开空间访问

		if qn.url == "error" || qn.keyfile == "error" {
			flag.Usage()
			return
		}

		publicAccessURL := storage.MakePublicURL(qn.url, qn.keyfile)
		fmt.Println(publicAccessURL)

	} else if style == "private" {
		// 私有空间访问

		if qn.url == "error" || qn.keyfile == "error" || qn.secretKey == "***" || qn.accessKey == "***" {
			flag.Usage()
			return
		}

		mac := auth.New(qn.accessKey, qn.secretKey)
		deadline := time.Now().Add(time.Second * 3600).Unix() //1小时有效期
		privateAccessURL := storage.MakePrivateURL(mac, qn.url, qn.keyfile, deadline)
		fmt.Println(privateAccessURL)
	} else {
		flag.Usage()
	}

}

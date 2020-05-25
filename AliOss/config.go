package main

import "log"

type configJson struct {
	Endpoint        string `json:"endpoint"`          // 存储桶的Endpoint信息
	AccessKeyId     string `json:"access_key_id"`     // ak
	AccessKeySecret string `json:"access_key_secret"` // sk
	BucketName      string `json:"bucket_name"`       // 桶名
	ObjectName      string `json:"object_name"`       // OSS存储名
	LocalFilename   string `json:"local_filename"`    // 本地文件和地址
}

var (
	InfoLog    *log.Logger
	ErrorLog   *log.Logger
	routing    int
	tmpFileDir string
	module     string
)

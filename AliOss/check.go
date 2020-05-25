package main

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"hash/crc64"
	"io"
	"os"
	"strconv"
)

func check(id configJson) {
	InfoLog.Println("start to check", id.LocalFilename, "and", id.ObjectName, "in", id.BucketName)
	// 创建OSSClient实例。
	client, err := oss.New(id.Endpoint, id.AccessKeyId, id.AccessKeySecret)
	if err != nil {
		ErrorLog.Println("Error:", err)
		os.Exit(-1)
	}

	objectName := id.ObjectName

	// 获取存储空间。
	bucket, err := client.Bucket(id.BucketName)
	if err != nil {
		ErrorLog.Println("Error:", err)
		os.Exit(-1)
	}

	// 列举包含指定前缀的文件。
	ret, err := bucket.ListObjectVersions(oss.Prefix(objectName))
	if err != nil {
		ErrorLog.Println("Error:", err)
		os.Exit(-1)
	}

	// 获取指定版本文件的部分元信息。
	props, err := bucket.GetObjectMeta(objectName, oss.VersionId(ret.ObjectVersions[0].VersionId))
	if err != nil {
		ErrorLog.Println("bucket.GetObjectMeta Error:", err)
		os.Exit(-1)
	}
	objectCrc64ecma := props.Get("X-Oss-Hash-Crc64ecma")
	InfoLog.Println("Object Hash Crc64ecma:", objectCrc64ecma)

	// 计算文件的crc64值
	file, err := os.Open(id.LocalFilename)
	if err != nil {
		ErrorLog.Println("os.Open Error:", err)
		os.Exit(-1)
	}
	defer func() {
		_ = file.Close()
	}()
	hash := crc64.New(crc64.MakeTable(crc64.ECMA))
	if _, err = io.Copy(hash, file); err != nil {
		ErrorLog.Println("io.Copy Error:", err)
		os.Exit(-1)
	}
	localCrc64ecma := strconv.FormatUint(hash.Sum64(), 10)
	InfoLog.Println("Local File Hash Crc64ecma:", localCrc64ecma)

	if objectCrc64ecma != localCrc64ecma {
		ErrorLog.Println("check Crc64ecma is different")
		os.Exit(-1)
	}
}

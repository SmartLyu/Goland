package main

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
	"net/http"
	"os"
)

// 定义进度条监听器。
type OssProgressListener struct {
}

// 定义进度变更事件处理函数。
func (listener *OssProgressListener) ProgressChanged(event *oss.ProgressEvent) {
	switch event.EventType {
	case oss.TransferStartedEvent:
		fmt.Printf("Transfer Started, ConsumedBytes: %d, TotalBytes %d.\n",
			event.ConsumedBytes, event.TotalBytes)
	case oss.TransferDataEvent:
		fmt.Printf("\rTransfer Data, ConsumedBytes: %d, TotalBytes %d, %d%%.",
			event.ConsumedBytes, event.TotalBytes, event.ConsumedBytes*100/event.TotalBytes)
	case oss.TransferCompletedEvent:
		fmt.Printf("\nTransfer Completed, ConsumedBytes: %d, TotalBytes %d.\n",
			event.ConsumedBytes, event.TotalBytes)
	case oss.TransferFailedEvent:
		fmt.Printf("\nTransfer Failed, ConsumedBytes: %d, TotalBytes %d.\n",
			event.ConsumedBytes, event.TotalBytes)
	default:
	}
}

func upload(id configJson) {
	InfoLog.Println("start", routing, "processes to pull", id.LocalFilename, "to", id.ObjectName, "in", id.BucketName)

	// 创建目录
	dirInfo, err := os.Stat(tmpFileDir)
	if err != nil && os.IsNotExist(err) {
		err := os.MkdirAll(tmpFileDir, 0700)
		if err != nil {
			ErrorLog.Println(tmpFileDir, "os.Mkdir error:", err)
		}
	} else {
		if err != nil {
			ErrorLog.Println(tmpFileDir, "stat error:", err)
		}
		if !dirInfo.IsDir() {
			ErrorLog.Println(tmpFileDir, "is not dirtory")
		}
	}

	// 创建OSSClient实例。
	client, err := oss.New(id.Endpoint, id.AccessKeyId, id.AccessKeySecret)
	if err != nil {
		ErrorLog.Println("oss.New Error: ", err)
		os.Exit(-1)
	}
	bucketName := id.BucketName
	objectName := id.ObjectName
	locaFilename := id.LocalFilename
	// 用oss.GetResponseHeader获取返回header。
	var retHeader http.Header

	// 获取存储空间。
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		ErrorLog.Println("client.Bucket Error:", err)
		os.Exit(-1)
	}

	var tmpfile = tmpFileDir + "/" + fmt.Sprintf("%x", md5.Sum([]byte(locaFilename))) + ".process"
	InfoLog.Println("tmp file write in", tmpfile)

	// 计算md5，校验上传情况
	file, err := os.Open(id.LocalFilename)
	if err != nil {
		ErrorLog.Println("os.Open Error:", err)
		os.Exit(-1)
	}
	defer func() {
		_ = file.Close()
	}()
	hash := md5.New()
	if _, err = io.Copy(hash, file); err != nil {
		ErrorLog.Println("io.Copy Error:", err)
		os.Exit(-1)
	}
	strMd5 := base64.StdEncoding.EncodeToString(hash.Sum(nil))

	// 分片大小100K，routing个协程并发上传分片，使用断点续传。
	err = bucket.UploadFile(objectName, locaFilename, 100*1024, oss.Routines(routing),
		oss.Checkpoint(true, tmpfile+".upload"), oss.Progress(&OssProgressListener{}),
		oss.GetResponseHeader(&retHeader), oss.ContentMD5(strMd5))

	if err != nil {
		ErrorLog.Println("bucket.UploadFile Error:", err)
		os.Exit(-1)
	}

	// 打印x-oss-version-id。
	InfoLog.Println("this case's x-oss-version-id:", oss.GetVersionId(retHeader), " upload successfully")

	/* 分片操作
	chunks, err := oss.SplitFileByPartNum(locaFilename, 3)
	if err != nil {
		ErrorLog.Println("oss.SplitFileByPartNum Error:", err)
		os.Exit(-1)
	}

	fd, err := os.Open(locaFilename)
	if err != nil {
		ErrorLog.Println("os.Open Error:", err)
		os.Exit(-1)
	}

	defer func() {
		_ = fd.Close()
	}()

	// 步骤1：初始化一个分片件。
	InfoLog.Println("start to init upload case")
	imur, err := bucket.InitiateMultipartUpload(objectName)
	if err != nil {
		ErrorLog.Println("bucket.InitiateMultipartUpload Error:", err)
		os.Exit(-1)
	}

	// 步骤2：上传分片。
	InfoLog.Println("start uploading seek")
	var parts []oss.UploadPart
	for _, chunk := range chunks {
		_, err = fd.Seek(chunk.Offset, os.SEEK_SET)
		if err != nil {
			ErrorLog.Println("fd.Seek Error:", err)
			os.Exit(-1)
		}
		// 对每个分片调用UploadPart方法上传。
		part, err := bucket.UploadPart(imur, fd, chunk.Size, chunk.Number)
		if err != nil {
			ErrorLog.Println("bucket.UploadPart Error:", err)
			os.Exit(-1)
		}
		parts = append(parts, part)
	}

	// 步骤3：完成分片上传。
	InfoLog.Println("complete uploading")
	cmur, err := bucket.CompleteMultipartUpload(imur, parts, oss.GetResponseHeader(&retHeader))
	if err != nil {
		ErrorLog.Println("bucket.CompleteMultipartUpload Error:", err)
		os.Exit(-1)
	}
	// 打印cmur的值。
	InfoLog.Println("this case's cmur:", cmur)
	// 打印x-oss-version-id。
	InfoLog.Println("this case's x-oss-version-id:", oss.GetVersionId(retHeader))
	*/
}

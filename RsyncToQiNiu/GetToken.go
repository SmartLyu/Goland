package main

import (
	"github.com/qiniu/api.v7/auth"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
)

var (
	accessKey string
	secretKey string
)

// 获取上传token
func UpDataGetToken(uf UploadFile) string {
	// 设置上传凭证有效期
	putPolicy := storage.PutPolicy{
		Scope: uf.Bucket,
	}
	mac := auth.New(accessKey, secretKey)
	putPolicy.Expires = 7200 //示例2小时有效期

	upToken := putPolicy.UploadToken(mac)
	return upToken
}

// 获取查看文件信息的认证
func GetFileData(uf UploadFile) *storage.BucketManager {
	accessKeyDefault := ""
	secretKeyDefault := ""

	mac := qbox.NewMac(accessKeyDefault, secretKeyDefault)
	cfg := storage.Config{
		// 是否使用https域名进行资源管理
		UseHTTPS: false,
	}
	return storage.NewBucketManager(mac, &cfg)
}

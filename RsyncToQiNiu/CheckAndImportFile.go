package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/qiniu/api.v7/storage"
	"io"
	"os"
	"strings"
)

const (
	BLOCK_BITS = 22 // Indicate that the blocksize is 4M
	BLOCK_SIZE = 1 << BLOCK_BITS
)

func BlockCount(fsize int64) int {
	return int((fsize + (BLOCK_SIZE - 1)) >> BLOCK_BITS)
}

func CalSha1(b []byte, r io.Reader) ([]byte, error) {

	h := sha1.New()
	_, err := io.Copy(h, r)
	if err != nil {
		return nil, err
	}
	return h.Sum(b), nil
}

// 计算哈希值
func GetEtag(filename string) (etag string, err error) {

	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer func() {
		_ = f.Close()
	}()

	fi, err := f.Stat()
	if err != nil {
		return
	}

	fsize := fi.Size()
	blockCnt := BlockCount(fsize)
	sha1Buf := make([]byte, 0, 21)

	if blockCnt <= 1 { // file size <= 4M
		sha1Buf = append(sha1Buf, 0x16)
		sha1Buf, err = CalSha1(sha1Buf, f)
		if err != nil {
			return
		}
	} else { // file size > 4M
		sha1Buf = append(sha1Buf, 0x96)
		sha1BlockBuf := make([]byte, 0, blockCnt*20)
		for i := 0; i < blockCnt; i ++ {
			body := io.LimitReader(f, BLOCK_SIZE)
			sha1BlockBuf, err = CalSha1(sha1BlockBuf, body)
			if err != nil {
				return
			}
		}
		sha1Buf, _ = CalSha1(sha1Buf, bytes.NewReader(sha1BlockBuf))
	}
	etag = base64.URLEncoding.EncodeToString(sha1Buf)
	return
}

// 判断文件是否上传
func CheckFile(uf UploadFile, file string, bucketManager *storage.BucketManager) (bool, error) {

	bucket := uf.Bucket
	key := uf.KeyName
	fileInfo, err := bucketManager.Stat(bucket, key)
	if err != nil {
		return true, nil
	}
	fileHash, err := GetEtag(file)
	if err != nil {
		return false, err
	}

	if fileInfo.Hash != fileHash {
		return true, nil
	}

	return false, nil
}

// 判断是否上传后，上传文件
func ImportFile(uf UploadFile, file string, dir string, upToken string, bucketManager *storage.BucketManager, addname string) error {

	tmpString := file
	if dir != "" {
		tmpString = strings.TrimPrefix(file, dir+"/")
	}

	uf.LocalFile = file
	uf.KeyName = tmpString
	uf.KeyName = addname + uf.KeyName

	fmt.Println("start to upload " + uf.LocalFile + " to " + uf.KeyName)
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

// 判断是否在七牛已存在该文件
func CheckDir(uf UploadFile, file string, dir string, upToken string, bucketManager *storage.BucketManager, addname string) (bool, error) {

	tmpString := file
	if dir != "" {
		tmpString = strings.TrimPrefix(file, dir+"/")
	}

	uf.LocalFile = file
	uf.KeyName = tmpString
	uf.KeyName = addname + uf.KeyName
	if uf.KeyName == "" {
		return true, errors.New("the dir sets error")
	}
	isUpload, err := CheckFile(uf, file, bucketManager)
	if err != nil {
		return true, err
	}

	return isUpload, nil

}

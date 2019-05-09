package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"context"
	"github.com/qiniu/api.v7/storage"
)

// 上传文件的基本数据
type UploadFile struct {
	LocalFile string
	KeyName   string
	Bucket    string
}

// 获取md5值
func md5Hex(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

type ProgressRecord struct {
	Progresses []storage.BlkputRet `json:"progresses"`
}

// 上传文件
func UpDataFile(uf UploadFile, upToken string) error {
	localFile := uf.LocalFile
	key := uf.KeyName
	bucket := uf.Bucket

	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Zone = &storage.ZoneHuadong
	// 是否使用https域名
	cfg.UseHTTPS = false
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false

	// 必须仔细选择一个能标志上传唯一性的 recordKey 用来记录上传进度
	// 我们这里采用 md5(bucket+key+local_path+local_file_last_modified)+".progress" 作为记录上传进度的文件名
	fileInfo, statErr := os.Stat(localFile)
	if statErr != nil {
		return statErr
	}

	fileSize := fileInfo.Size()
	fileLmd := fileInfo.ModTime().UnixNano()
	recordKey := md5Hex(fmt.Sprintf("%s:%s:%s:%s", bucket, key, localFile, fileLmd)) + ".progress"

	// 指定的进度文件保存目录，实际情况下，请确保该目录存在，而且只用于记录进度文件
	recordDir := "/Users/jemy/Temp/progress"
	mErr := os.MkdirAll(recordDir, 0755)
	if mErr != nil {
		return errors.New("mkdir for record dir error:" + mErr.Error())

	}

	recordPath := filepath.Join(recordDir, recordKey)
	progressRecord := ProgressRecord{}
	// 尝试从旧的进度文件中读取进度
	recordFp, openErr := os.Open(recordPath)
	if openErr == nil {
		progressBytes, readErr := ioutil.ReadAll(recordFp)
		if readErr == nil {
			mErr := json.Unmarshal(progressBytes, &progressRecord)
			if mErr == nil {
				// 检查context 是否过期，避免701错误
				for _, item := range progressRecord.Progresses {
					if storage.IsContextExpired(item) {
						fmt.Println(item.ExpiredAt)
						progressRecord.Progresses = make([]storage.BlkputRet, storage.BlockCount(fileSize))
						break
					}
				}
			}
		}
		_ = recordFp.Close()
	}

	if len(progressRecord.Progresses) == 0 {
		progressRecord.Progresses = make([]storage.BlkputRet, storage.BlockCount(fileSize))
	}

	resumeUploader := storage.NewResumeUploader(&cfg)
	ret := storage.PutRet{}
	progressLock := sync.RWMutex{}

	putExtra := storage.RputExtra{
		Progresses: progressRecord.Progresses,
		Notify: func(blkIdx int, blkSize int, ret *storage.BlkputRet) {
			progressLock.Lock()
			defer progressLock.Unlock()
			//将进度序列化，然后写入文件
			progressRecord.Progresses[blkIdx] = *ret
			progressBytes, _ := json.Marshal(progressRecord)
			fmt.Println("write progress file", blkIdx, recordPath)
			wErr := ioutil.WriteFile(recordPath, progressBytes, 0644)
			if wErr != nil {
				fmt.Println("write progress file error,", wErr)
			}
		},
	}
	err := resumeUploader.PutFile(context.Background(), &ret, upToken, key, localFile, &putExtra)
	if err != nil {
		return err
	}
	//上传成功之后，一定记得删除这个进度文件
	_ = os.Remove(recordPath)
	fmt.Println(ret.Key, ret.Hash)
	return nil
}

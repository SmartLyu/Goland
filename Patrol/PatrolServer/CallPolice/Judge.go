package CallPolice

import (
	"../Global"
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"time"
)

func Judge(monitorjson Global.MonitorJson){
	if ! monitorjson.Status {
		pwd := Global.DataFileDir
		des := pwd + time.Now().Format("2006-01/02") + Global.DataFileName

		_, err := FindWorkInFile(des, time.Now().Add(-time.Minute * 1).Format("2006-01-02 15:04"),
			monitorjson.Info ,monitorjson.IP , "false")
		if err == nil {
			CallPolice(monitorjson.Hostname+" 的 "+ monitorjson.Info + " 异常 \n具体服务器信息：\n   " +monitorjson.IP)
		}
	}
}

//在文件中查找关键字
func FindWorkInFile(path string, word ...string) (string, error) {

	// 遍历获取所有word值
	re := make(map[int][]byte)
	reNum := 0
	for _, j := range word {
		if j != "" {
			re[reNum] = []byte(j)
			reNum ++
		}
	}

	findWord := ""
	tmpword := ""

	// 访问文件
	File, err := os.Open(path)
	defer func() {
		_ = File.Close()
	}()

	if err != nil {
		return findWord, err
	}

	r := bufio.NewReader(File)
	for {
		b, _, err := r.ReadLine()

		// 排除二进制文件
		//if bytes.Contains(b, []byte{0}) {
		//  return findWord, nil
		//}

		if err != nil {
			if err == io.EOF {
				if tmpword = findOrNot(b, re); tmpword != "" {
					findWord = findWord + "\n" + tmpword
				}
				if findWord == "" {
					err = errors.New(" Cannot Find Suitable Value ")
					return findWord, err
				}
				return findWord, nil
			}
			return findWord, err
		}

		if tmpword = findOrNot(b, re); tmpword != "" {
			findWord = findWord + "\n" + tmpword
		}
	}
}

// 进行二进制字段比对
func findOrNot(b []byte, re map[int][]byte) string {
	findWord := ""
	for _, j := range re {
		if bytes.Contains(b, j) {
			findWord = string(bytes.TrimSpace(b))
		} else {
			findWord = ""
			break
		}
	}

	return findWord
}
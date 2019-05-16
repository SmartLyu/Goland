package File

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
)

//在文件中查找关键字
func FindWorkInFile(path string, word string) (string, error) {
	re := []byte(word)
	findWord := ""

	File, err := os.Open(path)
	defer func() {
		_ = File.Close()
	}()

	if err != nil {
		fmt.Println(path, ":", err)
		return findWord, err
	}

	r := bufio.NewReader(File)
	for {
		b, _, err := r.ReadLine()

		// 排除二进制文件
	//	if bytes.Contains(b, []byte{0}) {
	//		return findWord, nil
	//	}
		if err != nil {
			if err == io.EOF {
				if bytes.Contains(b, re) {
					findWord = findWord + "\n" + string(bytes.TrimSpace(b))
				}
				if findWord == "" {
					err = errors.New("Cannot Find Suitable Value ")
					return findWord, err
				}
				return findWord, nil
			}
			return findWord, err
		}
		if bytes.Contains(b, re) {
			findWord = findWord + "\n" + string(bytes.TrimSpace(b))
		}
	}
}
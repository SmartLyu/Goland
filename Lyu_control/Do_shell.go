package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

// 执行shell命令
func myCmd(bash string, shell ...string) error {
	contentArray := make([]string, 0, 5)
	cmd := exec.Command(bash, shell...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		WriteFile(fmt.Sprint(cmd.Stderr) + " error=>" + err.Error())
	}

	_ = cmd.Start()

	reader := bufio.NewReader(stdout)

	contentArray = contentArray[0:0]
	var index int
	//实时循环读取输出流中的一行内容
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		WriteFile(line)
		index++
		contentArray = append(contentArray, line)
	}
	err = cmd.Wait()

	defer func() {
		if err = cmd.Process.Kill() ; err != nil {
			WriteFile("Kill Error : " +err.Error())
		}
	}()

	if err != nil {
		WriteFile("Execute Shell: " + strings.Join(shell," "))
		return errors.New("failed with error:" + err.Error())
	}

	return nil
}

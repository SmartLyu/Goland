package File

import (
	"bufio"
	"errors"
	"io"
	"os/exec"
)

// 执行shell命令
func MyCmd(bash string, shell ...string) error {
	contentArray := make([]string, 0, 5)
	cmd := exec.Command(bash, shell...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		WriteErrorLog("error=>" + err.Error())
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
		WriteErrorLog(line)
		index++
		contentArray = append(contentArray, line)
	}
	err = cmd.Wait()

	if err != nil {
		shellstring := ""
		for _, i := range shell {
			shellstring = shellstring + i
		}
		WriteInfoLog("Execute Shell: " + shellstring + ";")
		return errors.New("failed with error:" + err.Error())
	}

	return nil
}

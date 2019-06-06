package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os/exec"
	"strings"
)

var installUrl = "http://139.159.233.66:8888/download/falcon-agent.tar.gz"
var falconDir = "/work/servers/FalconAgent"

// 获取变量值
var (
	hostname string
	ip       string
)

// 执行shell命令
func myCmd(bash string, shell ...string) error {
	contentArray := make([]string, 0, 5)
	cmd := exec.Command(bash, shell...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(cmd.Stderr, "error=>", err.Error())
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
		fmt.Print(line)
		index++
		contentArray = append(contentArray, line)
	}
	err = cmd.Wait()

	if err != nil {
		fmt.Printf("Execute Shell %s: ", shell)
		return errors.New("failed with error:" + err.Error())
	}

	return nil
}

// 更具需求替换文件中的内容
func myChange(file string, word string, change string) error {
	input, err := ioutil.ReadFile(file)
	if err != nil {
		return errors.New("failed with error:" + err.Error())
	}
	lines := strings.Split(string(input), "\n")
	for i, line := range lines {
		if strings.Contains(line, word) {
			lines[i] = change
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(file, []byte(output), 0644)
	if err != nil {
		return errors.New("failed with error:" + err.Error())
	}
	return nil
}

// 检查指定端口的占用情况
func myCheckPort(port string) error {

	listener, err := net.Listen("tcp", "127.0.0.1:"+port)
	if err != nil {
		return err
	}

	if err = listener.Close(); err != nil {
		return err
	}
	return nil
}

func main() {

	flag.StringVar(&hostname, "host", "Error", "hostname")
	flag.StringVar(&ip, "falcon", "Error", "open-falcon ip")
	flag.Parse()

	if hostname == "Error" || ip == "Error" {
		flag.Usage()
		return
	}

	// 检查agent是否开启了
	if err := myCheckPort("1988"); err != nil {
		fmt.Println("1988 port is exist")
		log.Fatal(err)
	}

	// 创建目录
	if err := myCmd("/usr/bin/mkdir", "-p", falconDir); err != nil {
		if err := myCmd("/usr/bin/ls", falconDir); err != nil {
			fmt.Println("mkdir falcon dir error")
			log.Fatal(err)
		}
	}

	if err := myCmd("/usr/bin/mkdir", "-p", "/work/src"); err != nil {
		if err := myCmd("/usr/bin/ls", "/work/src"); err != nil {
			fmt.Println("mkdir work dir error")
			log.Fatal(err)
		}
	}

	// 下载agent服务包
	if err := myCmd("/usr/bin/wget", "-q", installUrl, "-O", "/work/src/FalconAgent.tar.gz"); err != nil {
		fmt.Println("wget tar file error")
		log.Fatal(err)
	}

	// 解压数据包
	if err := myCmd("/bin/tar", "-xf", "/work/src/FalconAgent.tar.gz", "-C", falconDir); err != nil {
		fmt.Println("untar is error")
		log.Fatal(err)
	}

	// 修改配置
	if err := myChange(falconDir+"/agent/config/cfg.json", "hostname", "    \"hostname\": \""+hostname+"\","); err != nil {
		fmt.Println("change file for hostname is error")
		log.Fatal(err)
	}

	if err := myChange(falconDir+"/agent/config/cfg.json", "6030", "        \"addr\": \""+ip+":6030\","); err != nil {
		fmt.Println("change file for hostname is error")
		log.Fatal(err)
	}

	if err := myChange(falconDir+"/agent/config/cfg.json", "8433", "            \""+ip+":8433\""); err != nil {
		fmt.Println("change file for hostname is error")
		log.Fatal(err)
	}

	// 启动agent
	if err := myCmd("/bin/bash", falconDir+"/star-agent"); err != nil {
		fmt.Println("start agent is error")
		log.Fatal(err)
	}

	// 检查是否启动成功
	if err := myCheckPort("1988"); err != nil {
		fmt.Printf("\n %c[1;40;32m%s%c[0m\n\n", 0x1B, "falcon is ready", 0x1B)
		return
	} else {
		fmt.Println("falcon start error")
	}
}

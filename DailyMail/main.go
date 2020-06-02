package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"
)

func main() {
	var title string
	flag.StringVar(&title, "title", "Error", "Title")
	flag.Parse()
	if title == "Error" {
		flag.Usage()
		return
	}

	// 通知发信
	SendMessage(time.Now().Format("2006年01月02日 15时04分05秒") + "\n" + "准备发送" + title)

	m.SetHeader("Subject", title) // 主题

	r, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
		return
	}
	err = SendMail(string(r))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("发送邮件成功")
}

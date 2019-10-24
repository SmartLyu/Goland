package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
)

func main()  {
	var title string
	flag.StringVar(&title, "title", "Error", "Title")
	flag.Parse()
	if title == "Error" {
		flag.Usage()
		return
	}

	m.SetHeader("Subject", title) // 主题

	r , err:= ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
		return
	}
	err = test_SendMail(string(r))
	if  err != nil {
		log.Fatal(err)
	}
	fmt.Println("发送邮件成功")
}
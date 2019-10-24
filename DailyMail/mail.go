package main

import (
	"errors"
	"fmt"
	"github.com/go-gomail/gomail"
)

func SendMail(body string) error {

	mInit()

	m.SetBody("text/html",body) // 正文

	if err := d.DialAndSend(m); err != nil {
		return errors.New("发送失败: " + err.Error())
	}

	return nil
}

func test_SendMail(body string) error {

	test_d := gomail.NewPlainDialer("smtp.qq.com", 465,
		"1543376365@qq.com", "fmcwitvwusoajibh") // 发送邮件服务器、端口、发件人账号、发件人密码

	m.SetAddressHeader("From",
		"1543376365@qq.com" /*"发件人地址"*/, "yuzhiyuan") // 发件人

	m.SetHeader("To",
		m.FormatAddress("1543376365@qq.com", "yuzhiyuan")) // 收件人

	m.SetBody("text/html",body) // 正文
	fmt.Println(body)

	if err := test_d.DialAndSend(m); err != nil {
		return errors.New("发送失败: " + err.Error())
	}

	return nil
}
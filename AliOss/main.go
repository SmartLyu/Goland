package main

import (
	"flag"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"log"
	"os"
	"reflect"
)

func main() {
	var conf configJson
	var p = make([]string, reflect.ValueOf(conf).NumField(), reflect.ValueOf(conf).NumField())

	InfoLog = log.New(os.Stdout,
		"Info ",
		log.Ldate|log.Ltime|log.Lshortfile)

	ErrorLog = log.New(os.Stderr,
		"Error ",
		log.Ldate|log.Ltime|log.Lshortfile)
	InfoLog.Println("OSS Go SDK Version: ", oss.Version)

	// 反射对象configJson，遍历获取所有值
	for i := 0; i < reflect.ValueOf(conf).NumField(); i++ {
		name := reflect.TypeOf(conf).Field(i).Tag.Get("json")
		flag.StringVar(&p[i], name, "error", "输入"+name+"的值")
	}
	flag.IntVar(&routing, "route", 3, "输入上传并发数")
	flag.StringVar(&tmpFileDir, "tmpDir", "/work/tmp/uploadToOss", "输入临时文件存放位置")
	flag.StringVar(&module, "module", "check", "输入使用的功能模块(upload、check、all)")

	// 检验是否所有值都存在
	flag.Parse()
	for i := 0; i < reflect.ValueOf(conf).NumField(); i++ {
		reflect.ValueOf(&conf).Elem().Field(i).SetString(p[i])
		if reflect.ValueOf(conf).Field(i).String() == "error" {
			flag.Usage()
			ErrorLog.Println("please input: ", reflect.TypeOf(conf).Field(i).Tag.Get("json"))
			os.Exit(-1)
		}
	}

	switch module {
	case "upload":
		upload(conf)
	case "check":
		err := check(conf)
		if err != nil {
			ErrorLog.Println("check is error", err)
			os.Exit(-1)
		}
	case "all":
		err := check(conf)
		if err != nil {
			upload(conf)
		}
	default:
		flag.Usage()
		ErrorLog.Println("please input right module is not: ", module)
		os.Exit(-1)
	}
}

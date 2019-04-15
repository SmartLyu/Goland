package main

import (
	"flag"
	"fmt"
	"log"
)

var (
	kind string
	url  string
	body string
)

func main() {
	flag.StringVar(&kind, "kind", "Error", "post or get")
	flag.StringVar(&url, "url", "Error", "the web's ip and url")
	flag.StringVar(&body, "body", "Error", "post json")
	flag.Parse()

	if url == "Error" || kind == "Error" {
		flag.Usage()
		return
	}

	if kind == "get" {
		data, hea, err := httpGet(url)
		if err != nil {
			log.Panic(err)
		}
		fmt.Println("header is", hea)
		ShowJson(data)

	} else if kind == "post" && body != "Error" {
		data, hea, err := httpPostJson(url, body)
		if err != nil {
			log.Panic(err)
		}
		fmt.Println("header is", hea)
		fmt.Println(data)

	} else {
		flag.Usage()

	}

}

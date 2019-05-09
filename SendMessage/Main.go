package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	id := SecretId

	flag.StringVar(&id.content, "message", "Hello World", "the message to Lyu weixin")
	flag.Parse()

	if err := SendWeiXinMessage(id); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Message:")
		fmt.Println(id.content)
		fmt.Println("to ", id.agentid, "successfully")
	}
}

package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	id := SecretId
	id.content = "Hello World!"
	flag.StringVar(&id.content, "message", "Hello World", "the message to Lyu weixin")

	if err := SendWeiXinMessage(id); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Message:")
		fmt.Println(id.content)
		fmt.Println("to ", id.agentid, "successfully")
	}
}

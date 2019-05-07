package main

import (
	"fmt"
	"log"
)

func main() {
	id := SecretId
	id.content = "Hello World!"

	if err := SendWeiXinMessage(id); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Message:")
		fmt.Println(id.content)
		fmt.Println("to ", id.agentid, "successfully")
	}
}

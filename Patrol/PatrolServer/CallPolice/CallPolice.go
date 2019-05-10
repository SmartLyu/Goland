package CallPolice

import (
	"fmt"
	"log"
)

func CallPolice(message string) {
	id := SecretId
	id.content = message

	if err := SendWeiXinMessage(id); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Message:")
		fmt.Println(id.content)
		fmt.Println("to ", id.agentid, "successfully")
	}
}

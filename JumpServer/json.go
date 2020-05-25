package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type PostTokenJson struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func ReadProjectJsonString(jsonData string) {
	var v interface{}
	if err := json.Unmarshal([]byte(jsonData), &v); err != nil {
		log.Fatalln(err)
	}
	data := v.(map[string]interface{})
	data_2 := data["data"].(map[string]interface{})
	data_3 := data_2["data"].([]interface{})
	for k, v := range data_3 {
		fmt.Println(k, v)
	}
}

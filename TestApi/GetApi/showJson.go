package main

import (
	"encoding/json"
	"fmt"
)

func ShowJson(jsonstr string) {
	fmt.Println(jsonstr)
	js, _ := json.Marshal(&jsonstr)
	jsIndent, _ := json.MarshalIndent(&jsonstr, "", "   ")
	fmt.Println("\njs:\n", string(js), "\n\njsIndent:\n", string(jsIndent))
}

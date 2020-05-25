package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func httpJson(method string, url string, jsonStr string) {

	jsonbyte := []byte(jsonStr)

	// method 为post、get请求
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonbyte))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		ErrorLog.Println(err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			ErrorLog.Println(err)
		}
	}()

	hea := resp.Header
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ErrorLog.Println(err)
	}

	InfoLog.Println("head is ", hea, "; body is ", body)

}

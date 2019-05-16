package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
)

func httpPostJson( str string , status string) (string, http.Header, error) {

	url := MonitorUrl
	jsonbyte := []byte(PostJson(str,status))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonbyte))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", http.Header{}, err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Panic(err)
		}
	}()

	hea := resp.Header
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", http.Header{}, err
	}

	return string(body), hea, nil
}

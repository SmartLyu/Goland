package DispatchApi

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func httpJson(url string, httpType string, jsonbyte []byte) (string, http.Header, error) {

	fmt.Println(time.Now().Format("2006-01-02 15:04:05 +0800 CST")+" ", "Start ", httpType, url, string(jsonbyte))
	req, err := http.NewRequest(httpType, url, bytes.NewBuffer(jsonbyte))
	if err != nil {
		return "", http.Header{}, err
	}

	req.Header.Set("content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", http.Header{}, err
	}

	hea := resp.Header
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", http.Header{}, err
	}

	if err := resp.Body.Close(); err != nil {
		return "", http.Header{}, err
	}
	return string(body), hea, nil
}
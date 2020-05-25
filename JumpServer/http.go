package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func httpPostJson(url string, jsonbyte []byte) (string, error) {

	var resp *http.Response
	var err error
	resp, err = http.Post(url,
		"application/json",
		bytes.NewBuffer(jsonbyte))
	if err != nil {
		ErrorLog.Println(err)
		return "", err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			ErrorLog.Println(err)
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ErrorLog.Println(err)
		return "", err
	}

	InfoLog.Println("head is ", resp.Header, "; status is ", resp.Status, "; body is ", string(body))
	return string(body), nil
}

func httpGetJson(url string, keyVale map[string]string) (string, error) {
	var resp *http.Response
	var err error

	req, err := http.NewRequest("get", url, nil)
	if err != nil {
		ErrorLog.Println(err)
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range keyVale {
		req.Header.Add(k, v)
	}

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		ErrorLog.Println(err)
		return "", err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			ErrorLog.Println(err)
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ErrorLog.Println(err)
		return "", err
	}

	InfoLog.Println("head is ", resp.Header, "; status is ", resp.Status, "; body is ", string(body))
	return string(body), nil
}

package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

func httpGet(url string) (string, http.Header, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", http.Header{}, err
	}

	defer func() {
		if err = resp.Body.Close(); err != nil {
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

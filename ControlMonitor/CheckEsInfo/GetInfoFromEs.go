package CheckEsInfo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func httpJSON(url string, httpType string, jsonbyte []byte) (string, http.Header, error) {
	fmt.Println("Start " + httpType + " " + url + " : " + string(jsonbyte))
	req, err := http.NewRequest(httpType, url, bytes.NewBuffer(jsonbyte))
	if err != nil {
		return "", http.Header{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(HTTPUserName, HTTPPassword)
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

// ReadJSONString 读json数据
func ReadJSONString(jsonData string) (float64, error) {
	var v interface{}
	if err := json.Unmarshal([]byte(jsonData), &v); err != nil {
		return 0, errors.New("解析失败" + err.Error())
	}
	data1, ok := v.(map[string]interface{})
	if !ok {
		return 0, errors.New("返回数据无法解析")
	}
	data2, ok := data1["hits"].(map[string]interface{})
	if !ok {
		return 0, errors.New("返回数据无法解析")
	}
	total, ok := data2["total"].(float64)
	if !ok {
		return 0, errors.New("返回数据无法解析")
	}
	data3, ok := data2["hits"].([]interface{})
	if !ok {
		return 0, errors.New("返回数据无法解析")
	}
	for _, v := range data3 {
		data4, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		data, ok := data4["_source"].(map[string]interface{})
		if !ok {
			continue
		}
		ipAndPort := fmt.Sprintf("%v", data["serverip"])
		serveripAndport := strings.Split(ipAndPort, ":")
		if len(serveripAndport) != 2 {
			return 0, errors.New("获取json异常" + ipAndPort)
		}
		_, ok = allInfo[ipAndPort]
		if !ok {
			InitTask(serveripAndport[0], serveripAndport[1])
		}
		SetTask(ipAndPort, fmt.Sprintf("%v", data["task"]), fmt.Sprintf("%v", data["status"]))
	}
	return total, nil
}

// GetInfo 获取info数据并解析处理
func GetInfo(body []byte) error {
	bodystr, _, err := httpJSON(GetURL(), "GET", body)
	if err != nil {
		return errors.New("获取json数据异常： " + bodystr + err.Error())
	}
	_, err = ReadJSONString(bodystr)
	if err != nil {
		return errors.New("获取json数据异常： " + bodystr + err.Error())
	}
	return nil
}

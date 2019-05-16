package CallCoco

import (
	"../File"
	"../Global"
	"net/http"
	"strings"
)

func httpPostJson(jsonStr string,port string) error {

	url := Global.CocoUrl
	//req, err := http.NewRequest("POST", url+"?nat="+jsonStr, nil)
	//req.Header.Set("Content-Type", "application/json")

	//client := &http.Client{}
	resp, err := http.Post(url,
		"application/x-www-form-urlencoded",
		strings.NewReader("nat="+jsonStr+"&port="+port))
	if err != nil {
		return err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			File.WriteErrorLog(err.Error())
		}
	}()

	return nil
}
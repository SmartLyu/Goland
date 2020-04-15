package CallPolice

import (
	"../Global"
	"encoding/json"
	"io/ioutil"
	"regexp"
)

func ReadDingdingAtFile(hostname string, defaultMobiles []int) (DingdingAtMobiles []int) {
	var jsons []Global.DingdingAtJson
	DingdingAtMobiles = defaultMobiles

	data, err := ioutil.ReadFile(Global.DingdingAtFile)
	if err != nil {
		CallMessage(Global.DingdingAtFile + "json file can not read , error is " + err.Error())
		return
	}

	if err := json.Unmarshal(data, &jsons); err != nil {
		CallMessage(Global.DingdingAtFile + "json is error , error is " + err.Error())
		return
	}
	for _, key := range jsons {
		if check, _ := regexp.Match(key.Hostname, []byte(hostname)); check {
			DingdingAtMobiles = append(DingdingAtMobiles, key.Mobiles...)
		}
	}
	return
}

package CallPolice

import (
	"../Global"
	"encoding/json"
	"io/ioutil"
	"regexp"
)

func ReadDingdingAtFile(hostname string) (DingdingAtMobiles []int) {
	var atJsons []Global.DingdingAtJson
	var mobileJsons []Global.DingdingMobilesJson

	data, err := ioutil.ReadFile(Global.DingdingAtFile)
	if err != nil {
		CallMessage(Global.DingdingAtFile + " json file can not read , error is " + err.Error())
		return
	}

	if err := json.Unmarshal(data, &atJsons); err != nil {
		CallMessage(Global.DingdingAtFile + " json is error , error is " + err.Error())
		return
	}

	data, err = ioutil.ReadFile(Global.DingdingMobilesFile)
	if err != nil {
		CallMessage(Global.DingdingMobilesFile + " json file can not read , error is " + err.Error())
		return
	}

	if err := json.Unmarshal(data, &mobileJsons); err != nil {
		CallMessage(Global.DingdingMobilesFile + " json is error , error is " + err.Error())
		return
	}

	for _, key := range atJsons {
		if check, _ := regexp.Match(key.Hostname, []byte(hostname)); check {
			for _, m := range mobileJsons {
				for _, k := range key.Members {
					if m.Member == k {
						DingdingAtMobiles = append(DingdingAtMobiles, m.Mobiles...)
					}
				}
			}
		}
	}

	if DingdingAtMobiles == nil {
		for _, m := range mobileJsons {
			if m.Member == Global.DingdingDefaultAt {
				DingdingAtMobiles = m.Mobiles
			}
		}
	}
	return
}

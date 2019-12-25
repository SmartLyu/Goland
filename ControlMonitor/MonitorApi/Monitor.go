package MonitorApi

import (
	"../Log"
	"encoding/json"
	"errors"
	"strconv"
	"time"
)

func Monitor(agent string, metric string, value float32) error {
	var addJson = make([]PostDataJson, 1)
	addJson[0] = PostDataJson{
		Dagent:     agent,
		Dmetric:    metric,
		Dvalue:     value,
		Dtimestamp: time.Now().Unix(),
	}
	js, _ := json.Marshal(&addJson)
	Log.DebugLog.Println("Add Agents: " + string(js))

	var body string
	var err error
	if body, _, err = httpJson("http://"+URL+"/v1/data-report/", "POST", js); err != nil {
		return errors.New(err.Error())
	}
	if ReadAddDataJsonString(body) {
		Log.InfoLog.Println("Add monitor data successfully, agent: " + agent + ", metric: " + metric + ", value: " +
			strconv.FormatFloat(float64(value), 'f', 0, 32))
	}
	return nil
}

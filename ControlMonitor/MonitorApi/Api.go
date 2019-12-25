package MonitorApi

import (
	"../Log"
	"encoding/json"
	"errors"
	"time"
)

func GetProject(name string) (string, error) {
	var body string
	var err error
	if body, _, err = httpTokenJson("http://120.92.156.253:8888/partition/projects/?query=name:"+name+
		"&query-type=0", "GET", []byte("")); err != nil {
		return "", errors.New(err.Error())

	}
	//fmt.Println(body)
	code, datamap, err := ReadUUidJsonString(body)
	if err != nil {
		return "", err
	}
	if code != "200" || len(datamap) == 0 {
		return "", errors.New("get project error: " + body)
	}
	if len(datamap) > 1 {
		return "", errors.New("获取太多project信息：" + body)
	}
	return datamap[0], nil
}

func GetAgents(projectUuid string, name string) (string, string, error) {
	var body string
	var err error
	if body, _, err = httpTokenJson("http://"+URL+"/v1/agents?query=project:"+projectUuid+",name:"+name+
		"&query-type=0", "GET", []byte("")); err != nil {
		return "", "", err
	}
	code, datamap, err := ReadUUidJsonString(body)
	if err != nil {
		return "", "", err
	}
	if code != "200" {
		return code, body, nil
	}
	if len(datamap) == 0 {
		return "200200", body, nil
	}
	if len(datamap) > 1 {
		return "", "", errors.New("获取太多agents的信息：" + body)
	}
	return code, datamap[0], nil
}

func AddAgents(projectUuid string, name string) (string, error) {
	var addJson = AddProjectJson{
		name,
		projectUuid,
	}
	js, _ := json.Marshal(&addJson)
	Log.DebugLog.Println("Add Agents: " + string(js))

	var body string
	var err error
	if body, _, err = httpJson("http://"+URL+"/v1/agents", "POST", js); err != nil {
		return "", errors.New(err.Error())
	}
	_, datamap, err := ReadAddAgentsJsonString(body)
	if err != nil {
		return "", err
	}
	return datamap, nil
}

func PutAgents(isMaintain bool, agentId string) error {
	var putJson PutAgentJson
	if isMaintain {
		putJson = PutAgentJson{
			"0001-01-01T00:00:00Z",
			time.Now().Add(1 * time.Hour).Format("2006-01-02T15:04:05Z"),
		}
	} else {
		putJson = PutAgentJson{
			"0001-01-01T00:00:00Z",
			"0001-01-01T00:00:00Z",
		}
	}

	js, _ := json.Marshal(&putJson)
	Log.DebugLog.Println("PUT Agents: " + string(js))

	var body string
	var err error
	if body, _, err = httpTokenJson("http://"+URL+"/v1/agents/"+agentId,
		"PUT", js); err != nil {
		return err
	}
	if ReadAddDataJsonString(body) {
		return nil
	} else {
		return errors.New("Put change maintain time errror")
	}
}

func AddStrategy(addJson AddStrategyJson) (string, error) {
	js, _ := json.Marshal(&addJson)
	Log.DebugLog.Println("Add Strategy: " + string(js))

	var body string
	var err error
	if body, _, err = httpJson("http://"+URL+"/v1/strategies", "POST", js); err != nil {
		return "", err
	}
	_, datamap, err := ReadAddAgentsJsonString(body)
	if err != nil {
		return "", err
	}
	return datamap, nil
}

func GetStrategy(agentUuid string, metric string) (string, string, error) {
	var body string
	var err error
	if body, _, err = httpTokenJson("http://"+URL+"/v1/strategies?query=agent:"+agentUuid+
		",metric:"+metric+"&query-type=0", "GET", []byte("")); err != nil {
		return "", "", err
	}
	code, datamap, err := ReadUUidJsonString(body)
	if err != nil {
		return "", "", err
	}
	if code != "200" {
		return code, body, nil
	}
	if len(datamap) == 0 {
		return "200200", body, nil
	}
	if len(datamap) > 1 {
		return "", "", errors.New("获取太多Strategy的信息：" + body)
	}
	return code, datamap[0], nil
}

func GetRuleTemplete(rule_name string) (string, error) {
	var body string
	var err error
	if body, _, err = httpTokenJson("http://"+URL+"/v1/ruletpls?query=name:"+rule_name+
		"&query-type=0", "GET", []byte("")); err != nil {
		return "", err
	}
	code, datamap, err := ReadUUidJsonString(body)
	if err != nil {
		return "", err
	}
	if code != "200" || len(datamap) == 0 {
		return "", errors.New("get rule templete error: " + body)
	}
	if len(datamap) > 1 {
		return "", errors.New("获取太多RuleTemplete的信息：" + body)
	}
	return datamap[0], nil
}

func AddRule(ruletpl string, rule_name string, projects []string) (string, error) {
	var addJson = AddRoleJson{
		ruletpl,
		rule_name,
		projects,
	}
	js, _ := json.Marshal(&addJson)
	Log.DebugLog.Println("Add Rule: " + string(js))

	var body string
	var err error
	if body, _, err = httpJson("http://"+URL+"/v1/rule-based-tpl", "POST", js); err != nil {
		return "", err
	}
	returnStr, err := ReadAddRuleJsonString(body)
	if err != nil {
		return "", err
	}
	return returnStr, nil
}

func GetRule(rule_name string, project string) (string, string, error) {
	var body string
	var err error
	if body, _, err = httpTokenJson("http://"+URL+"/v1/rules?query=name:"+rule_name+",project:"+project+
		"&query-type=0", "GET", []byte("")); err != nil {
		return "", "", err
	}
	code, datamap, err := ReadUUidJsonString(body)
	if err != nil {
		return "", "", err
	}
	if code != "200" {
		return code, body, nil
	}
	if len(datamap) == 0 {
		return "200200", body, nil
	}
	if len(datamap) > 1 {
		return "", "", errors.New("获取太多Rule的信息：" + body)
	}
	return code, datamap[0], nil
}

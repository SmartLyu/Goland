package MonitorApi

import (
	"../Log"
)

// 监控2.0使用的模板
func DoDefault(projectst string, placestr string, isMaintain bool) (string, string, error) {
	projectName = projectst
	place = placestr

	if place == "海外" {
		URL = OverseaUrl
	} else {
		URL = ChineseUrl
	}

	AgentsName = "单包行为监控"
	project, err := GetProject(projectName)
	if err != nil {
		return "", "", err
	}
	Log.InfoLog.Println("Successfully get project \"" + projectName + "\" uuid: " + project)

	code, agent, err := GetAgents(project, AgentsName)
	if err != nil {
		return "", "", err
	}
	if code != "200" {
		agent, err = AddAgents(project, AgentsName)
		if err != nil {
			return "", "", err
		}
	}

	err = PutAgents(isMaintain, agent)
	if err != nil {
		return "", "", err
	}
	Log.InfoLog.Println("Successfully get agent \"" + AgentsName + "\" uuid: " + agent)

	if isMaintain {
		return project, agent, nil
	}

	strategy = AddStrategyJson{
		Smetric:          "",
		Sop:              "=",
		Sfunc:            "all(#1)",
		Sright_value:     1,
		Snote:            "同一策略模拟失败后，触发30份同策略行为模拟，失败率达到20%",
		Sagent:           agent,
		Snodata:          false,
		Snodata_value:    1,
		Snodata_interval: 1800,
		Sproject:         project,
	}

	strategySmetrics := [...]string{
		"启动游戏",
		"登录",
		"加载服务器成功",
		"进入游戏",
		"打开充值界面",
		"充值行为",
	}
	for _, s := range strategySmetrics {
		if s == "启动游戏" {
			strategy.Snodata = true
		} else {
			strategy.Snodata = false
		}
		strategy.Smetric = s
		code, strategyUUid, err := GetStrategy(agent, strategy.Smetric)
		if err != nil {
			return "", "", err
		}
		if code != "200" {
			strategyUUid, err = AddStrategy(strategy)
			if err != nil {
				return "", "", err
			}
		}
		Log.DebugLog.Println("Successfully get strategy metric \"" + strategy.Smetric + "\" uuid: " + strategyUUid)
	}

	rules := [...]string{
		"游戏登录",
		"游戏选服页面",
		"游戏在线",
		"游戏充值界面监控",
		"游戏获取支付监控",
	}

	for _, r := range rules {
		ruleName = r
		ruletpl = r
		ruleTmp, err := GetRuleTemplete(ruletpl)
		if err != nil {
			return "", "", err
		}
		Log.DebugLog.Println("Successfully get rule templete \"" + ruletpl + "\" uuid: " + ruleTmp)

		code, rule, err := GetRule(ruleName, project)
		if err != nil {
			return "", "", err
		}
		if code != "200" {
			var projects = make([]string, 1, 2)
			projects[0] = project
			rule, err = AddRule(ruleTmp, ruleName, projects)
			if err != nil {
				return "", "", err
			}
		}
		Log.DebugLog.Println("Successfully get rule \"" + ruleName + "\" uuid: " + rule)
	}

	return project, agent, nil
}

package MonitorApi

var (
	// 项目名必填
	projectName = "测试"
	URL         = ""
	place       = ""

	// agent(数据源)、strategy(策略)
	AgentsName = "于志远测试agent"
	strategy   = AddStrategyJson{
		Smetric:          "测试指标",
		Sop:              "<",
		Sfunc:            "all(#2)",
		Sright_value:     3,
		Snote:            "测试备注",
		Sagent:           "",
		Snodata:          true,
		Snodata_value:    1,
		Snodata_interval: 5,
		Sproject:         "",
	}

	// rule规则, ruletpl规则模板需要事先配好
	ruleName = "测试规则"
	ruletpl  = "测试"
)

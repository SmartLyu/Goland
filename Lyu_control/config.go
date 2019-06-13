package main

var (
	DataFileDir    = "/work/logs/control/"             // 存放历史监控信息
	DataFileName   = ".control.log"                    // 监控后缀名
	ListenSig      = make(chan int)                    // 监听后台阻塞信号
	CrtFile        = "/work/ssl/ijunhai.com.pem"       // https的pem文件
	KeyFile        = "/work/ssl/ijunhai.com.key"       // https的key文件
	ShellFile      = "/work/sh/LyuControl.sh"          // 微信连接的脚本
	ShowFile       = "/work/data/control/showfile"     // 展示反馈内容文件
	LockFile       = "/work/data/control/control.lock" // 检查命令执行情况
	WaiteTime      = 7                                 // 最大等待时间70s
	NatShellFile   = "/work/sh/NatShell.sh"            // nat执行的脚本


)

var SecretId = corpText{

	touser:     "YuZhiYuan",
	agentid:    1000002,
	content:    "Test Api",
	safe:       "1",
}

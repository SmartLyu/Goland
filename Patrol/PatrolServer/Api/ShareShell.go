package Api

import (
	"../Global"
	"io/ioutil"
	"net/http"
	"os"
)

// 返回监控脚本
func ReturnMonitorShell(w http.ResponseWriter, r *http.Request) {
	//解析参数，默认是不会解析的
	err := r.ParseForm()
	if err != nil {
		Global.ErrorLog.Println("Recv:" + r.RemoteAddr)
	}

	des := Global.MonitorShellFile
	desStat, err := os.Stat(des)
	if err != nil {
		Global.ErrorLog.Println("File Not Exit " + des)
		http.NotFoundHandler().ServeHTTP(w, r)
	} else if desStat.IsDir() {
		Global.ErrorLog.Println("File Is Dir" + des)
		http.NotFoundHandler().ServeHTTP(w, r)
	} else {
		data, err := ioutil.ReadFile(des)
		if err != nil {
			Global.ErrorLog.Println("Read File Err: " + err.Error())
		} else {
			Global.InfoLog.Println("Send File:" + des)
			_, err = w.Write([]byte(data))
			if err != nil {
				Global.ErrorLog.Println("http writed Err " + err.Error())
			}
		}
	}
}

// 返回nat机器提权脚本
func ReturnNatShell(w http.ResponseWriter, r *http.Request) {
	//解析参数，默认是不会解析的
	err := r.ParseForm()
	if err != nil {
		Global.ErrorLog.Println("Recv:" + r.RemoteAddr)
	}

	des := Global.NatShellFile
	desStat, err := os.Stat(des)
	if err != nil {
		Global.ErrorLog.Println("File Not Exit " + des)
		http.NotFoundHandler().ServeHTTP(w, r)
	} else if desStat.IsDir() {
		Global.ErrorLog.Println("File Is Dir" + des)
		http.NotFoundHandler().ServeHTTP(w, r)
	} else {
		data, err := ioutil.ReadFile(des)
		if err != nil {
			Global.ErrorLog.Println("Read File Err: " + err.Error())
		} else {
			Global.InfoLog.Println("Send File: " + des)
			_, err = w.Write([]byte(data))
			if err != nil {
				Global.ErrorLog.Println("http writed Err " + err.Error())
			}
		}
	}
}

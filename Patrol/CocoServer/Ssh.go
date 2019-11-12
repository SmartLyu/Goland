package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

func connect(user, password, host, key string, port int, cipherList []string) (*ssh.Session, *ssh.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		config       ssh.Config
		session      *ssh.Session
		err          error
	)
	// get auth method , 用秘钥或者密码连接
	auth = make([]ssh.AuthMethod, 0)
	if key == "" {
		auth = append(auth, ssh.Password(password))
	} else {
		pemBytes, err := ioutil.ReadFile(key)
		if err != nil {
			return nil, nil, err
		}

		var signer ssh.Signer
		if password == "" {
			signer, err = ssh.ParsePrivateKey(pemBytes)
		} else {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(pemBytes, []byte(password))
		}
		if err != nil {
			return nil, nil, err
		}
		auth = append(auth, ssh.PublicKeys(signer))
	}

	if len(cipherList) == 0 {
		config = ssh.Config{
			Ciphers: []string{"aes128-ctr", "aes192-ctr",
				"aes256-ctr", "aes128-gcm@openssh.com", "arcfour256",
				"arcfour128", "aes128-cbc", "3des-cbc", "aes192-cbc", "aes256-cbc"},
		}
	} else {
		config = ssh.Config{
			Ciphers: cipherList,
		}
	}

	clientConfig = &ssh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: 10 * time.Second,
		Config:  config,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	// ssh 连接
	addr = fmt.Sprintf("%s:%d", host, port)

	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, nil, err
	}

	// create session
	if session, err = client.NewSession(); err != nil {
		return nil, nil, err
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		return nil, nil, err
	}

	return session, client, nil
}

var (
	username = SshUser
	password = ""
	key      = Sshkey
)

func sshDoShell(ip string, port int, cmd string) error {
	ciphers := []string{}
	session, client, err := connect(username, password, ip, key, port, ciphers)

	defer func() {
		_ = session.Close()
		errc := client.Close()
		if errc != nil {
			log.Println("close client has error " + errc.Error())
		}
	}()

	if err != nil {
		return errors.New("连接 " + ip + " 异常" + err.Error())
	}

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	err = session.Run(cmd)
	if err != nil {
		log.Println("error has stopped " + err.Error())
		return errors.New(err.Error())
	}

	return nil
}

func SshToNat(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		WriteLog("Recv:" + r.RemoteAddr)
	}

	// 读取用户输入的参数信息
	getNat := r.Form.Get("nat")
	getPort, err := strconv.Atoi(r.Form.Get("port"))

	if getNat == "" || err != nil {
		w.WriteHeader(http.StatusFailedDependency)
		return
	}

	err = sshDoShell(getNat, getPort, "wget -q "+NatShellDownloadUrl+" --timeout 10 -O /tmp/patrol-tmp.sh&&"+
		"/usr/bin/nohup /bin/bash /tmp/patrol-tmp.sh --nat "+getNat+" &> /dev/null &&"+
		"rm -f /tmp/patrol-tmp.sh")
	if err != nil {
		WriteLog("ssh error:" + r.RemoteAddr)
		body, _, err := httpPostJson("ssh-to-nat-"+getNat, "false")
		if err != nil {
			WriteLog("post error:" + r.RemoteAddr + " - " + err.Error())
		}
		WriteLog(body)
	}

	// 返回信息
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(getNat); err != nil {
		WriteLog(err.Error())
	}

	body, _, err := httpPostJson("ssh-to-nat-"+getNat, "true")
	if err != nil {
		WriteLog("post error:" + r.RemoteAddr + " - " + err.Error())
	}
	WriteLog(body)
	WriteLog(time.Now().Format("2006.01.02 15:04") + "\t ssh " + getNat + " successfully")
	return
}

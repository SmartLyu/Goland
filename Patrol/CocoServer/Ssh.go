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

func connect(user, password, host, key string, port int, cipherList []string) (*ssh.Session, error) {
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
			return nil, err
		}

		var signer ssh.Signer
		if password == "" {
			signer, err = ssh.ParsePrivateKey(pemBytes)
		} else {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(pemBytes, []byte(password))
		}
		if err != nil {
			return nil, err
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
		Timeout: 30 * time.Second,
		Config:  config,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	// ssh 连接
	addr = fmt.Sprintf("%s:%d", host, port)

	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	// create session
	if session, err = client.NewSession(); err != nil {
		return nil, err
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		return nil, err
	}

	return session, nil
}

var (
	username = SshUser
	password = ""
	key      = Sshkey
)

func sshDoShell(ip string, port int, cmd string) error {
	ciphers := []string{}
	session, err := connect(username, password, ip, key, port, ciphers)

	if err != nil {
		fmt.Println("连接 ", ip, " 异常")
		log.Panic(err)
	}

	defer func() {
		_ = session.Close()
	}()

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	err = session.Run(cmd)
	if err != nil {
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
	getPort , err:= strconv.Atoi(r.Form.Get("port"))

	fmt.Println(getNat)
	fmt.Println(getPort)

	if getNat == "" || err != nil {
		w.WriteHeader(http.StatusFailedDependency)
	}

	err = sshDoShell(getNat, getPort, "wget "+NatShellDownloadUrl+" --timeout 10 -O /tmp/patrol-tmp.sh;"+
		"/usr/bin/nohup /bin/bash /tmp/patrol-tmp.sh --nat "+getNat+";"+
		"rm -f /tmp/patrol-tmp.sh")
	if err != nil {
		WriteLog("ssh error:" + r.RemoteAddr)
		body, _, err := httpPostJson("ssh error:", "false")
		if err != nil {
			WriteLog("post error:" + r.RemoteAddr)
		}
		WriteLog(body)
	}

	// 返回信息
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(getNat); err != nil {
		WriteLog(err.Error())
	}
}

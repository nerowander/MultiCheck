package Plugins

import (
	"InformationScan/common"
	"InformationScan/config"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net"
	"strings"
	"time"
)

func sshScan(info *config.InfoScan) (reserr error) {
	if !info.Brute {
		return
	}
	starttime := time.Now().Unix()
	for _, user := range config.UsernameDict["ssh"] {
		for _, pass := range config.PasswordDict {
			pass = strings.Replace(pass, "{user}", user, -1)
			flag, err := SshConn(info, user, pass)
			if flag == true && err == nil {
				return err
			} else {
				errlog := fmt.Sprintf("[-] ssh %v:%v %v %v %v", info.Hosts, info.Ports, user, pass, err)
				common.LogError(errlog)
				reserr = err
				if common.CheckErrMessages(err) {
					return err
				}
				if time.Now().Unix()-starttime > (int64(len(config.UsernameDict["ssh"])*len(config.PasswordDict)) * config.Timeout) {
					return err
				}
			}
			if info.SshKey != "" {
				return err
			}
		}
	}
	return reserr
}

func SshConn(info *config.InfoScan, user string, pass string) (flag bool, err error) {
	flag = false
	Host, Port, Username, Password := info.Hosts, info.Ports, user, pass
	var Auth []ssh.AuthMethod
	if info.SshKey != "" {
		var pemBytes []byte
		pemBytes, err = ioutil.ReadFile(info.SshKey)
		if err != nil {
			return false, errors.New("read key failed" + err.Error())
		}
		signer, err1 := ssh.ParsePrivateKey(pemBytes)
		if err1 != nil {
			return false, errors.New("parse key failed" + err.Error())
		}
		Auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}
	} else {
		Auth = []ssh.AuthMethod{ssh.Password(Password)}
	}

	sshConfig := &ssh.ClientConfig{
		User:    Username,
		Auth:    Auth,
		Timeout: time.Duration(config.Timeout) * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%v:%v", Host, Port), sshConfig)
	if err == nil {
		defer client.Close()
		var session *ssh.Session
		session, err = client.NewSession()
		if err == nil {
			defer session.Close()
			flag = true
			var logMessage string
			if config.Command != "" {
				combo, _ := session.CombinedOutput(config.Command)
				logMessage = fmt.Sprintf("[+] SSH %v:%v:%v %v \n %v", Host, Port, Username, Password, string(combo))
				if info.SshKey != "" {
					logMessage = fmt.Sprintf("[+] SSH %v:%v sshkey correct \n %v", Host, Port, string(combo))
				}
				common.LogSuccess(logMessage)
			} else {
				logMessage = fmt.Sprintf("[+] SSH %v:%v:%v %v", Host, Port, Username, Password)
				if info.SshKey != "" {
					logMessage = fmt.Sprintf("[+] SSH %v:%v sshkey correct", Host, Port)
				}
				common.LogSuccess(logMessage)
			}
		}
	}
	return flag, err

}

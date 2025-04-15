package Plugins

import (
	"InformationScan/common"
	"InformationScan/config"
	"errors"
	"fmt"
	"github.com/stacktitan/smb/smb"
	"strings"
	"time"
)

func smbScan(info *config.InfoScan) (reserr error) {
	if !info.Brute {
		return nil
	}
	starttime := time.Now().Unix()
	for _, user := range config.UsernameDict["smb"] {
		for _, pass := range config.PasswordDict {
			pass = strings.Replace(pass, "{user}", user, -1)
			flag, err := doWithTimeOut(info, user, pass)
			if flag == true && err == nil {
				var result string
				if info.Domain != "" {
					result = fmt.Sprintf("[+] SMB %v:%v:%v\\%v %v", info.Hosts, info.Ports, info.Domain, user, pass)
				} else {
					result = fmt.Sprintf("[+] SMB %v:%v:%v %v", info.Hosts, info.Ports, user, pass)
				}
				common.LogSuccess(result)
				return err
			} else {
				errlog := fmt.Sprintf("[-] smb %v:%v %v %v %v", info.Hosts, 445, user, pass, err)
				errlog = strings.Replace(errlog, "\n", "", -1)
				common.LogError(errlog)
				reserr = err
				if common.CheckErrMessages(err) {
					return err
				}
				if time.Now().Unix()-starttime > (int64(len(config.UsernameDict["smb"])*len(config.PasswordDict)) * config.Timeout) {
					return err
				}
			}
		}
	}
	return reserr
}

func SmbConn(info *config.InfoScan, user string, pass string, signal chan struct{}) (flag bool, err error) {
	flag = false
	Host, Username, Password := info.Hosts, user, pass
	options := smb.Options{
		Host:        Host,
		Port:        445,
		User:        Username,
		Password:    Password,
		Domain:      info.Domain,
		Workstation: "",
	}

	session, err := smb.NewSession(options, false)
	if err == nil {
		session.Close()
		if session.IsAuthenticated {
			flag = true
		}
	}
	signal <- struct{}{}
	return flag, err
}

func doWithTimeOut(info *config.InfoScan, user string, pass string) (flag bool, err error) {
	signal := make(chan struct{})
	go func() {
		flag, err = SmbConn(info, user, pass, signal)
	}()
	select {
	case <-signal:
		return flag, err
	case <-time.After(time.Duration(config.Timeout) * time.Second):
		return false, errors.New("smb connection time out")
	}
}

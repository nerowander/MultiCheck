package Plugins

import (
	"bytes"
	"fmt"
	"github.com/jlaffaye/ftp"
	"github.com/nerowander/MultiCheck/common"
	"github.com/nerowander/MultiCheck/config"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func ftpScan(info *config.InfoScan) (reserr error) {
	if info.Brute == false {
		return
	}
	startTime := time.Now().Unix()
	res, err := ftpConnect(info, "anonymous", "")
	if res && err == nil {
		return err
	} else {
		errlog := fmt.Sprintf("[-] ftp %v:%v %v %v", info.Hosts, info.Ports, "anonymous", err)
		common.LogError(errlog)
		reserr = err
		if common.CheckErrMessages(err) {
			return err
		}
	}
	// Brute
	for _, user := range config.UsernameDict["ftp"] {
		for _, pass := range config.PasswordDict {
			pass = strings.Replace(pass, "{user}", user, -1)
			flag, err := ftpConnect(info, user, pass)
			if flag && err == nil {
				return err
			} else {
				errlog := fmt.Sprintf("[-] ftp %v:%v %v %v %v", info.Hosts, info.Ports, user, pass, err)
				common.LogError(errlog)
				reserr = err
				if common.CheckErrMessages(err) {
					return err
				}
				// timeout
				if time.Now().Unix()-startTime > (int64(len(config.UsernameDict["ftp"])*len(config.PasswordDict)) * config.Timeout) {
					return err
				}
			}
		}
	}
	return reserr
}

func ftpConnect(info *config.InfoScan, user, pass string) (res bool, err error) {
	res = false
	host, port, username, password := info.Hosts, info.Ports, user, pass
	var conn *ftp.ServerConn
	// 无认证
	conn, err = ftp.Dial(host+":"+port, ftp.DialWithTimeout(time.Duration(config.Timeout)*time.Second))
	if err != nil {
		// 需要认证
		err = conn.Login(username, password)
		// 认证成功，执行操作
		if err == nil {
			res = true
			logMessage := fmt.Sprintf("[+] ftp connect success: %v:%v:%v %v", host, port, username, password)
			// store a file
			if info.FTPReadFile != "" {
				r, err1 := conn.Retr(info.FTPReadFile)
				if err1 != nil {
					logMessage += "[+] the chosen file to read doesn't exists or read file failed"
				}
				defer r.Close()
				buf, _ := ioutil.ReadAll(r)
				logMessage += fmt.Sprintf("[+] ftp read file %v : %v", info.FTPWriteFile, string(buf))
			}
			// Read a file
			if info.FTPWriteFile != "" {
				data, err2 := os.ReadFile(info.FTPWriteFile)
				if err2 != nil {
					logMessage += fmt.Sprintf("[+] the chosen file to write doesn't exists in your filesystem")
				} else {
					err3 := conn.Stor(info.FTPWriteFile, bytes.NewBufferString(string(data)))
					if err3 != nil {
						logMessage += fmt.Sprintf("[+] ftp write file %v failed", info.FTPWriteFile)
					}
				}
			}
			// read dir file
			dirs, err := conn.List("")
			if err == nil {
				if len(dirs) > 0 {
					for i := 0; i < len(dirs); i++ {
						if len(dirs[i].Name) > 50 {
							logMessage += "\n   [->]" + dirs[i].Name[:50]
						} else {
							logMessage += "\n   [->]" + dirs[i].Name
						}
						if i == 5 {
							break
						}
					}
				}
			}
		}
	}
	// common.LogSuccess(logMessage)
	// ftp.DialWithTimeout()
	return res, err
}

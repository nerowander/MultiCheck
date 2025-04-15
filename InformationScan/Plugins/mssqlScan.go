package Plugins

import (
	"database/sql"
	"fmt"
	"github.com/nerowander/MultiCheck/common"
	"github.com/nerowander/MultiCheck/config"
	"strings"
	"time"
)

func mssqlScan(info *config.InfoScan) (reserr error) {
	if info.Brute == false {
		return
	}
	starttime := time.Now().Unix()
	for _, user := range config.UsernameDict["mssql"] {
		for _, pass := range config.PasswordDict {
			pass = strings.Replace(pass, "{user}", user, -1)
			flag, err := MssqlConn(info, user, pass)
			if flag == true && err == nil {
				return err
			} else {
				errlog := fmt.Sprintf("[-] mssql %v:%v %v %v %v", info.Hosts, info.Ports, user, pass, err)
				common.LogError(errlog)
				reserr = err
				if common.CheckErrMessages(err) {
					return err
				}
				if time.Now().Unix()-starttime > (int64(len(config.UsernameDict["mssql"])*len(config.PasswordDict)) * config.Timeout) {
					return err
				}
			}
		}
	}
	return reserr
}

func MssqlConn(info *config.InfoScan, user string, pass string) (bool, error) {
	flag := false
	Host, Port, Username, Password := info.Hosts, info.Ports, user, pass
	dataSourceName := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%v;encrypt=disable;timeout=%v", Host, Username, Password, Port, time.Duration(config.Timeout)*time.Second)
	db, err := sql.Open("mssql", dataSourceName)
	if err == nil {
		db.SetConnMaxLifetime(time.Duration(config.Timeout) * time.Second)
		db.SetConnMaxIdleTime(time.Duration(config.Timeout) * time.Second)
		db.SetMaxIdleConns(0)
		defer db.Close()
		err = db.Ping()
		if err == nil {
			result := fmt.Sprintf("[+] mssql %v:%v:%v %v", Host, Port, Username, Password)
			common.LogSuccess(result)
			flag = true
		}
	}
	return flag, err
}

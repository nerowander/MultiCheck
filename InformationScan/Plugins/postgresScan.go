package Plugins

import (
	"InformationScan/common"
	"InformationScan/config"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

func postgresScan(info *config.InfoScan) (reserr error) {
	if info.Brute == false {
		return
	}
	starttime := time.Now().Unix()
	for _, user := range config.UsernameDict["postgres"] {
		for _, pass := range config.PasswordDict {
			pass = strings.Replace(pass, "{user}", user, -1)
			flag, err := PostgresConn(info, user, pass)
			if flag == true && err == nil {
				return err
			} else {
				errlog := fmt.Sprintf("[-] postgres %v:%v %v %v %v", info.Hosts, info.Ports, user, pass, err)
				common.LogError(errlog)
				reserr = err
				if common.CheckErrMessages(err) {
					return err
				}
				if time.Now().Unix()-starttime > (int64(len(config.UsernameDict["postgres"])*len(config.PasswordDict)) * config.Timeout) {
					return err
				}
			}
		}
	}
	return reserr
}

func PostgresConn(info *config.InfoScan, user string, pass string) (bool, error) {
	flag := false
	Host, Port, Username, Password := info.Hosts, info.Ports, user, pass
	dataSourceName := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=%v", Username, Password, Host, Port, "postgres", "disable")
	db, err := sql.Open("postgres", dataSourceName)
	if err == nil {
		db.SetConnMaxLifetime(time.Duration(config.Timeout) * time.Second)
		db.SetConnMaxIdleTime(time.Duration(config.Timeout) * time.Second)
		db.SetMaxIdleConns(0)
		defer db.Close()
		err = db.Ping()
		if err == nil {
			result := fmt.Sprintf("[+] postgres %v:%v:%v %v", Host, Port, Username, Password)
			common.LogSuccess(result)
			flag = true
		}
	}
	return flag, err
}

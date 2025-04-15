package Plugins

import (
	"InformationScan/common"
	"InformationScan/config"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"time"
)

func mysqlScan(info *config.InfoScan) (reserr error) {
	if info.Brute == false {
		return
	}
	starttime := time.Now().Unix()
	for _, user := range config.UsernameDict["mysql"] {
		for _, pass := range config.PasswordDict {
			pass = strings.Replace(pass, "{user}", user, -1)
			flag, err := MysqlConn(info, user, pass)
			if flag == true && err == nil {
				return err
			} else {
				errlog := fmt.Sprintf("[-] mysql %v:%v %v %v %v", info.Hosts, info.Ports, user, pass, err)
				common.LogError(errlog)
				reserr = err
				if common.CheckErrMessages(err) {
					return err
				}
				if time.Now().Unix()-starttime > (int64(len(config.UsernameDict["mysql"])*len(config.PasswordDict)) * config.Timeout) {
					return err
				}
			}
		}
	}
	return reserr
}

func MysqlConn(info *config.InfoScan, user string, pass string) (bool, error) {
	flag := false
	Host, Port, Username, Password := info.Hosts, info.Ports, user, pass
	dataSourceName := fmt.Sprintf("%v:%v@tcp(%v:%v)/mysql?charset=utf8&timeout=%v", Username, Password, Host, Port, time.Duration(config.Timeout)*time.Second)
	db, err := sql.Open("mysql", dataSourceName)
	if err == nil {
		db.SetConnMaxLifetime(time.Duration(config.Timeout) * time.Second)
		db.SetConnMaxIdleTime(time.Duration(config.Timeout) * time.Second)
		db.SetMaxIdleConns(0)
		defer db.Close()
		err = db.Ping()
		if err == nil {
			result := fmt.Sprintf("[+] mysql %v:%v:%v %v", Host, Port, Username, Password)
			common.LogSuccess(result)
			flag = true
		}
	}
	return flag, err
}

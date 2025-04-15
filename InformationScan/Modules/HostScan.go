package Modules

import (
	"FinalProject/config"
	"InformationScan/Plugins"
	"InformationScan/WebScan/lib"
	"InformationScan/common"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

func HostScan(info *config.InfoScan) {
	// 判断是否是容器模式，如果是则传参数到另一个函数，请求InformationScan module容器
	if config.EnableInfoContainer == true {
		//config.EnablePocContainer = true
		module := "InformationScan"
		url := fmt.Sprintf("http://%s-service/infoscan?hosts=%s&brute=%s&ftpreadfile=%s&ftpwritefile=%s&sshkey=%s&domain=%s&"+
			"skipredis=%s&redissshfile=%s&rediswebshellfile=%s&rediscronhost=%s&remotepath=%s&logwaittime=%s&printlog=%s"+
			"savelogtojson=%s&savelogtohtml=%s", module, info.Hosts, strconv.FormatBool(info.Brute), info.FTPReadFile, info.FTPWriteFile,
			info.SshKey, info.Domain, strconv.FormatBool(info.SkipRedis), info.RedisSshFile, info.RedisWebshellFile, info.RedisCronHost,
			info.RemotePath, strconv.FormatInt(common.LogWaitTime, 10), strconv.FormatBool(common.PrintLog), strconv.FormatBool(common.SaveLogToJSON), strconv.FormatBool(common.SaveLogToHTML))

		fmt.Println(url)
		jsonPayload := map[string]interface{}{
			"HostFile":     config.HostFile,
			"Threads":      config.Threads,
			"ScanType":     config.ScanType,
			"Ports":        config.Ports,
			"PortsFile":    config.PortsFile,
			"Timeout":      config.Timeout,
			"URL":          config.URL,
			"URLFile":      config.URLFile,
			"Username":     config.Username,
			"UsernameFile": config.UsernameFile,
			"Password":     config.Password,
			"PasswordFile": config.PasswordFile,
			"HashFile":     config.HashFile,
			"Hash":         config.Hash,
			"AddPorts":     config.AddPorts,
			"AddPassWords": config.AddPassWords,
			"AddUserNames": config.AddUserNames,
			"BruteThreads": config.BruteThreads,
			"Command":      config.Command,
			"Cookie":       config.Cookie,
			"PocNum":       config.PocNum,
			"PocPath":      config.PocPath,
			"PocName":      config.PocName,
			"PocType":      config.PocType,
			"ExpNum":       config.ExpNum,
			"ExpPath":      config.ExpPath,
			"ExpType":      config.ExpType,
			"ExpName":      config.ExpName,
			"WebTimeout":   config.WebTimeout,
			"NoPOC":        config.NoPOC,
			"NoExploit":    config.NoExploit,
			"DnsLog":       config.DnsLog,
			"CeyeToken":    config.CeyeToken,
			"CeyeURL":      config.CeyeURL,
			"FullPOC":      config.FullPOC,
			"FullEXP":      config.FullEXP,
			"SaveResult":   config.SaveResult,
			"OutPutFile":   config.OutPutFile,
		}
		var body []byte
		var err error
		body, err = json.Marshal(jsonPayload)
		if err != nil {
			fmt.Printf("Error marshaling JSON: %v\n", err)
			return
		}
		var resp *http.Response
		resp, err = http.Post(url, "application/json", bytes.NewBuffer(body))
		if err != nil {
			fmt.Printf("Error calling module %s: %v\n", module, err)
			return
		}
		defer resp.Body.Close()
		resBody, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(resBody))
	} else {
		fmt.Println("Start host scan")
		hosts, err := common.ParseIP(info)
		if err != nil {
			fmt.Println("Parse IP error :", err)
		}
		//fmt.Println(config.Threads)
		lib.InitHTTP()
		var ch = make(chan struct{}, config.Threads)
		var wg = sync.WaitGroup{}
		var web = strconv.Itoa(config.PortInt["web"])
		if len(hosts) > 0 || len(info.HostPort) > 0 {
			if len(hosts) > 1 || config.ScanType == "ping" {
				// host
				hosts = CheckHostLive(hosts)
				fmt.Println("[*] Alive hosts len is: ", len(hosts))
			}
			// 扫描类型指定为ping，下面的portscan同理
			if config.ScanType == "ping" {
				config.WG.Wait()
				return
			}
			var AlivePorts []string
			// port
			AlivePorts = CheckPortLive(hosts, config.Ports, config.Timeout)
			fmt.Println("[*] alive ports len is:", len(AlivePorts))
			if config.ScanType == "portscan" {
				config.WG.Wait()
				return
			}

			if len(info.HostPort) > 0 {
				// ip:port
				AlivePorts = append(AlivePorts, info.HostPort...)
				AlivePorts = common.RemoveDuplicateHosts(AlivePorts)
				info.HostPort = nil
				fmt.Println("[*] AliveHostPorts len is: ", len(AlivePorts))
			}
			var serverPorts []string
			for _, serverPort := range config.PortInt {
				serverPorts = append(serverPorts, strconv.Itoa(serverPort))
			}

			for _, target := range AlivePorts {
				info.Hosts, info.Ports = strings.Split(target, ":")[0], strings.Split(target, ":")[1]
				if config.ScanType == "all" {
					// todo: add 135
					// todo: add ms17010
					switch {
					case info.Ports == "9000":
						makeScan(web, info, ch, &wg)        // webtitle
						makeScan(info.Ports, info, ch, &wg) // fastcgi
					case IsContain(serverPorts, info.Ports):
						makeScan(info.Ports, info, ch, &wg) // plugins
					default:
						makeScan(web, info, ch, &wg) // webtitle
					}
				} else {
					scanType := strconv.Itoa(config.PortInt[config.ScanType])
					makeScan(scanType, info, ch, &wg)
				}
			}
		}
		for _, url := range config.Urls {
			info.Url = url
			makeScan(web, info, ch, &wg)
		}
		wg.Wait()
		config.WG.Wait()
		close(common.LogResults)
		fmt.Printf("已完成扫描任务 %v/%v\n", common.End, common.Num)
	}
}

var Mutex = &sync.Mutex{}

func makeScan(scanType string, info *config.InfoScan, ch chan struct{}, wg *sync.WaitGroup) {
	ch <- struct{}{}
	wg.Add(1)
	go func(infoCopy config.InfoScan) {
		defer wg.Done()
		defer func() { <-ch }()
		Mutex.Lock()
		common.Num += 1
		Mutex.Unlock()
		ScanInvoke(&scanType, &infoCopy)
		Mutex.Lock()
		common.End += 1
		Mutex.Unlock()
		//wg.Done()
		//<-ch
	}(*info)
}

func ScanInvoke(name *string, info *config.InfoScan) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("[-] %v:%v scan error: %v\n\n", info.Hosts, info.Ports, err)
		}
	}()
	p := reflect.ValueOf(Plugins.PluginList[*name])
	v := []reflect.Value{reflect.ValueOf(info)}
	p.Call(v)
}

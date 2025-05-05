package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nerowander/MultiCheck/InformationScan/Modules"
	"github.com/nerowander/MultiCheck/common"
	"github.com/nerowander/MultiCheck/config"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

var (
	taskResults sync.Map
	taskWg      sync.WaitGroup
	info        config.InfoScan
)

type ScanConfig struct {
	Username     string   `json:"username"`
	UsernameFile string   `json:"usernamefile"`
	URL          string   `json:"url"`
	URLFile      string   `json:"urlfile"`
	HostFile     string   `json:"hostfile"`
	Ports        string   `json:"ports"`
	PortsFile    string   `json:"portsfile"`
	Hash         string   `json:"hash"`
	HashFile     string   `json:"hashfile"`
	AddPorts     string   `json:"addports"`
	AddUserNames string   `json:"addusernames"`
	AddPassWords string   `json:"addpasswords"`
	Password     string   `json:"password"`
	PasswordFile string   `json:"passwordfile"`
	Threads      int      `json:"threads"`
	ScanType     string   `json:"scantype"`
	Timeout      int64    `json:"timeout"`
	Urls         []string `json:"urls"`
	Command      string   `json:"command"`
	BruteThreads int      `json:"brutethreads"`
	Cookie       string   `json:"cookie"`
	PocNum       int      `json:"pocnum"`
	PocPath      string   `json:"pocpath"`
	PocName      string   `json:"pocname"`
	PocType      string   `json:"poctype"`
	ExpNum       int      `json:"expnum"`
	ExpPath      string   `json:"exppath"`
	ExpType      string   `json:"exptype"`
	ExpName      string   `json:"expname"`
	WebTimeout   int64    `json:"webtimeout"`
	NoPOC        bool     `json:"nopoc"`
	NoExploit    bool     `json:"noexploit"`
	DnsLog       bool     `json:"dnslog"`
	CeyeToken    string   `json:"ceyettoken"`
	CeyeURL      string   `json:"ceyeurl"`
	FullPOC      bool     `json:"fullpoc"`
	FullEXP      bool     `json:"fullexp"`
	SaveResult   bool     `json:"saveresult"`
	OutPutFile   string   `json:"outputfile"`
}

func scanTask(taskID string) {
	defer func() {
		taskWg.Done()
		common.ClearLogChannel(common.LogResults)
	}()
	// 任务完成后减少计数
	fmt.Printf("Start infoscanning target: %s\n", info.Hosts)
	startTime := time.Now()
	config.UseContainer = true
	Modules.HostScan(&info)
	common.GetSugestions()
	result := fmt.Sprintf("Scan complete for target: %s, time used: %s", info.Hosts, time.Since(startTime))
	taskResults.Store(taskID, result)
	// 定时清理ID避免占用内存
	time.AfterFunc(60*time.Second, func() {
		taskResults.Delete(taskID)
		fmt.Printf("Task result for %s deleted after 60s\n", taskID)
	})
}

func scanHandler(w http.ResponseWriter, r *http.Request) {
	// 接收扫描config参数,info赋值
	// 由命令行传参转为url传参
	//config.EnablePocContainer = true
	if err := decodeJSONBody(r); err != nil {
		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
		return
	}
	info.Hosts = r.URL.Query().Get("hosts")
	info.Brute, _ = strconv.ParseBool(r.URL.Query().Get("brute"))
	info.FTPReadFile = r.URL.Query().Get("ftpreadfile")
	info.FTPWriteFile = r.URL.Query().Get("ftpwritefile")
	info.SshKey = r.URL.Query().Get("sshkey")
	info.Domain = r.URL.Query().Get("domain")
	info.SkipRedis, _ = strconv.ParseBool(r.URL.Query().Get("skipredis"))
	info.RedisSshFile = r.URL.Query().Get("redissshfile")
	info.RedisWebshellFile = r.URL.Query().Get("rediswebshellfile")
	info.RedisCronHost = r.URL.Query().Get("rediscronhost")
	info.RemotePath = r.URL.Query().Get("remotepath")
	common.LogWaitTime, _ = strconv.ParseInt(r.URL.Query().Get("logwaittime"), 10, 64)
	common.PrintLog, _ = strconv.ParseBool(r.URL.Query().Get("printlog"))
	common.SaveLogToJSON, _ = strconv.ParseBool(r.URL.Query().Get("savelogtojson"))
	common.SaveLogToHTML, _ = strconv.ParseBool(r.URL.Query().Get("savelogtohtml"))
	// 解析并设定一些默认值

	config.Threads = map[bool]int{true: config.Threads, false: 500}[config.Threads > 0]
	//config.ScanType = map[bool]string{true: config.ScanType, false: "all"}[config.ScanType != ""]
	config.Timeout = map[bool]int64{true: config.Timeout, false: 5}[config.Timeout > 0]
	config.WebTimeout = map[bool]int64{true: config.WebTimeout, false: 5}[config.WebTimeout > 0]
	config.PocNum = map[bool]int{true: config.PocNum, false: 20}[config.PocNum > 0]
	config.ExpNum = map[bool]int{true: config.ExpNum, false: 20}[config.ExpNum > 0]
	common.LogWaitTime = map[bool]int64{true: common.LogWaitTime, false: 60}[common.LogWaitTime > 0]

	common.ParseInit(&info)
	//fmt.Printf("config.Ports: %s", config.Ports)
	taskID := fmt.Sprintf("%d", time.Now().UnixNano()) // 生成任务 ID

	taskWg.Add(1) // 记录任务
	go scanTask(taskID)

	fmt.Fprintf(w, "InfoScan started. Check status with task_id: %s, for example: /infoscanresult?task_id=%s\n", taskID, taskID)
}

func resultHandler(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query().Get("task_id")
	if taskID == "" {
		http.Error(w, "Missing task_id parameter", http.StatusBadRequest)
		return
	}

	if result, found := taskResults.Load(taskID); found {
		fmt.Fprintf(w, "%s", result)
		taskResults.Delete(taskID) // 读取后删除，避免占用内存
	} else {
		fmt.Fprintf(w, "Task %s is still running or not found.", taskID)
	}
}
func decodeJSONBody(r *http.Request) error {
	if r.Header.Get("Content-Type") != "application/json" {
		return fmt.Errorf("invalid content type")
	}
	var scanConfig ScanConfig
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&scanConfig); err != nil {
		return err
	}
	//fmt.Printf("scanconfig.ports: %s", scanConfig.Ports)
	config.HostFile = scanConfig.HostFile
	config.Threads = scanConfig.Threads
	config.ScanType = scanConfig.ScanType
	config.Ports = scanConfig.Ports
	config.PortsFile = scanConfig.PortsFile
	config.Timeout = scanConfig.Timeout
	config.URL = scanConfig.URL
	config.URLFile = scanConfig.URLFile
	config.Username = scanConfig.Username
	config.UsernameFile = scanConfig.UsernameFile
	config.Password = scanConfig.Password
	config.PasswordFile = scanConfig.PasswordFile
	config.HashFile = scanConfig.HashFile
	config.Hash = scanConfig.Hash
	config.AddPorts = scanConfig.AddPorts
	config.AddPassWords = scanConfig.AddPassWords
	config.AddUserNames = scanConfig.AddUserNames
	config.BruteThreads = scanConfig.BruteThreads
	config.Command = scanConfig.Command
	config.Cookie = scanConfig.Cookie
	config.PocNum = scanConfig.PocNum
	config.PocPath = scanConfig.PocPath
	config.PocName = scanConfig.PocName
	config.PocType = scanConfig.PocType
	config.ExpNum = scanConfig.ExpNum
	config.ExpPath = scanConfig.ExpPath
	config.ExpType = scanConfig.ExpType
	config.ExpName = scanConfig.ExpName
	config.WebTimeout = scanConfig.WebTimeout
	config.NoPOC = scanConfig.NoPOC
	config.NoExploit = scanConfig.NoExploit
	config.DnsLog = scanConfig.DnsLog
	config.CeyeToken = scanConfig.CeyeToken
	config.CeyeURL = scanConfig.CeyeURL
	config.FullPOC = scanConfig.FullPOC
	config.FullEXP = scanConfig.FullEXP
	config.SaveResult = scanConfig.SaveResult
	config.OutPutFile = scanConfig.OutPutFile
	return nil
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Infoscan is running.")
}

func main() {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	http.HandleFunc("/infoscan", scanHandler)
	http.HandleFunc("/infoscanresult", resultHandler)
	http.HandleFunc("/ping", pingHandler)
	go func() {
		taskWg.Wait() // 等待所有后台任务完成
		fmt.Println("All InfoScan tasks completed.")
	}()
	server := &http.Server{Addr: ":8080"}
	go func() {
		fmt.Println("Server started on :8080")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()
	// 加一个close channel的监听
	//fmt.Println("Server started on :8080")
	//http.ListenAndServe(":8080", nil)
	_ = <-shutdown
	fmt.Println("InfoScan module exit")
	close(common.LogResults)
	os.Exit(0)
}

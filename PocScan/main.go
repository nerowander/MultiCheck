package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nerowander/MultiCheck/PocScan/Modules"
	"github.com/nerowander/MultiCheck/WebScan/lib"
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

type PocConfig struct {
	PocNum     int    `json:"pocnum"`
	WebTimeout int64  `json:"webtimeout"`
	PocType    string `json:"poctype"`
	PocName    string `json:"pocname"`
	PocPath    string `json:"pocpath"`
	Cookie     string `json:"cookie"`
}

func scanTask(taskID string) {
	defer func() {
		taskWg.Done()
		common.ClearLogChannel(common.LogResults)
	}() // 任务完成后减少计数
	fmt.Printf("Start pocscanning target: %s\n", info.Url)
	startTime := time.Now()
	config.UseContainer = true
	lib.InitHTTP()
	Modules.WebPocScan(&info)
	common.GetSugestions()
	result := fmt.Sprintf("PocScan complete for target: %s, time used: %s", info.Hosts, time.Since(startTime))
	taskResults.Store(taskID, result)
	time.AfterFunc(60*time.Second, func() {
		taskResults.Delete(taskID)
		fmt.Printf("Task result for %s deleted after 60s\n", taskID)
	})
}

func scanHandler(w http.ResponseWriter, r *http.Request) {
	// 接收扫描config参数,info赋值
	if err := decodeJSONBody(r); err != nil {
		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
		return
	}
	info.Hosts = r.URL.Query().Get("hosts")
	info.Url = r.URL.Query().Get("url")
	common.LogWaitTime, _ = strconv.ParseInt(r.URL.Query().Get("logwaittime"), 10, 64)
	common.PrintLog, _ = strconv.ParseBool(r.URL.Query().Get("printlog"))
	common.SaveLogToJSON, _ = strconv.ParseBool(r.URL.Query().Get("savelogtojson"))
	common.SaveLogToHTML, _ = strconv.ParseBool(r.URL.Query().Get("savelogtohtml"))

	config.WebTimeout = map[bool]int64{true: config.WebTimeout, false: 5}[config.WebTimeout > 0]
	config.PocNum = map[bool]int{true: config.PocNum, false: 20}[config.PocNum > 0]
	common.LogWaitTime = map[bool]int64{true: common.LogWaitTime, false: 60}[common.LogWaitTime > 0]

	taskID := fmt.Sprintf("%d", time.Now().UnixNano()) // 生成任务 ID

	taskWg.Add(1) // 记录任务
	go scanTask(taskID)

	fmt.Fprintf(w, "PocScan started. Check status with task_id: %s, for example: /pocscanresult?task_id=%s\n", taskID, taskID)
}

func decodeJSONBody(r *http.Request) error {
	if r.Header.Get("Content-Type") != "application/json" {
		return fmt.Errorf("invalid content type")
	}
	var pocConfig PocConfig
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&pocConfig); err != nil {
		return err
	}
	config.PocNum = pocConfig.PocNum
	config.PocName = pocConfig.PocName
	config.PocPath = pocConfig.PocPath
	config.PocType = pocConfig.PocType
	config.Cookie = pocConfig.Cookie
	config.WebTimeout = pocConfig.WebTimeout

	return nil
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

func pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "PocScan is running.")
}

func main() {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	http.HandleFunc("/pocscan", scanHandler)
	http.HandleFunc("/pocscanresult", resultHandler)
	http.HandleFunc("/ping", pingHandler)
	go func() {
		taskWg.Wait() // 等待所有后台任务完成
		fmt.Println("All PocScan tasks completed.")
	}()
	server := &http.Server{Addr: ":8080"}
	go func() {
		fmt.Println("Server started on :8080")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()
	_ = <-shutdown
	fmt.Println("PocScan module exit")
	close(common.LogResults)
	os.Exit(0)
}

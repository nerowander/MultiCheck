package common

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/nerowander/MultiCheck/config"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// collect results to log and save results to the log files
var Num int64
var End int64
var LogResults = make(chan *string)
var LogSuccessTime int64
var LogErrorTime int64
var LogWaitTime int64 = 60
var LogWG sync.WaitGroup
var PrintLog bool = true
var SaveLogToJSON bool = false
var SaveLogToHTML bool = false
var FileTime int64
var LogFileName string

type JsonLog struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

var htmlLogs []string

func init() {
	log.SetOutput(io.Discard)
	LogSuccessTime = time.Now().Unix()
	FileTime = time.Now().Unix()

	fileExt := filepath.Ext(config.OutPutFile)
	fileNameWithoutExt := strings.TrimSuffix(config.OutPutFile, fileExt)
	LogFileName = fmt.Sprintf("%s_%d%s", fileNameWithoutExt, FileTime, fileExt)

	go logSave()
}
func LogSuccess(res string) {
	LogWG.Add(1)
	LogSuccessTime = time.Now().Unix()
	LogResults <- &res
}
func LogError(errInfo interface{}) {
	if LogWaitTime == 0 {
		fmt.Printf("Finished %v/%v %v \n", End, Num, errInfo)
	} else if (time.Now().Unix()-LogSuccessTime) > LogWaitTime && (time.Now().Unix()-LogErrorTime) > LogWaitTime {
		fmt.Printf("Finished %v/%v %v \n", End, Num, errInfo)
		LogErrorTime = time.Now().Unix()
	}
}

func logSave() {
	for result := range LogResults {
		config.ResLogs = append(config.ResLogs, *result)
		if PrintLog {
			if strings.HasPrefix(*result, "[+] InfoScan") {
				color.Cyan(*result)
			} else if strings.HasPrefix(*result, "[+] PocScan") {
				color.Blue(*result)
			} else if strings.HasPrefix(*result, "[+] ExpScan") {
				color.Magenta(*result)
			} else if strings.HasPrefix(*result, "[+]") {
				color.Green(*result)
			} else {
				fmt.Println(*result)
			}
		}
		if config.SaveResult {
			WriteLogToFile(*result, LogFileName)
		}
		LogWG.Done()
	}

}

// todo: 后面可以弄一个只打开一次文件的优化
func WriteLogToFile(result string, filename string) {
	// json or html
	// 最好加一个时间戳的文件
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Open file %s error, %v\n", filename, err)
		return
	}
	ext := filepath.Ext(filename)

	if SaveLogToJSON && strings.EqualFold(ext, ".json") {
		var scanType string
		var text string
		// for example
		// [+] ftp connect success: %v:%v:%v %v
		if strings.HasPrefix(result, "[+]") || strings.HasPrefix(result, "[*]") || strings.HasPrefix(result, "[-]") {
			//找到第二个空格的位置
			index := strings.Index(result[4:], " ")
			if index == -1 {
				scanType = "msg"
				text = result[4:]
			} else {
				scanType = result[4 : 4+index]
				text = result[4+index+1:]
			}
		} else {
			scanType = "msg"
			text = result
		}
		jsonText := JsonLog{
			Type: scanType,
			Text: text,
		}
		var jsonData []byte
		// 序列化，转换成json对象
		jsonData, err = json.Marshal(jsonText)
		if err != nil {
			fmt.Println(err)
			jsonText = JsonLog{
				Type: "msg",
				Text: result,
			}
			jsonData, err = json.Marshal(jsonText)
			if err != nil {
				fmt.Println(err)
				jsonData = []byte(result)
			}
		}
		jsonData = append(jsonData, []byte(",\n")...)
		_, err = file.Write(jsonData)
	} else if SaveLogToHTML && strings.EqualFold(ext, ".html") {
		// 如果文件为空，添加HTML头部
		if fileInfo, _ := file.Stat(); fileInfo.Size() == 0 {
			_, err = file.WriteString(`
			<html>
			<head>
				<title>Service Logs</title>
				<style>
					<style>
        body {
            font-family: 'Arial', sans-serif;
            background-color: #f4f7fc;
            color: #333;
            margin: 0;
            padding: 20px;
        }

        h2 {
            text-align: center;
            color: #4CAF50;
            font-size: 24px;
            margin-top: 30px;
        }

        table {
            width: 100%;
            margin-top: 20px;
            border-collapse: collapse;
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
            background-color: #fff;
            table-layout: fixed; /* 固定表格布局 */
        }

        th, td {
            padding: 12px 20px;
            text-align: left;
            border-bottom: 1px solid #ddd;
            word-wrap: break-word; /* 强制内容换行 */
        }

        th {
            background-color: #4CAF50;
            color: white;
        }

        tr:hover {
            background-color: #f1f1f1;
        }

        .log-entry {
            font-size: 14px;
        }

        .log-time {
            color: #888;
            max-width: 200px;
            overflow: hidden;
            text-overflow: ellipsis; /* 超过宽度的文本省略显示 */
            white-space: nowrap;
        }

        .log-message {
            color: #555;
            max-width: 600px; /* 设置最大宽度 */
        }

        .container {
            max-width: 1000px;
            margin: 0 auto;
            background-color: #fff;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
        }

        .footer {
            text-align: center;
            font-size: 12px;
            color: #888;
            margin-top: 40px;
        }

        .header {
            text-align: center;
            margin-bottom: 20px;
        }

        .header h1 {
            font-size: 36px;
            color: #4CAF50;
        }
    </style>
			</head>
			<body>
				<div class="container">
					<div class="header">
						<h1>Service Logs</h1>
						<p>Latest logs from your services (ftp, FCGI, Mongodb)</p>
					</div>
		`)
			if err != nil {
				fmt.Println("Error writing to file:", err)
				//return
			}
		}
		// 定义不同服务的日志
		logsByService := map[string][]Log{
			"ftp":       {},
			"FCGI":      {},
			"Mongodb":   {},
			"Redis":     {},
			"Memcached": {},
			"mssql":     {},
			"mysql":     {},
			"orcale":    {},
			"postgres":  {},
			"RDP":       {},
			"SMB":       {},
			"SSH":       {},
			"InfoScan":  {},
			"PocScan":   {},
			"ExpScan":   {},
			"Err":       {},
		}
		htmlLogs = append(htmlLogs, result)
		// 将日志分类
		for _, log := range htmlLogs {
			service, message := parseLog(log)
			logsByService[service] = append(logsByService[service], Log{Service: service, Message: message})
		}

		// 将每个服务的日志写入HTML
		for service, serviceLogs := range logsByService {
			_, err = file.WriteString(fmt.Sprintf("<h2>%s Logs</h2><table><tr><th>Time</th><th>Message</th></tr>", service))
			if err != nil {
				fmt.Println("Error writing to file:", err)
				//return
			}
			for _, log := range serviceLogs {
				_, err = file.WriteString(fmt.Sprintf("<tr><td class='log-time'>%s</td><td class='log-message'>%s</td></tr>", time.Now().Format("2006-01-02 15:04:05"), log.Message))
				if err != nil {
					fmt.Println("Error writing to file:", err)
					//return
				}
			}
			_, err = file.WriteString("</table><br>")
			if err != nil {
				fmt.Println("Error writing to file:", err)
				//return
			}
		}

		// 结束HTML文件
		_, err = file.WriteString(`
		<div class="footer">
			<p>Logs updated at <strong>` + time.Now().Format("2006-01-02 15:04:05") + `</strong></p>
		</div>
		</div>
	</body>
	</html>
	`)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			//return
		}
	} else {
		// txt情况，非json或html
		_, err = file.Write([]byte(result + "\n"))
	}
	file.Close()
	if err != nil {
		fmt.Printf("Write json %s error, %v\n", filename, err)
	}
}

func CheckErrMessages(err error) bool {
	if err == nil {
		return false
	}
	errs := []string{
		"closed by the remote host", "too many connections",
		"i/o timeout", "EOF", "A connection attempt failed",
		"established connection failed", "connection attempt failed",
		"Unable to read", "is not allowed to connect to this",
		"no pg_hba.conf entry",
		"No connection could be made",
		"invalid packet size",
		"bad connection",
	}
	for _, key := range errs {
		if strings.Contains(strings.ToLower(err.Error()), strings.ToLower(key)) {
			return true
		}
	}
	return false
}

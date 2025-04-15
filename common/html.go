package common

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type Log struct {
	Service string
	Message string
}

func startTest() {
	// 模拟日志输入
	logs := []string{
		"[+] ftp connect success: 192.168.1.1:21 2023-03-04 12:00:00",
		"[+] FCGI 192.168.1.1:9000 \nService start",
		"[+] ftp write file /tmp/file.txt failed",
		"[+] Mongodb 127.0.0.1 unauthorized",
		"[+] redis tests",
		"[+] ftp connect success: 192.168.1.1:21 2023-03-04 12:00:00",
		"k",
	}

	// 打开或创建HTML文件
	file, err := os.OpenFile("logs.html", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

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
			return
		}
	}

	// 定义不同服务的日志
	logsByService := map[string][]Log{
		"ftp":     {},
		"FCGI":    {},
		"Mongodb": {},
		"redis":   {},
		"k":       {},
	}

	// 将日志分类
	for _, log := range logs {
		service, message := parseLog(log)
		logsByService[service] = append(logsByService[service], Log{Service: service, Message: message})
	}

	// 将每个服务的日志写入HTML
	for service, serviceLogs := range logsByService {
		_, err = file.WriteString(fmt.Sprintf("<h2>%s Logs</h2><table><tr><th>Time</th><th>Message</th></tr>", service))
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
		for _, log := range serviceLogs {
			_, err = file.WriteString(fmt.Sprintf("<tr><td class='log-time'>%s</td><td class='log-message'>%s</td></tr>", time.Now().Format("2006-01-02 15:04:05"), log.Message))
			if err != nil {
				fmt.Println("Error writing to file:", err)
				return
			}
		}
		_, err = file.WriteString("</table><br>")
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
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
		return
	}

	fmt.Println("Logs have been written to logs.html")
}

// 解析日志信息，返回服务名称和消息
func parseLog(log string) (string, string) {
	if strings.Contains(log, "ftp") {
		return "ftp", log
	} else if strings.Contains(log, "FCGI") {
		return "FCGI", log
	} else if strings.Contains(log, "Mongodb") {
		return "Mongodb", log
	} else if strings.Contains(log, "Redis") {
		return "Redis", log
	} else if strings.Contains(log, "Memcached") {
		return "Memcached", log
	} else if strings.Contains(log, "mssql") {
		return "mssql", log
	} else if strings.Contains(log, "mysql") {
		return "mysql", log
	} else if strings.Contains(log, "orcale") {
		return "orcale", log
	} else if strings.Contains(log, "postgres") {
		return "postgres", log
	} else if strings.Contains(log, "RDP") {
		return "RDP", log
	} else if strings.Contains(log, "SMB") {
		return "SMB", log
	} else if strings.Contains(log, "SSH") {
		return "SSH", log
	} else if strings.Contains(log, "PocScan") {
		return "PocScan", log
	} else if strings.Contains(log, "ExpScan") {
		return "ExpScan", log
	} else if strings.Contains(log, "[*]") {
		return "InfoScan", log
	} else if strings.Contains(log, "[-]") {
		return "Err", log
	}
	return "unknown", log
}

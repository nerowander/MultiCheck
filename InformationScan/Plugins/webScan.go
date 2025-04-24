package Plugins

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"encoding/json"
	"fmt"
	Modules2 "github.com/nerowander/MultiCheck/ExploitScan/Modules"
	"github.com/nerowander/MultiCheck/PocScan/Modules"
	"github.com/nerowander/MultiCheck/WebScan"
	"github.com/nerowander/MultiCheck/WebScan/lib"
	"github.com/nerowander/MultiCheck/common"
	"github.com/nerowander/MultiCheck/config"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

func webScan(info *config.InfoScan) error {

	// 判断是否是容器模式，如果是则传参数到另一个函数，判断请求poc容器/exp容器，或者2个都有
	if config.EnableVulContainer == true {
		if config.ScanType == "pocscan" || config.ScanType == "all" {
			module := "pocscan"
			url := fmt.Sprintf("http://%s-service/pocscan?hosts=%s&url=%s&logwaittime=%s&printlog=%s&savelogtojson=%s&savelogtohtml=%s",
				module, info.Hosts, info.Url, strconv.FormatInt(common.LogWaitTime, 10),
				strconv.FormatBool(common.PrintLog), strconv.FormatBool(common.SaveLogToJSON),
				strconv.FormatBool(common.SaveLogToHTML))

			fmt.Println(url)
			jsonPayload := map[string]interface{}{
				"pocnum":     config.PocNum,
				"webtimeout": config.WebTimeout,
				"poctype":    config.PocType,
				"pocname":    config.PocName,
				"pocpath":    config.PocPath,
				"cookie":     config.Cookie,
			}
			var body []byte
			var err error
			body, err = json.Marshal(jsonPayload)
			if err != nil {
				fmt.Printf("Error marshaling JSON: %v\n", err)
				return nil
			}
			var resp *http.Response
			resp, err = http.Post(url, "application/json", bytes.NewBuffer(body))
			if err != nil {
				fmt.Printf("Error calling module %s: %v\n", module, err)
				return nil
			}
			defer resp.Body.Close()
			resBody, _ := ioutil.ReadAll(resp.Body)
			fmt.Println(string(resBody))
			//return nil
		} else if config.ScanType == "exploit" || config.ScanType == "all" {
			module := "exploit"
			url := fmt.Sprintf("http://%s-service/expscan?hosts=%s&url=%s&logwaittime=%s&printlog=%s&savelogtojson=%s&savelogtohtml=%s",
				module, info.Hosts, info.Url, strconv.FormatInt(common.LogWaitTime, 10),
				strconv.FormatBool(common.PrintLog), strconv.FormatBool(common.SaveLogToJSON),
				strconv.FormatBool(common.SaveLogToHTML))
			fmt.Println(url)
			jsonPayload := map[string]interface{}{
				"expnum":     config.ExpNum,
				"webtimeout": config.WebTimeout,
				"exptype":    config.ExpType,
				"expname":    config.ExpName,
				"exppath":    config.ExpPath,
				"cookie":     config.Cookie,
			}
			var body []byte
			var err error
			body, err = json.Marshal(jsonPayload)
			if err != nil {
				fmt.Printf("Error marshaling JSON: %v\n", err)
				return nil
			}
			var resp *http.Response
			resp, err = http.Post(url, "application/json", bytes.NewBuffer(body))
			if err != nil {
				fmt.Printf("Error calling module %s: %v\n", module, err)
				return nil
			}
			defer resp.Body.Close()
			resBody, _ := ioutil.ReadAll(resp.Body)
			fmt.Println(string(resBody))
			//return nil
		}
		//} else if config.ScanType == "all" {
		//	// 后面测试时再补充
		//	return nil
		//}
		return nil
	} else {
		err, CheckData := GOWebTitle(info)
		info.WebInfo = WebScan.InfoCheck(info.Url, &CheckData)
		//不扫描打印机
		for _, v := range info.WebInfo {
			if v == "打印机" {
				return nil
			}
		}
		if config.ScanType == "pocscan" {
			// poc扫描 web+iot
			Modules.WebPocScan(info)
			return nil
			// 后续看这个return nil要不要去掉
		} else if config.ScanType == "exploit" {
			// exploit利用 web+iot
			Modules2.WebExploit(info)
			return nil
		} else if config.ScanType == "all" {
			Modules.WebPocScan(info)
			Modules2.WebExploit(info)
			return nil
		}
		if config.NoPOC == false && err == nil {
			Modules.WebPocScan(info)
		} else if config.NoExploit == false && err == nil {
			Modules2.WebExploit(info)
		} else {
			errlog := fmt.Sprintf("[-] webtitle %v %v", info.Url, err)
			common.LogError(errlog)
		}
		return err
	}
}
func GOWebTitle(info *config.InfoScan) (err error, CheckData []WebScan.CheckDatas) {
	if info.Url == "" {
		switch info.Ports {
		case "80":
			info.Url = fmt.Sprintf("http://%s", info.Hosts)
		case "443":
			info.Url = fmt.Sprintf("https://%s", info.Hosts)
		default:
			host := fmt.Sprintf("%s:%s", info.Hosts, info.Ports)
			protocol := GetProtocol(host, config.Timeout)
			info.Url = fmt.Sprintf("%s://%s:%s", protocol, info.Hosts, info.Ports)
		}
	} else {
		if !strings.Contains(info.Url, "://") {
			host := strings.Split(info.Url, "/")[0]
			protocol := GetProtocol(host, config.Timeout)
			info.Url = fmt.Sprintf("%s://%s", protocol, info.Url)
		}
	}
	var result string
	err, result, CheckData = getURL(info, 1, CheckData)
	if err != nil && !strings.Contains(err.Error(), "EOF") {
		return
	}

	//redirect
	if strings.Contains(result, "://") {
		info.Url = result
		err, result, CheckData = getURL(info, 3, CheckData)
		if err != nil {
			return
		}
	}

	if result == "https" && !strings.HasPrefix(info.Url, "https://") {
		info.Url = strings.Replace(info.Url, "http://", "https://", 1)
		err, result, CheckData = getURL(info, 1, CheckData)
		//有跳转
		if strings.Contains(result, "://") {
			info.Url = result
			err, _, CheckData = getURL(info, 3, CheckData)
			if err != nil {
				return
			}
		}
	}
	//是否访问图标
	//err, _, CheckData = geturl(info, 2, CheckData)
	//if err != nil {
	//	return
	//}
	return
}

func getURL(info *config.InfoScan, flag int, CheckData []WebScan.CheckDatas) (error, string, []WebScan.CheckDatas) {
	//flag 1 first try
	//flag 2 /favicon.ico
	//flag 3 302
	//flag 4 400 -> https

	Url := info.Url
	if flag == 2 {
		URL, err := url.Parse(Url)
		if err == nil {
			Url = fmt.Sprintf("%s://%s/favicon.ico", URL.Scheme, URL.Host)
		} else {
			Url += "/favicon.ico"
		}
	}
	req, err := http.NewRequest("GET", Url, nil)
	if err != nil {
		return err, "", CheckData
	}
	req.Header.Set("User-agent", config.UserAgent)
	req.Header.Set("Accept", config.Accept)
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	if config.Cookie != "" {
		req.Header.Set("Cookie", config.Cookie)
	}
	req.Header.Set("Connection", "close")
	var client *http.Client
	if flag == 1 {
		// 不处理重定向
		client = lib.ClientNoRedirect
	} else {
		// 处理重定向
		client = lib.Client
	}
	var resp *http.Response
	resp, err = client.Do(req)

	// get https
	if err != nil {
		return err, "https", CheckData
	}

	defer resp.Body.Close()
	var title string
	var body []byte
	body, err = getRespBody(resp)
	// get https
	if err != nil {
		return err, "https", CheckData
	}
	CheckData = append(CheckData, WebScan.CheckDatas{body, fmt.Sprintf("%s", resp.Header)})
	var redirectUrl string
	if flag != 2 {
		if !utf8.Valid(body) {
			body, _ = simplifiedchinese.GBK.NewDecoder().Bytes(body)
		}
		// webtitle
		title = getTitle(body)
		length := resp.Header.Get("Content-Length")
		if length == "" {
			length = fmt.Sprintf("%v", len(body))
		}
		redirURL, err1 := resp.Location()
		if err1 == nil {
			redirectUrl = redirURL.String()
		}
		result := fmt.Sprintf("[*] WebTitle %-25v code:%-3v len:%-6v title:%v", resp.Request.URL, resp.StatusCode, length, title)
		if redirectUrl != "" {
			result += fmt.Sprintf(" 跳转url: %s", redirectUrl)
		}
		common.LogSuccess(result)
	}
	if redirectUrl != "" {
		return nil, redirectUrl, CheckData
	}
	if resp.StatusCode == 400 && !strings.HasPrefix(info.Url, "https") {
		return nil, "https", CheckData
	}
	return nil, "", CheckData
}

func getRespBody(oResp *http.Response) ([]byte, error) {
	var body []byte
	if oResp.Header.Get("Content-Encoding") == "gzip" {
		gr, err := gzip.NewReader(oResp.Body)
		if err != nil {
			return nil, err
		}
		defer gr.Close()
		for {
			buf := make([]byte, 1024)
			n, err := gr.Read(buf)
			if err != nil && err != io.EOF {
				return nil, err
			}
			if n == 0 {
				break
			}
			body = append(body, buf...)
		}
	} else {
		raw, err := io.ReadAll(oResp.Body)
		if err != nil {
			return nil, err
		}
		body = raw
	}
	return body, nil
}

func getTitle(body []byte) (title string) {
	re := regexp.MustCompile("(?ims)<title.*?>(.*?)</title>")
	find := re.FindSubmatch(body)
	if len(find) > 1 {
		title = string(find[1])
		title = strings.TrimSpace(title)
		title = strings.Replace(title, "\n", "", -1)
		title = strings.Replace(title, "\r", "", -1)
		title = strings.Replace(title, "&nbsp;", " ", -1)
		if len(title) > 100 {
			title = title[:100]
		}
		if title == "" {
			title = "\"\"" //空格
		}
	} else {
		title = "None" //没有title
	}
	return
}

func GetProtocol(host string, Timeout int64) (protocol string) {
	protocol = "http"
	//如果端口是80或443,跳过Protocol判断
	if strings.HasSuffix(host, ":80") || !strings.Contains(host, ":") {
		return
	} else if strings.HasSuffix(host, ":443") {
		protocol = "https"
		return
	}

	socksconn, err := common.TestTCPWithTimeout("tcp", host, time.Duration(Timeout)*time.Second)
	if err != nil {
		return
	}
	conn := tls.Client(socksconn, &tls.Config{MinVersion: tls.VersionTLS10, InsecureSkipVerify: true})
	defer func() {
		if conn != nil {
			defer func() {
				if err := recover(); err != nil {
					common.LogError(err)
				}
			}()
			conn.Close()
		}
	}()
	conn.SetDeadline(time.Now().Add(time.Duration(Timeout) * time.Second))
	err = conn.Handshake()
	if err == nil || strings.Contains(err.Error(), "handshake failure") {
		protocol = "https"
	}
	return protocol
}

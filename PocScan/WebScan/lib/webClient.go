package lib

import (
	"PocScan/config"
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

var (
	Client           *http.Client
	ClientNoRedirect *http.Client
	dialTimout       = 5 * time.Second
	keepAlive        = 5 * time.Second
)

func InitHTTP() {
	if config.PocNum == 0 {
		config.PocNum = 20
	}
	if config.WebTimeout == 0 {
		config.WebTimeout = 5
	}
	err := InitHttpClient(config.PocNum, time.Duration(config.WebTimeout)*time.Second)
	// todo: 考虑加proxy
	if err != nil {
		panic(err)
	}
}
func InitHttpClient(ThreadsNum int, Timeout time.Duration) error {
	// todo: 考虑加proxy
	type DialContext = func(ctx context.Context, network, addr string) (net.Conn, error)
	dialer := &net.Dialer{
		Timeout:   dialTimout,
		KeepAlive: keepAlive,
	}

	tr := &http.Transport{
		DialContext:         dialer.DialContext,
		MaxConnsPerHost:     5,
		MaxIdleConns:        0,
		MaxIdleConnsPerHost: ThreadsNum * 2,
		IdleConnTimeout:     keepAlive,
		TLSClientConfig:     &tls.Config{MinVersion: tls.VersionTLS10, InsecureSkipVerify: true},
		TLSHandshakeTimeout: 5 * time.Second,
		DisableKeepAlives:   false,
	}

	Client = &http.Client{
		Transport: tr,
		Timeout:   Timeout,
	}

	// 若重定向则返回上一次的响应
	ClientNoRedirect = &http.Client{
		Transport:     tr,
		Timeout:       Timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse },
	}
	return nil
}

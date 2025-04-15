package Modules

import (
	"fmt"
	"github.com/nerowander/MultiCheck/common"
	"github.com/nerowander/MultiCheck/config"
	"sync"
	"time"
)

type Addr struct {
	host string
	port int
}

func CheckPortLive(hosts []string, ports string, timeout int64) []string {
	var resAddress []string
	checkPorts := common.ParsePort(ports)
	if len(checkPorts) == 0 {
		fmt.Printf("[-] parse port %s error, the port format: 80 | 80,443 | 80-443\n", ports)
		return resAddress
	}
	// todo：添加noports：不扫描的端口
	workers := config.Threads
	addrs := make(chan Addr, 100)
	results := make(chan string, 100)
	var wg sync.WaitGroup

	// receive results
	go func() {
		for found := range results {
			resAddress = append(resAddress, found)
			wg.Done()
		}
	}()
	// scanports
	for i := 0; i < workers; i++ {
		go func() {
			for addr := range addrs {
				connectPort(addr, results, timeout, &wg)
				wg.Done()
			}
		}()
	}
	// add ports
	for _, port := range checkPorts {
		for _, host := range hosts {
			wg.Add(1)
			addrs <- Addr{
				host,
				port,
			}
		}
	}
	wg.Wait()
	close(addrs)
	close(results)
	return resAddress
}

func connectPort(addr Addr, results chan string, timeout int64, wg *sync.WaitGroup) {
	host, port := addr.host, addr.port
	conn, err := common.TestTCPWithTimeout("tcp4", fmt.Sprintf("%s:%d", host, port), time.Duration(timeout)*time.Second)
	if err == nil {
		defer conn.Close()
		address := fmt.Sprintf("%s:%d", host, port)
		result := fmt.Sprintf("%s is open", address)
		common.LogSuccess(result)
		wg.Add(1)
		results <- address
	}
}

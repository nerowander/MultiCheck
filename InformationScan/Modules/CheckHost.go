package Modules

import (
	"bytes"
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"sync"
)

var (
	liveWg     sync.WaitGroup
	ExistHosts = make(map[string]struct{})
	AliveHosts []string
)

func CheckHostLive(hosts []string) (resHosts []string) {
	scanHosts := make(chan string, len(hosts))
	go func() {
		for ip := range scanHosts {
			if _, ok := ExistHosts[ip]; !ok && IsContain(hosts, ip) {
				ExistHosts[ip] = struct{}{}
				fmt.Printf("(ping) Target %s is alive\n", ip)
				AliveHosts = append(AliveHosts, ip)
			}
			liveWg.Done()
		}
	}()

	RunPing(hosts, scanHosts)
	liveWg.Wait()
	close(scanHosts)

	//if len(hosts) > 1000 {
	//
	//}
	//
	//if len(hosts) > 256 {
	//
	//}
	// todo 统计ip数量
	// todo：考虑添加icmp扫描模式？
	// todo: collect scan result to file and log
	// todo: common.LogSuccess
	return AliveHosts
}

func RunPing(hosts []string, scanHosts chan string) {
	var wg sync.WaitGroup
	limiter := make(chan struct{}, 50)
	for _, host := range hosts {
		wg.Add(1)
		limiter <- struct{}{}
		go func(host string) {
			if CheckWithCommandPing(host) {
				liveWg.Add(1)
				scanHosts <- host
			}
			<-limiter
			wg.Done()
		}(host)
	}
	wg.Wait()
}

func CheckWithCommandPing(host string) bool {
	var cmd *exec.Cmd
	checkIP := net.ParseIP(host)
	checkHost, _ := net.LookupHost(host) // 例如www.baidu.com不带IP的情况
	if checkIP == nil && len(checkHost) == 0 {
		return false
	}
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("ping", "-n", "1", "-w", "1", host)
		//cmd = exec.Command("ping", "/c", "ping -n 1 -w 1 "+host+" && echo true || echo false")
	case "darwin":
		cmd = exec.Command("ping", "-c", "1", "-W", "1", host)
		//cmd = exec.Command("/bin/bash", "-c", "ping -c 1 -W 1 "+host+" && echo true || echo false")
	default:
		cmd = exec.Command("ping", "-c", "1", "-w", "1", host)
	}
	//cmd = exec.Command("/bin/bash", "-c", "ping -c 1 -w 1 "+host+" && echo true || echo false")
	outinfo := bytes.Buffer{}
	cmd.Stdout = &outinfo
	err := cmd.Run()
	if err != nil {
		return false
	} else {
		//fmt.Println(outinfo.String())
		//if strings.Count(outinfo.String(), host) >= 2 {
		//	return true
		//} else {
		//	return false
		//}
		return true
	}
}
func IsContain(hosts []string, ip string) bool {
	for _, host := range hosts {
		if host == ip {
			return true
		}
	}
	return false
}

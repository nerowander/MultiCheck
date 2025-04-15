package Plugins

import (
	"fmt"
	"github.com/nerowander/MultiCheck/common"
	"github.com/nerowander/MultiCheck/config"
	"net"
	"strings"
	"time"
)

func memcachedScan(info *config.InfoScan) (err error) {
	realhost := fmt.Sprintf("%s:%v", info.Hosts, info.Ports)
	var client net.Conn
	client, err = common.TestTCPWithTimeout("tcp", realhost, time.Duration(config.Timeout)*time.Second)
	defer func() {
		if client != nil {
			client.Close()
		}
	}()
	if err == nil {
		err = client.SetDeadline(time.Now().Add(time.Duration(config.Timeout) * time.Second))
		if err == nil {
			_, err = client.Write([]byte("stats\n"))
			if err == nil {
				rev := make([]byte, 1024)
				var n int
				n, err = client.Read(rev)
				if err == nil {
					if strings.Contains(string(rev[:n]), "STAT") {
						result := fmt.Sprintf("[+] Memcached %s unauthorized", realhost)
						common.LogSuccess(result)
					}
				} else {
					errlog := fmt.Sprintf("[-] Memcached unauthorized failed %v:%v %v", info.Hosts, info.Ports, err)
					common.LogError(errlog)
				}
			}
		}
	}
	return err
}

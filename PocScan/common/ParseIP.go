package common

import (
	"PocScan/config"
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func ParseIP(info *config.InfoScan) (hosts []string, err error) {
	if config.HostFile == "" {
		// ip:port
		if strings.Contains(info.Hosts, ":") {
			if strings.Contains(info.Hosts, ",") {
				// ip1,ip2,ip3:port
				hostList := strings.Split(info.Hosts, ",")
				for _, target := range hostList {
					// ip3:port
					if strings.Contains(target, ":") {
						hostsPart, portsPart, _ := net.SplitHostPort(target)
						//info.Ports = portsPart
						config.Ports = portsPart
						resHosts := ParseIPs(hostsPart)
						hosts = append(hosts, resHosts...)

					} else {
						// ip1 ip2
						resHosts := ParseIPs(target)
						hosts = append(hosts, resHosts...)
					}
				}
			} else {
				// ip1:port
				hostPart, portsPart, _ := net.SplitHostPort(info.Hosts)
				//info.Ports = portsPart
				config.Ports = portsPart
				hosts = ParseIPs(hostPart)
			}
			//hostsPart := strings.Split(info.Hosts, ":")
			//
			//hosts = ParseIPs(hostsPart[0])
			//info.Ports = hostsPart[1]
		} else {
			// ip
			hosts = ParseIPs(info.Hosts)
		}
	} else {
		var filehost []string
		filehost, _ = readIPFile(info)
		hosts = append(hosts, filehost...)
	}
	hosts = RemoveDuplicateHosts(hosts)
	if len(hosts) == 0 && len(info.HostPort) == 0 && config.HostFile == "" && info.Hosts == "" {
		return nil, errors.New("Ip format error, the format:" +
			"192.168.0.1\n" +
			"1.2.3.4/24\n" +
			"192.168.0.1\n" +
			"192.168.0.1,192.168.0.2\n" +
			"1.2.3.4-255\n" +
			"1.2.3.4-1.2.3.255\n")
	}
	return hosts, nil
}

func ParseIPs(ip string) (hosts []string) {
	//if strings.Contains(ip, ",") {
	//	IPList := strings.Split(ip, ",")
	//	var ips []string
	//
	//	for _, ip := range IPList {
	//		ips = parseDifferentIP(ip)
	//		hosts = append(hosts, ips...)
	//	}
	//} else {
	//	hosts = parseDifferentIP(ip)
	//}
	hosts = parseDifferentIP(ip)
	return hosts

}

func parseDifferentIP(ip string) (ips []string) {
	reg := regexp.MustCompile(`[a-zA-Z]+`)
	switch {
	case ip == "192":
		return parseIP2("192.168.0.0/16")
	case ip == "172":
		return parseIP2("172.16.0.0/12")
	case ip == "10":
		return parseIP8("10.0.0.0/8")
	case strings.HasSuffix(ip, "/8"):
		return parseIP8(ip)
	case strings.Contains(ip, "/"):
		return parseIP2(ip)
	case strings.Contains(ip, "-"):
		return parseIPRange(ip)
	case reg.MatchString(ip):
		return []string{ip}
	default:
		resIP := net.ParseIP(ip)
		if resIP == nil {
			return nil
		} else {
			return []string{ip}
		}
	}

}

func readIPFile(info *config.InfoScan) ([]string, error) {
	file, err := os.Open(config.HostFile)
	if err != nil {
		return nil, errors.New("Open IP file error")
	}
	defer file.Close()
	var res []string
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			// ip:port
			text := strings.Split(line, ":")
			if len(text) == 2 {
				port := text[1]
				var num int
				num, err = strconv.Atoi(port)
				if err != nil || (num < 1 || num > 65535) {
					continue
				}
				hosts := ParseIPs(text[0])
				for _, host := range hosts {
					info.HostPort = append(info.HostPort, fmt.Sprintf("%s:%s", host, port))
				}
			} else {
				// ip
				hosts := ParseIPs(line)
				res = append(res, hosts...)
			}
		}
	}
	return res, nil

}

func RemoveDuplicateHosts(hosts []string) []string {
	seen := make(map[string]bool)
	var uniqueIPs []string
	for _, ipStr := range hosts {
		ip := net.ParseIP(ipStr)
		if ip != nil && !seen[ipStr] {
			seen[ipStr] = true
			uniqueIPs = append(uniqueIPs, ipStr)
		}
	}
	return uniqueIPs
}

func parseIP8(ip string) (ips []string) {
	resIP := ip[:len(ip)-2]
	checkIP := net.ParseIP(resIP)
	if checkIP == nil {
		return nil
	}
	IPRange := strings.Split(ip, ".")[0]
	var resAllIP []string
	for i := 0; i <= 255; i++ {
		for j := 0; j <= 255; j++ {
			resAllIP = append(resAllIP, fmt.Sprintf("%s.%d.%d.%d", IPRange, i, j, 1))
			resAllIP = append(resAllIP, fmt.Sprintf("%s.%d.%d.%d", IPRange, i, j, 2))
			resAllIP = append(resAllIP, fmt.Sprintf("%s.%d.%d.%d", IPRange, i, j, 4))
			resAllIP = append(resAllIP, fmt.Sprintf("%s.%d.%d.%d", IPRange, i, j, 5))
			resAllIP = append(resAllIP, fmt.Sprintf("%s.%d.%d.%d", IPRange, i, j, randInt(6, 55)))
			resAllIP = append(resAllIP, fmt.Sprintf("%s.%d.%d.%d", IPRange, i, j, randInt(56, 100)))
			resAllIP = append(resAllIP, fmt.Sprintf("%s.%d.%d.%d", IPRange, i, j, randInt(101, 150)))
			resAllIP = append(resAllIP, fmt.Sprintf("%s.%d.%d.%d", IPRange, i, j, randInt(151, 200)))
			resAllIP = append(resAllIP, fmt.Sprintf("%s.%d.%d.%d", IPRange, i, j, randInt(201, 253)))
			resAllIP = append(resAllIP, fmt.Sprintf("%s.%d.%d.%d", IPRange, i, j, 254))
		}
	}
	return resAllIP

}

func parseIP2(ip string) (ips []string) {
	_, resIP, err := net.ParseCIDR(ip)
	if err != nil {
		return nil
	}
	ips = parseIPRange(IPRange(resIP))
	return ips
}

func parseIPRange(ipRange string) (ips []string) {
	// 解析 IP 范围
	parts := strings.Split(ipRange, "-")
	if len(parts) != 2 {
		return nil
	}

	startIP := net.ParseIP(parts[0])
	var endIP net.IP
	if startIP == nil {
		return nil
	}

	// 处理 "192.168.1.1-255" 这种格式
	if len(parts[1]) <= 3 { // 假设最后部分是 1-255
		lastOctet, err := strconv.Atoi(parts[1])
		if err != nil || lastOctet < 1 || lastOctet > 255 {
			return nil
		}
		endIP = net.ParseIP(fmt.Sprintf("%s.%d", parts[0][:strings.LastIndex(parts[0], ".")], lastOctet))
		if endIP == nil {
			return nil
		}
	} else {
		// 192.168.1.1-192.168.1.255
		endIP = net.ParseIP(parts[1])
		if endIP == nil {
			return nil
		}
	}

	startInt := ipToInt(startIP)
	endInt := ipToInt(endIP)

	for i := startInt; i <= endInt; i++ {
		ip := intToIP(uint32(i))
		ips = append(ips, ip.String())
	}

	return ips
}

func ipToInt(ip net.IP) uint32 {
	// 确保 IP 地址是 IPv4
	if len(ip) > net.IPv4len {
		ip = ip[12:]
	}
	return uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3])
}

func intToIP(n uint32) net.IP {
	return net.IPv4(byte(n>>24), byte(n>>16), byte(n>>8), byte(n))
}
func randInt(min int, max int) (res int) {
	rand.Seed(time.Now().UnixNano())
	// 生成随机数
	return rand.Intn(max-min+1) + min
}

func IPRange(c *net.IPNet) (res string) {
	start := c.IP.String()
	mask := c.Mask
	broadcast := make(net.IP, len(c.IP))
	copy(broadcast, c.IP)
	for i := 0; i < len(mask); i++ {
		ipIndex := len(broadcast) - i - 1
		broadcast[ipIndex] = c.IP[ipIndex] | ^mask[len(mask)-i-1]
	}
	end := broadcast.String()
	return fmt.Sprintf("%s-%s", start, end)
}

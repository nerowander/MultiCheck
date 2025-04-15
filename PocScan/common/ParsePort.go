package common

import (
	"PocScan/config"
	"strconv"
	"strings"
)

func ParsePort(ports string) (resPorts []int) {
	if ports == "" {
		return
	}
	parts := strings.Split(ports, ",")
	for _, port := range parts {
		port = strings.TrimSpace(port)
		if port == "" {
			continue
		}
		if config.PortMap[port] != "" {
			port = config.PortMap[port]
			resPorts = append(resPorts, ParsePort(port)...)
			continue
		}
		upper := port
		if strings.Contains(port, "-") {
			ranges := strings.Split(port, "-")
			if len(ranges) < 2 {
				continue
			}
			startPart, _ := strconv.Atoi(ranges[0])
			endPart, _ := strconv.Atoi(ranges[1])
			if startPart < endPart {
				port = ranges[0]
				upper = ranges[1]
			} else {
				port = ranges[1]
				upper = ranges[0]
			}
		}
		startPort, _ := strconv.Atoi(port)
		endPort, _ := strconv.Atoi(upper)
		for i := startPort; i <= endPort; i++ {
			if i > 65535 || i < 1 {
				continue
			}
			resPorts = append(resPorts, i)

		}
	}
	resPorts = removeDuplicatePorts(resPorts)
	return resPorts
}

func removeDuplicatePorts(input []int) []int {
	result := []int{}
	tmp := map[int]struct{}{}
	for _, port := range input {
		if _, ok := tmp[port]; !ok {
			tmp[port] = struct{}{}
			result = append(result, port)
		}
	}
	return result
}

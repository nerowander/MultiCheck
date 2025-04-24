package common

import (
	"bufio"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/nerowander/MultiCheck/config"
	"os"
	"strings"
)

func ParseInit(info *config.InfoScan) {
	parseUsername()
	parsePassword()
	parseUrl()
	parsePorts()
	parseUserInput(info)
}

func parseUsername() {
	if config.Username == "" && config.UsernameFile == "" {
		return
	}
	var Usernames []string
	if config.Username != "" {
		Usernames = strings.Split(config.Username, ",")
	}

	if config.UsernameFile != "" {
		users, err := readFile(config.UsernameFile)
		if err == nil {
			for _, user := range users {
				if user != "" {
					Usernames = append(Usernames, user)
				}
			}
		}
	}

	Usernames = removeDuplicateData(Usernames)
	for name := range config.UsernameDict {
		config.UsernameDict[name] = Usernames
	}
}
func parseUrl() {
	if config.URL != "" {
		urls := strings.Split(config.URL, ",")
		tmpUrls := make(map[string]struct{})
		for _, url := range urls {
			if _, ok := tmpUrls[url]; !ok {
				tmpUrls[url] = struct{}{}
				if url != "" {
					config.Urls = append(config.Urls, url)
				}
			}
		}
	}
	if config.URLFile != "" {
		urls, err := readFile(config.URLFile)
		if err == nil {
			TmpUrls := make(map[string]struct{})
			for _, url := range urls {
				if _, ok := TmpUrls[url]; !ok {
					TmpUrls[url] = struct{}{}
					if url != "" {
						config.Urls = append(config.Urls, url)
					}
				}
			}
		}
	}
}

func readFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Open %s error, %v\n", filename, err)
		os.Exit(0)
	}
	defer file.Close()
	var content []string
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if text != "" {
			content = append(content, scanner.Text())
		}
	}
	return content, nil
}

func parsePassword() {
	var PwdList []string
	if config.Password != "" {
		passwords := strings.Split(config.Password, ",")
		for _, pass := range passwords {
			if pass != "" {
				PwdList = append(PwdList, pass)
			}
		}
		config.PasswordDict = PwdList
	}
	if config.PasswordFile != "" {
		passwords, err := readFile(config.PasswordFile)
		if err == nil {
			for _, pass := range passwords {
				if pass != "" {
					PwdList = append(PwdList, pass)
				}
			}
			config.PasswordDict = PwdList
		}
	}
	if config.HashFile != "" {
		hashs, err := readFile(config.HashFile)
		if err == nil {
			for _, line := range hashs {
				if line == "" {
					continue
				}
				if len(line) == 32 {
					config.Hashs = append(config.Hashs, line)
				} else {
					fmt.Println("[-] " + line + "len(hash) != 32 ")
				}
			}
		}
	}

}

func parsePorts() {
	if config.PortsFile != "" {
		ports, err := readFile(config.PortsFile)
		if err == nil {
			newport := ""
			for _, port := range ports {
				if port != "" {
					newport += port + ","
				}
			}
			config.Ports = newport
		}
	}
}

func parseUserInput(info *config.InfoScan) {
	if info.Hosts == "" && config.HostFile == "" && config.URL == "" && config.URLFile == "" {
		fmt.Println("no hosts detected")
		flag.Usage()
		os.Exit(0)
	}

	if config.BruteThreads <= 0 {
		config.BruteThreads = 1
	}

	if config.Ports == config.DefaultPorts {
		config.Ports += "," + config.Webports
	}

	// extra ports
	if config.AddPorts != "" {
		config.Ports += "," + config.AddPorts
	}

	// extra username
	if config.AddUserNames != "" {
		user := strings.Split(config.AddUserNames, ",")
		for a := range config.UsernameDict {
			config.UsernameDict[a] = append(config.UsernameDict[a], user...)
			config.UsernameDict[a] = removeDuplicateData(config.UsernameDict[a])
		}
	}
	// extra password
	if config.AddPassWords != "" {
		pass := strings.Split(config.AddPassWords, ",")
		config.PasswordDict = append(config.PasswordDict, pass...)
		config.PasswordDict = removeDuplicateData(config.PasswordDict)
	}

	// todo: 后续添加proxy功能
	if config.Hash != "" && len(config.Hash) != 32 {
		fmt.Println("[-] input Hash is error,len(hash) must be 32")
		os.Exit(0)
	} else {
		config.Hashs = append(config.Hashs, config.Hash)
	}
	config.Hashs = removeDuplicateData(config.Hashs)
	for _, hash := range config.Hashs {
		hashByte, err := hex.DecodeString(hash)
		// hashbyte, err := hex.DecodeString(Hash)
		if err != nil {
			fmt.Println("[-] Hash is error,hex decode error ", hash)
			continue
		} else {
			config.HashBytes = append(config.HashBytes, hashByte)
		}
	}
	config.Hashs = []string{}
}
func removeDuplicateData(data []string) []string {
	result := []string{}
	temp := map[string]struct{}{}
	for _, item := range data {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

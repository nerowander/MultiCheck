package shell_cli

import (
	"fmt"
	"github.com/nerowander/MultiCheck/common"
	"github.com/nerowander/MultiCheck/config"
	"reflect"
	"strconv"
	"strings"
)

func handleInfoScanThreads(value interface{}) (result string) {
	v := reflect.ValueOf(value)

	if v.Kind() == reflect.String {
		intValue, err := strconv.Atoi(v.String())
		if err != nil {
			result = fmt.Sprintf("Invalid value for InfoScanThreads, type error: %v", err)
			return result
		}
		// 同时还要修改modules组的值
		config.Threads = intValue
		modules["Config"]["InfoScanThreads"] = strconv.Itoa(config.Threads)
		result = fmt.Sprintf("Set the value of InfoScanThreads: %d", config.Threads)
		return result
	} else if v.Kind() == reflect.Int {
		config.Threads = int(v.Int())
		modules["Config"]["InfoScanThreads"] = strconv.Itoa(config.Threads)
		result = fmt.Sprintf("Set the value of InfoScanThreads: %d", config.Threads)
		return result
	} else {
		result = fmt.Sprintf("Invalid value for InfoScanThreads, type error")
		return result
	}
}

func handleScanType(value interface{}) (result string) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String {
		config.ScanType = v.String()
		modules["Config"]["ScanType"] = strconv.Itoa(config.Threads)
		result = fmt.Sprintf("Set the value of ScanType: %s", config.ScanType)
		return result
	} else {
		result = fmt.Sprintf("Invalid value for ScanType, type error")
		return result
	}
}

func handleTimeout(value interface{}) (result string) {
	v := reflect.ValueOf(value)

	if v.Kind() == reflect.String {
		intValue, err := strconv.Atoi(v.String())
		if err != nil {
			result = fmt.Sprintf("Invalid value for Timeout, type error: %v", err)
			return result
		}
		// 同时还要修改modules组的值
		config.Timeout = int64(intValue)
		modules["Config"]["Timeout"] = strconv.Itoa(int(config.Timeout))
		result = fmt.Sprintf("Set the value of Timeout: %d", config.Timeout)
		return result
	} else if v.Kind() == reflect.Int {
		config.Timeout = v.Int()
		modules["Config"]["Timeout"] = strconv.Itoa(int(config.Timeout))
		result = fmt.Sprintf("Set the value of Timeout: %d", config.Timeout)
		return result
	} else {
		result = fmt.Sprintf("Invalid value for Timeout, type error")
		return result
	}
}

func handleCommand(value interface{}) (result string) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String {
		config.Command = v.String()
		modules["Config"]["Command"] = config.Command
		result = fmt.Sprintf("Set the value of Command: %s", config.Command)
		return result
	} else {
		result = fmt.Sprintf("Invalid value for Command, type error")
		return result
	}
}

func handlePorts(value interface{}) (result string) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String {
		config.Ports = v.String()
		modules["Config"]["Ports"] = config.Ports
		result = fmt.Sprintf("Set the value of Ports: %s", config.Ports)
		return result
	} else {
		result = fmt.Sprintf("Invalid value for Ports, type error")
		return result
	}
}

func handlePortsFile(value interface{}) (result string) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String {
		config.PortsFile = v.String()
		modules["Config"]["PortsFile"] = config.PortsFile
		result = fmt.Sprintf("Set the value of PortsFile: %s", config.PortsFile)
		return result
	} else {
		result = fmt.Sprintf("Invalid value for PortsFile, type error")
		return result
	}
}

func handleURL(value interface{}) (result string) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String {
		config.URL = v.String()
		modules["Config"]["URL"] = config.URL
		result = fmt.Sprintf("Set the value of URL: %s", config.URL)
		return result
	} else {
		result = fmt.Sprintf("Invalid value for URL, type error")
		return result
	}
}

func handleURLFile(value interface{}) (result string) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String {
		config.URLFile = v.String()
		modules["Config"]["URLFile"] = config.URLFile
		result = fmt.Sprintf("Set the value of URLFile: %s", config.URLFile)
		return result
	} else {
		result = fmt.Sprintf("Invalid value for URLFile, type error")
		return result
	}
}

func handleAddPorts(value interface{}) (result string) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String {
		config.AddPorts = v.String()
		modules["Config"]["AddPorts"] = config.AddPorts
		result = fmt.Sprintf("Set the value of AddPorts: %s", config.AddPorts)
		return result
	} else {
		result = fmt.Sprintf("Invalid value for AddPorts, type error")
		return result
	}
}

func handleAddUserNames(value interface{}) (result string) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String {
		config.AddUserNames = v.String()
		modules["Config"]["AddUserNames"] = config.AddUserNames
		result = fmt.Sprintf("Set the value of AddUserNames: %s", config.AddUserNames)
		return result
	} else {
		result = fmt.Sprintf("Invalid value for AddUserNames, type error")
		return result
	}
}

func handleAddPassWords(value interface{}) (result string) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String {
		config.AddPassWords = v.String()
		modules["Config"]["AddPassWords"] = config.AddPassWords
		result = fmt.Sprintf("Set the value of AddPassWords: %s", config.AddPassWords)
		return result
	} else {
		result = fmt.Sprintf("Invalid value for AddPassWords, type error")
		return result
	}
}

func handleHash(value interface{}) (result string) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String {
		config.Hash = v.String()
		modules["Config"]["Hash"] = config.Hash
		result = fmt.Sprintf("Set the value of Hash: %s", config.Hash)
		return result
	} else {
		result = fmt.Sprintf("Invalid value for Hash, type error")
		return result
	}
}

func handleBruteThreads(value interface{}) (result string) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String {
		intValue, err := strconv.Atoi(v.String())
		if err != nil {
			result = fmt.Sprintf("Invalid value for BruteThreads, type error: %v", err)
			return result
		}
		// 同时还要修改modules组的值
		config.BruteThreads = intValue
		modules["Config"]["BruteThreads"] = strconv.Itoa(config.BruteThreads)
		result = fmt.Sprintf("Set the value of BruteThreads: %d", config.BruteThreads)
		return result
	} else if v.Kind() == reflect.Int {
		config.BruteThreads = int(v.Int())
		modules["Config"]["BruteThreads"] = strconv.Itoa(config.BruteThreads)
		result = fmt.Sprintf("Set the value of BruteThreads: %d", config.BruteThreads)
		return result
	} else {
		result = fmt.Sprintf("Invalid value for BruteThreads, type error")
		return result
	}
}

func handlePocPath(value interface{}) (result string) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String {
		config.PocPath = v.String()
		modules["Config"]["PocPath"] = config.PocPath
		result = fmt.Sprintf("Set the value of PocPath: %s", config.PocPath)
		return result
	} else {
		result = fmt.Sprintf("Invalid value for PocPath, type error")
		return result
	}
}

func handleCookie(value interface{}) (result string) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String {
		config.Cookie = v.String()
		modules["Config"]["Cookie"] = config.Cookie
		result = fmt.Sprintf("Set the value of Cookie: %s", config.Cookie)
		return result
	} else {
		result = fmt.Sprintf("Invalid value for Cookie, type error")
		return result
	}
}

func handleWebTimeout(value interface{}) (result string) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String {
		intValue, err := strconv.Atoi(v.String())
		if err != nil {
			result = fmt.Sprintf("Invalid value for WebTimeout, type error: %v", err)
			return result
		}
		// 同时还要修改modules组的值
		config.WebTimeout = int64(intValue)
		modules["Config"]["BruteThreads"] = strconv.Itoa(int(config.WebTimeout))
		result = fmt.Sprintf("Set the value of WebTimeout: %d", config.WebTimeout)
		return result
	} else if v.Kind() == reflect.Int {
		config.WebTimeout = v.Int()
		modules["Config"]["WebTimeout"] = strconv.Itoa(int(config.WebTimeout))
		result = fmt.Sprintf("Set the value of WebTimeout: %d", config.WebTimeout)
		return result
	} else {
		result = fmt.Sprintf("Invalid value for WebTimeout, type error")
		return result
	}
}

func handleUserName(value interface{}) (result string) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String {
		config.Username = v.String()
		modules["Config"]["Username"] = config.Username
		result = fmt.Sprintf("Set the value of Username: %s", config.Username)
		return result
	} else {
		result = fmt.Sprintf("Invalid value for Username, type error")
		return result
	}
}

func handleUserNameFile(value interface{}) (result string) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String {
		config.UsernameFile = v.String()
		modules["Config"]["UsernameFile"] = config.UsernameFile
		result = fmt.Sprintf("Set the value of UsernameFile: %s", config.UsernameFile)
		return result
	} else {
		result = fmt.Sprintf("Invalid value for UsernameFile, type error")
		return result
	}
}

func handlePassWord(value interface{}) (result string) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String {
		config.Password = v.String()
		modules["Config"]["Password"] = config.Password
		result = fmt.Sprintf("Set the value of Password: %s", config.Password)
		return result
	} else {
		result = fmt.Sprintf("Invalid value for Password, type error")
		return result
	}
}

func handlePassWordFile(value interface{}) (result string) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String {
		config.PasswordFile = v.String()
		modules["Config"]["PasswordFile"] = config.PasswordFile
		result = fmt.Sprintf("Set the value of PasswordFile: %s", config.PasswordFile)
		return result
	} else {
		result = fmt.Sprintf("Invalid value for PasswordFile, type error")
		return result
	}
}

func handleHashFile(value interface{}) (result string) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String {
		config.HashFile = v.String()
		modules["Config"]["HashFile"] = config.HashFile
		result = fmt.Sprintf("Set the value of HashFile: %s", config.HashFile)
		return result
	} else {
		result = fmt.Sprintf("Invalid value for HashFile, type error")
		return result
	}
}

func handlePocNum(value interface{}) (result string) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String {
		intValue, err := strconv.Atoi(v.String())
		if err != nil {
			result = fmt.Sprintf("Invalid value for PocNum, type error: %v", err)
			return result
		}
		// 同时还要修改modules组的值
		config.PocNum = intValue
		modules["Config"]["PocNum"] = strconv.Itoa(config.PocNum)
		result = fmt.Sprintf("Set the value of PocNum: %d", config.PocNum)
		return result
	} else if v.Kind() == reflect.Int {
		config.PocNum = int(v.Int())
		modules["Config"]["PocNum"] = strconv.Itoa(config.PocNum)
		result = fmt.Sprintf("Set the value of PocNum: %d", config.PocNum)
		return result
	} else {
		result = fmt.Sprintf("Invalid value for PocNum, type error")
		return result
	}
}

func handlePocName(value interface{}) (result string) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String {
		config.PocName = v.String()
		modules["Config"]["PocName"] = config.PocName
		result = fmt.Sprintf("Set the value of PocName: %s", config.PocName)
		return result
	} else {
		result = fmt.Sprintf("Invalid value for PocName, type error")
		return result
	}
}

func handlePocType(value interface{}) (result string) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String {
		config.PocType = v.String()
		modules["Config"]["PocType"] = config.PocType
		result = fmt.Sprintf("Set the value of PocType: %s", config.PocType)
		return result
	} else {
		result = fmt.Sprintf("Invalid value for PocType, type error")
		return result
	}
}

func handleDnsLog(value interface{}) (result string) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String {
		// 如果是字符串，尝试将其转换为 bool
		strValue := strings.ToLower(v.String())
		if strValue == "true" {
			config.DnsLog = true
		} else if strValue == "false" {
			config.DnsLog = false
		} else {
			result = fmt.Sprintf("Invalid value for DnsLog, expected 'true' or 'false', got: %s", strValue)
			return result
		}
		modules["Config"]["DnsLog"] = strconv.FormatBool(config.DnsLog)
		result = fmt.Sprintf("Set the value of DnsLog: %v", config.DnsLog)
		return result
	} else if v.Kind() == reflect.Bool {
		// 如果是 bool 类型，直接赋值
		config.DnsLog = v.Bool()
		modules["Config"]["DnsLog"] = strconv.FormatBool(config.DnsLog)
		result = fmt.Sprintf("Set the value of DnsLog: %v", config.DnsLog)
		return result
	} else {
		result = fmt.Sprintf("Invalid value for DnsLog, type error")
		return result
	}
}

func handleCeyeToken(value interface{}) (result string) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String {
		config.CeyeToken = v.String()
		modules["Config"]["CeyeToken"] = config.CeyeToken
		result = fmt.Sprintf("Set the value of CeyeToken: %s", config.CeyeToken)
		return result
	} else {
		result = fmt.Sprintf("Invalid value for CeyeToken, type error")
		return result
	}
}

func handleCeyeURL(value interface{}) (result string) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String {
		config.CeyeURL = v.String()
		modules["Config"]["CeyeURL"] = config.CeyeURL
		result = fmt.Sprintf("Set the value of CeyeURL: %s", config.CeyeURL)
		return result
	} else {
		result = fmt.Sprintf("Invalid value for CeyeURL, type error")
		return result
	}
}

func handleFullPOC(value interface{}) (result string) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String {
		// 如果是字符串，尝试将其转换为 bool
		strValue := strings.ToLower(v.String())
		if strValue == "true" {
			config.FullPOC = true
		} else if strValue == "false" {
			config.FullPOC = false
		} else {
			result = fmt.Sprintf("Invalid value for FullPOC, expected 'true' or 'false', got: %s", strValue)
			return result
		}
		modules["Config"]["FullPOC"] = strconv.FormatBool(config.FullPOC)
		result = fmt.Sprintf("Set the value of FullPOC: %v", config.FullPOC)
		return result
	} else if v.Kind() == reflect.Bool {
		// 如果是 bool 类型，直接赋值
		config.FullPOC = v.Bool()
		modules["Config"]["FullPOC"] = strconv.FormatBool(config.FullPOC)
		result = fmt.Sprintf("Set the value of FullPOC: %v", config.FullPOC)
		return result
	} else {
		result = fmt.Sprintf("Invalid value for FullPOC, type error")
		return result
	}
}

func handleFullEXP(value interface{}) (result string) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String {
		// 如果是字符串，尝试将其转换为 bool
		strValue := strings.ToLower(v.String())
		if strValue == "true" {
			config.FullEXP = true
		} else if strValue == "false" {
			config.FullEXP = false
		} else {
			result = fmt.Sprintf("Invalid value for FullEXP, expected 'true' or 'false', got: %s", strValue)
			return result
		}
		modules["Config"]["FullEXP"] = strconv.FormatBool(config.FullEXP)
		result = fmt.Sprintf("Set the value of FullEXP: %v", config.FullEXP)
		return result
	} else if v.Kind() == reflect.Bool {
		// 如果是 bool 类型，直接赋值
		config.FullEXP = v.Bool()
		modules["Config"]["FullEXP"] = strconv.FormatBool(config.FullEXP)
		result = fmt.Sprintf("Set the value of FullEXP: %v", config.FullEXP)
		return result
	} else {
		result = fmt.Sprintf("Invalid value for FullEXP, type error")
		return result
	}
}

func handleExpNum(value interface{}) (result string) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String {
		intValue, err := strconv.Atoi(v.String())
		if err != nil {
			result = fmt.Sprintf("Invalid value for ExpNum, type error: %v", err)
			return result
		}
		// 同时还要修改modules组的值
		config.ExpNum = intValue
		modules["Config"]["ExpNum"] = strconv.Itoa(config.ExpNum)
		result = fmt.Sprintf("Set the value of ExpNum: %d", config.ExpNum)
		return result
	} else if v.Kind() == reflect.Int {
		config.ExpNum = int(v.Int())
		modules["Config"]["PocNum"] = strconv.Itoa(config.ExpNum)
		result = fmt.Sprintf("Set the value of ExpNum: %d", config.ExpNum)
		return result
	} else {
		result = fmt.Sprintf("Invalid value for ExpNum, type error")
		return result
	}
}

func handleExpPath(value interface{}) (result string) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String {
		config.ExpPath = v.String()
		modules["Config"]["ExpPath"] = config.ExpPath
		result = fmt.Sprintf("Set the value of ExpPath: %s", config.ExpPath)
		return result
	} else {
		result = fmt.Sprintf("Invalid value for ExpPath, type error")
		return result
	}
}

func handleExpType(value interface{}) (result string) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String {
		config.ExpType = v.String()
		modules["Config"]["ExpType"] = config.ExpType
		result = fmt.Sprintf("Set the value of ExpType: %s", config.ExpType)
		return result
	} else {
		result = fmt.Sprintf("Invalid value for ExpType, type error")
		return result
	}
}

func handleExpName(value interface{}) (result string) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String {
		config.ExpName = v.String()
		modules["Config"]["ExpName"] = config.ExpName
		result = fmt.Sprintf("Set the value of ExpName: %s", config.ExpName)
		return result
	} else {
		result = fmt.Sprintf("Invalid value for ExpName, type error")
		return result
	}
}

func handleSaveResult(value interface{}) (result string) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String {
		// 如果是字符串，尝试将其转换为 bool
		strValue := strings.ToLower(v.String())
		if strValue == "true" {
			config.SaveResult = true
		} else if strValue == "false" {
			config.SaveResult = false
		} else {
			result = fmt.Sprintf("Invalid value for SaveResult, expected 'true' or 'false', got: %s", strValue)
			return result
		}
		modules["Config"]["SaveResult"] = strconv.FormatBool(config.SaveResult)
		result = fmt.Sprintf("Set the value of SaveResult: %v", config.SaveResult)
		return result
	} else if v.Kind() == reflect.Bool {
		// 如果是 bool 类型，直接赋值
		config.SaveResult = v.Bool()
		modules["Config"]["SaveResult"] = strconv.FormatBool(config.SaveResult)
		result = fmt.Sprintf("Set the value of SaveResult: %v", config.SaveResult)
		return result
	} else {
		result = fmt.Sprintf("Invalid value for SaveResult, type error")
		return result
	}
}

func handleOutPutFile(value interface{}) (result string) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String {
		config.OutPutFile = v.String()
		modules["Config"]["OutPutFile"] = config.OutPutFile
		result = fmt.Sprintf("Set the value of OutPutFile: %s", config.OutPutFile)
		return result
	} else {
		result = fmt.Sprintf("Invalid value for OutPutFile, type error")
		return result
	}
}

func handleHosts(value interface{}) (result string) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String {
		infoScan.Hosts = v.String()
		modules["InformationScan"]["Hosts"] = infoScan.Hosts
		result = fmt.Sprintf("Set the value of Hosts: %s", infoScan.Hosts)
		return result
	} else {
		result = fmt.Sprintf("Invalid value for Hosts, type error")
		return result
	}
}

func handleHostFile(value interface{}) (result string) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String {
		config.HostFile = v.String()
		modules["Config"]["HostFile"] = config.HostFile
		result = fmt.Sprintf("Set the value of HostFile: %s", config.HostFile)
		return result
	} else {
		result = fmt.Sprintf("Invalid value for HostFile, type error")
		return result
	}
}

func handleBrute(value interface{}) (result string) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String {
		// 如果是字符串，尝试将其转换为 bool
		strValue := strings.ToLower(v.String())
		if strValue == "true" {
			infoScan.Brute = true
		} else if strValue == "false" {
			infoScan.Brute = false
		} else {
			result = fmt.Sprintf("Invalid value for Brute, expected 'true' or 'false', got: %s", strValue)
			return result
		}
		modules["InformationScan"]["Brute"] = strconv.FormatBool(infoScan.Brute)
		result = fmt.Sprintf("Set the value of Brute: %v", infoScan.Brute)
		return result
	} else if v.Kind() == reflect.Bool {
		// 如果是 bool 类型，直接赋值
		infoScan.Brute = v.Bool()
		modules["InformationScan"]["Brute"] = strconv.FormatBool(infoScan.Brute)
		result = fmt.Sprintf("Set the value of Brute: %v", infoScan.Brute)
		return result
	} else {
		result = fmt.Sprintf("Invalid value for Brute, type error")
		return result
	}
}

func handleFTPReadFile(value interface{}) (result string) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String {
		infoScan.FTPReadFile = v.String()
		modules["InformationScan"]["FTPReadFile"] = infoScan.FTPReadFile
		result = fmt.Sprintf("Set the value of FTPReadFile: %s", infoScan.FTPReadFile)
		return result
	} else {
		result = fmt.Sprintf("Invalid value for FTPReadFile, type error")
		return result
	}
}

func handleFTPWriteFile(value interface{}) (result string) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String {
		infoScan.FTPWriteFile = v.String()
		modules["InformationScan"]["FTPWriteFile"] = infoScan.FTPWriteFile
		result = fmt.Sprintf("Set the value of FTPWriteFile: %s", infoScan.FTPWriteFile)
		return result
	} else {
		result = fmt.Sprintf("Invalid value for FTPWriteFile, type error")
		return result
	}
}

func handleDomain(value interface{}) (result string) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String {
		infoScan.Domain = v.String()
		modules["InformationScan"]["Domain"] = infoScan.Domain
		result = fmt.Sprintf("Set the value of Domain: %s", infoScan.Domain)
		return result
	} else {
		result = fmt.Sprintf("Invalid value for Domain, type error")
		return result
	}
}

func handleSkipRedis(value interface{}) (result string) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String {
		// 如果是字符串，尝试将其转换为 bool
		strValue := strings.ToLower(v.String())
		if strValue == "true" {
			infoScan.SkipRedis = true
		} else if strValue == "false" {
			infoScan.SkipRedis = false
		} else {
			result = fmt.Sprintf("Invalid value for SkipRedis, expected 'true' or 'false', got: %s", strValue)
			return result
		}
		modules["InformationScan"]["SkipRedis"] = strconv.FormatBool(infoScan.SkipRedis)
		result = fmt.Sprintf("Set the value of SkipRedis: %v", infoScan.SkipRedis)
		return result
	} else if v.Kind() == reflect.Bool {
		// 如果是 bool 类型，直接赋值
		infoScan.SkipRedis = v.Bool()
		modules["InformationScan"]["SkipRedis"] = strconv.FormatBool(infoScan.SkipRedis)
		result = fmt.Sprintf("Set the value of SkipRedis: %v", infoScan.SkipRedis)
		return result
	} else {
		result = fmt.Sprintf("Invalid value for SkipRedis, type error")
		return result
	}
}

func handleRedisSshFile(value interface{}) (result string) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String {
		infoScan.RedisSshFile = v.String()
		modules["InformationScan"]["RedisSshFile"] = infoScan.RedisSshFile
		result = fmt.Sprintf("Set the value of RedisSshFile: %s", infoScan.RedisSshFile)
		return result
	} else {
		result = fmt.Sprintf("Invalid value for RedisSshFile, type error")
		return result
	}
}

func handleRedisCronHost(value interface{}) (result string) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String {
		infoScan.RedisCronHost = v.String()
		modules["InformationScan"]["RedisCronHost"] = infoScan.RedisCronHost
		result = fmt.Sprintf("Set the value of RedisCronHost: %s", infoScan.RedisCronHost)
		return result
	} else {
		result = fmt.Sprintf("Invalid value for RedisCronHost, type error")
		return result
	}
}

func handleRedisWebshellFile(value interface{}) (result string) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String {
		infoScan.RedisWebshellFile = v.String()
		modules["InformationScan"]["RedisWebshellFile"] = infoScan.RedisWebshellFile
		result = fmt.Sprintf("Set the value of RedisWebshellFile: %s", infoScan.RedisWebshellFile)
		return result
	} else {
		result = fmt.Sprintf("Invalid value for RedisWebshellFile, type error")
		return result
	}
}

func handleRemotePath(value interface{}) (result string) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String {
		infoScan.RemotePath = v.String()
		modules["InformationScan"]["RemotePath"] = infoScan.RemotePath
		result = fmt.Sprintf("Set the value of RemotePath: %s", infoScan.RemotePath)
		return result
	} else {
		result = fmt.Sprintf("Invalid value for RemotePath, type error")
		return result
	}
}

func handleSshKey(value interface{}) (result string) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String {
		infoScan.SshKey = v.String()
		modules["InformationScan"]["SshKey"] = infoScan.SshKey
		result = fmt.Sprintf("Set the value of SshKey: %s", infoScan.SshKey)
		return result
	} else {
		result = fmt.Sprintf("Invalid value for SshKey, type error")
		return result
	}
}

func handleLogWaitTime(value interface{}) (result string) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String {
		intValue, err := strconv.Atoi(v.String())
		if err != nil {
			result = fmt.Sprintf("Invalid value for LogWaitTime, type error: %v", err)
			return result
		}
		// 同时还要修改modules组的值
		common.LogWaitTime = int64(intValue)
		modules["Common"]["LogWaitTime"] = strconv.Itoa(int(common.LogWaitTime))
		result = fmt.Sprintf("Set the value of LogWaitTime: %d", common.LogWaitTime)
		return result
	} else if v.Kind() == reflect.Int {
		common.LogWaitTime = v.Int()
		modules["Common"]["LogWaitTime"] = strconv.Itoa(int(common.LogWaitTime))
		result = fmt.Sprintf("Set the value of LogWaitTime: %d", common.LogWaitTime)
		return result
	} else {
		result = fmt.Sprintf("Invalid value for LogWaitTime, type error")
		return result
	}
}

func handlePrintLog(value interface{}) (result string) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String {
		// 如果是字符串，尝试将其转换为 bool
		strValue := strings.ToLower(v.String())
		if strValue == "true" {
			common.PrintLog = true
		} else if strValue == "false" {
			common.PrintLog = false
		} else {
			result = fmt.Sprintf("Invalid value for PrintLog, expected 'true' or 'false', got: %s", strValue)
			return result
		}
		modules["Common"]["PrintLog"] = strconv.FormatBool(common.PrintLog)
		result = fmt.Sprintf("Set the value of PrintLog: %v", common.PrintLog)
		return result
	} else if v.Kind() == reflect.Bool {
		// 如果是 bool 类型，直接赋值
		common.PrintLog = v.Bool()
		modules["Common"]["PrintLog"] = strconv.FormatBool(common.PrintLog)
		result = fmt.Sprintf("Set the value of PrintLog: %v", common.PrintLog)
		return result
	} else {
		result = fmt.Sprintf("Invalid value for PrintLog, type error")
		return result
	}
}

func handleSaveLogToJSON(value interface{}) (result string) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String {
		// 如果是字符串，尝试将其转换为 bool
		strValue := strings.ToLower(v.String())
		if strValue == "true" {
			common.SaveLogToJSON = true
		} else if strValue == "false" {
			common.SaveLogToJSON = false
		} else {
			result = fmt.Sprintf("Invalid value for SaveLogToJSON, expected 'true' or 'false', got: %s", strValue)
			return result
		}
		modules["Common"]["SaveLogToJSON"] = strconv.FormatBool(common.SaveLogToJSON)
		result = fmt.Sprintf("Set the value of SaveLogToJSON: %v", common.SaveLogToJSON)
		return result
	} else if v.Kind() == reflect.Bool {
		// 如果是 bool 类型，直接赋值
		common.SaveLogToJSON = v.Bool()
		modules["Common"]["SaveLogToJSON"] = strconv.FormatBool(common.SaveLogToJSON)
		result = fmt.Sprintf("Set the value of SaveLogToJSON: %v", common.SaveLogToJSON)
		return result
	} else {
		result = fmt.Sprintf("Invalid value for SaveLogToJSON, type error")
		return result
	}
}

func handleSaveLogToHTML(value interface{}) (result string) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String {
		// 如果是字符串，尝试将其转换为 bool
		strValue := strings.ToLower(v.String())
		if strValue == "true" {
			common.SaveLogToHTML = true
		} else if strValue == "false" {
			common.SaveLogToHTML = false
		} else {
			result = fmt.Sprintf("Invalid value for SaveLogToHTML, expected 'true' or 'false', got: %s", strValue)
			return result
		}
		modules["Common"]["SaveLogToHTML"] = strconv.FormatBool(common.SaveLogToHTML)
		result = fmt.Sprintf("Set the value of SaveLogToHTML: %v", common.SaveLogToHTML)
		return result
	} else if v.Kind() == reflect.Bool {
		// 如果是 bool 类型，直接赋值
		common.SaveLogToHTML = v.Bool()
		modules["Common"]["SaveLogToHTML"] = strconv.FormatBool(common.SaveLogToHTML)
		result = fmt.Sprintf("Set the value of SaveLogToHTML: %v", common.SaveLogToHTML)
		return result
	} else {
		result = fmt.Sprintf("Invalid value for SaveLogToHTML, type error")
		return result
	}
}

func handleEnableInfoContainer(value interface{}) (result string) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String {
		// 如果是字符串，尝试将其转换为 bool
		strValue := strings.ToLower(v.String())
		if strValue == "true" {
			config.EnableInfoContainer = true
		} else if strValue == "false" {
			config.EnableInfoContainer = false
		} else {
			result = fmt.Sprintf("Invalid value for EnableInfoContainer, expected 'true' or 'false', got: %s", strValue)
			return result
		}
		modules["Config"]["EnableInfoContainer"] = strconv.FormatBool(config.EnableInfoContainer)
		result = fmt.Sprintf("Set the value of EnableInfoContainer: %v", config.EnableInfoContainer)
		return result
	} else if v.Kind() == reflect.Bool {
		// 如果是 bool 类型，直接赋值
		config.EnableInfoContainer = v.Bool()
		modules["Config"]["EnableInfoContainer"] = strconv.FormatBool(config.EnableInfoContainer)
		result = fmt.Sprintf("Set the value of EnableInfoContainer: %v", config.EnableInfoContainer)
		return result
	} else {
		result = fmt.Sprintf("Invalid value for EnableInfoContainer, type error")
		return result
	}
}

func handleVulContainer(value interface{}) (result string) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String {
		// 如果是字符串，尝试将其转换为 bool
		strValue := strings.ToLower(v.String())
		if strValue == "true" {
			config.EnableVulContainer = true
		} else if strValue == "false" {
			config.EnableVulContainer = false
		} else {
			result = fmt.Sprintf("Invalid value for EnableVulContainer, expected 'true' or 'false', got: %s", strValue)
			return result
		}
		modules["Config"]["EnableVulContainer"] = strconv.FormatBool(config.EnableVulContainer)
		result = fmt.Sprintf("Set the value of EnableVulContainer: %v", config.EnableVulContainer)
		return result
	} else if v.Kind() == reflect.Bool {
		// 如果是 bool 类型，直接赋值
		config.EnableVulContainer = v.Bool()
		modules["Config"]["EnableVulContainer"] = strconv.FormatBool(config.EnableVulContainer)
		result = fmt.Sprintf("Set the value of EnableVulContainer: %v", config.EnableVulContainer)
		return result
	} else {
		result = fmt.Sprintf("Invalid value for EnableVulContainer, type error")
		return result
	}
}

// 提供参数描述
func getDescription(param string) string {
	descriptions := map[string]string{
		"InfoScanThreads":     "The nums of InfoScan thread",
		"ScanType":            "Choose scan type, for example: web | all | pocscan | exploit | ping | portscan",
		"Timeout":             "Set the InfoScan timeout",
		"Command":             "Exec command, for example: ssh",
		"Ports":               "Port to InfoScan, for example: 80 | 80,443 | 80-443",
		"PortsFile":           "Select port file to InfoScan",
		"URL":                 "Url to scan",
		"URLFile":             "Choose url file to scan",
		"AddPorts":            "Add extra ports,-pa 2333",
		"AddUserNames":        "Add extra username,-au admin",
		"AddPassWords":        "Add extra password,-apw admin",
		"Hash":                "Input hash",
		"BruteThreads":        "The nums of thread to brute",
		"PocPath":             "Poc file path",
		"Cookie":              "Set cookie",
		"WebTimeout":          "Set the Webclient timeout",
		"Username":            "Username account, for example: admin| admin,test",
		"UsernameFile":        "Username file",
		"Password":            "Password, for example: admin| admin,test",
		"PasswordFile":        "Password file",
		"HashFile":            "Hash file",
		"PocNum":              "Poc rate",
		"PocName":             "Choose poc name to scan, such as: sql",
		"PocType":             "Choose poc type: base(such as sql injection scan) | software (OA CMS scan) | iot(IOT scan) | all(web+iot scan)",
		"DnsLog":              "Enable ceye dnslog, then provide your ceye token and ceye url: set CeyeToken <ceye_token> | set CeyeURL <ceye_url>",
		"CeyeToken":           "Ceye token",
		"CeyeURL":             "Ceye url",
		"FullPOC":             "Poc full scan,as: use all shiro keys to scan",
		"FullEXP":             "Exp full scan,as: use all shiro keys to scan",
		"ExpNum":              "Exp rate",
		"ExpPath":             "Exp file path",
		"ExpType":             "Choose exp type: base(such as sql injection exploit) | software (OA CMS exploit) | iot(IOT exploit) | all(web+iot exploit)",
		"ExpName":             "Choose poc name to scan, such as: sql",
		"SaveResult":          "Save scan result to log file: json| html | txt(default)",
		"OutPutFile":          "Output log file name",
		"Hosts":               "IP address to scan, for example: 192.168.0.1 | 1.2.3.4/24 | 192.168.0.1,192.168.0.2 | 1.2.3.4-255 | 1.2.3.4-1.2.3.255",
		"HostFile":            "Select host file to scan",
		"Brute":               "Brute password or not",
		"FTPReadFile":         "Choose a file to read in ftp server",
		"FTPWriteFile":        "Choose a file to write in ftp server",
		"Domain":              "Smb | rdp domain",
		"SkipRedis":           "Skip redis exploit",
		"RedisSshFile":        "Redis file to write sshkey file, for example: id_rsa.pub",
		"RedisCronHost":       "Redis shell to write cron file for example: 192.168.1.1:2333",
		"RedisWebshellFile":   "Redis file to write webshell file, for example: shell.php",
		"RemotePath":          "Remote path of the target, for example: FastCGI filepath",
		"SshKey":              "Sshkey file (id_rsa)",
		"LogWaitTime":         "Logerr wait time",
		"PrintLog":            "Print scan log",
		"SaveLogToJSON":       "Save log to json file",
		"SaveLogToHTML":       "Save log to html file",
		"EnableInfoContainer": "Enable docker container to run scan",
	}
	return descriptions[param]
}

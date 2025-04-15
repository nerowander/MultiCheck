package shell_cli

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/nerowander/MultiCheck/InformationScan/Modules"
	"github.com/nerowander/MultiCheck/common"
	"github.com/nerowander/MultiCheck/config"
	"os"
	"strconv"
	"strings"
)

var infoScan config.InfoScan

func init() {
	infoScan.Brute = true
	infoScan.SkipRedis = true
}

// 定义模块参数
var modules = map[string]map[string]string{
	"Config": {
		"InfoScanThreads":        strconv.Itoa(config.Threads),
		"ScanType":               config.ScanType,
		"Timeout":                strconv.FormatInt(config.Timeout, 10),
		"Command":                config.Command,
		"HostFile":               config.HostFile,
		"Ports":                  config.Ports,
		"PortsFile":              config.PortsFile,
		"URL":                    config.URL,
		"URLFile":                config.URLFile,
		"AddPorts":               config.AddPorts,
		"AddUserNames":           config.AddUserNames,
		"AddPassWords":           config.AddPassWords,
		"Hash":                   config.Hash,
		"BruteThreads":           strconv.Itoa(config.BruteThreads),
		"PocPath":                config.PocPath,
		"Cookie":                 config.Cookie,
		"WebTimeout":             strconv.FormatInt(config.WebTimeout, 10),
		"Username":               config.Username,
		"UsernameFile":           config.UsernameFile,
		"Password":               config.Password,
		"PasswordFile":           config.PasswordFile,
		"HashFile":               config.HashFile,
		"PocNum":                 strconv.Itoa(config.PocNum),
		"PocName":                config.PocName,
		"PocType":                config.PocType,
		"DnsLog":                 strconv.FormatBool(config.DnsLog),
		"CeyeToken":              config.CeyeToken,
		"CeyeURL":                config.CeyeURL,
		"FullPOC":                strconv.FormatBool(config.FullPOC),
		"FullEXP":                strconv.FormatBool(config.FullEXP),
		"ExpNum":                 strconv.Itoa(config.ExpNum),
		"ExpPath":                config.ExpPath,
		"ExpType":                config.ExpType,
		"ExpName":                config.ExpName,
		"SaveResult":             strconv.FormatBool(config.SaveResult),
		"OutPutFile":             config.OutPutFile,
		"EnableContainer":        strconv.FormatBool(config.EnableInfoContainer),
		"EnablePocContainer":     strconv.FormatBool(config.EnablePocContainer),
		"EnableExploitContainer": strconv.FormatBool(config.EnableExploitContainer),
	},
	"InformationScan": {
		"Hosts":             infoScan.Hosts,
		"Brute":             strconv.FormatBool(infoScan.Brute),
		"FTPReadFile":       infoScan.FTPReadFile,
		"FTPWriteFile":      infoScan.FTPWriteFile,
		"Domain":            infoScan.Domain,
		"SkipRedis":         strconv.FormatBool(infoScan.SkipRedis),
		"RedisSshFile":      infoScan.RedisSshFile,
		"RedisCronHost":     infoScan.RedisCronHost,
		"RedisWebshellFile": infoScan.RedisWebshellFile,
		"RemotePath":        infoScan.RemotePath,
		"SshKey":            infoScan.SshKey,
	},
	//"PocScan":     {
	//	"Target":
	//},
	//"ExploitScan": {},
	"Common": {
		"LogWaitTime":   strconv.FormatInt(common.LogWaitTime, 10),
		"PrintLog":      strconv.FormatBool(common.PrintLog),
		"SaveLogToJSON": strconv.FormatBool(common.SaveLogToJSON),
		"SaveLogToHTML": strconv.FormatBool(common.SaveLogToHTML),
	},
}

var activeModule string = ""
var cliPrefix = "MutiCheck > "

func Start() {
	Execute()
}

func Execute() {
	ppt := prompt.New(
		executor,
		completer,
		prompt.OptionTitle("MyCLI"),
		prompt.OptionPrefix(cliPrefix), // 默认前缀
	)
	ppt.Run()
}

// 执行用户输入
func executor(input string) {
	input = strings.TrimSpace(input)
	args := strings.Split(input, " ")

	if len(args) == 0 {
		return
	}

	switch args[0] {
	case "exit":
		fmt.Println("Exiting...")
		os.Exit(0)
	case "use": // 使用某一模块
		if args[1] == "infoscan" {
			activeModule = "InformationScan"
			cliPrefix = "MultiCheck (InformationScan) > "
			fmt.Println("Selected module: InformationScan")
			Execute()
		} else {
			fmt.Println("Unsupported module: ", args[1])
		}
	case "show":
		if activeModule == "" {
			fmt.Println("No module selected. Use a module first.")
			return
		} else if len(args) > 1 && args[1] == "options" {
			showOptions()
		} else {
			fmt.Println("show options format error, for example: show options", input)
		}
	case "set":
		// set A 1
		if len(args) == 3 {
			results := getArgsVariable(args[1], args[2])
			fmt.Println(results)
		} else {
			fmt.Println("Set config format error, for example: set Hosts 192.168.1.1")
		}
	case "start":
		// 开始某一模块的任务
		if len(args) == 2 {
			if args[1] == "infoscan" {
				Modules.HostScan(&infoScan)
				common.GetSugestions()
				os.Exit(0)
			} else {
				fmt.Println("Start task format error, for example: start infoscan")
			}
		} else {
			fmt.Println("Start task format error, for example: start infoscan")
		}
	default:
		fmt.Printf("Unsupported command: %s \n", input)
	}
}

func showOptions() {
	if activeModule == "" {
		fmt.Println("No module selected. Use a module first.")
		return
	}
	fmt.Println(config.Threads)
	fmt.Println("\nconfig options:\n")
	fmt.Println("   Name       Current settings     Description                     ")
	fmt.Println("   ----       ----------------     -----------                     ")

	for key, value := range modules["Config"] {
		fmt.Printf("   %-14s %-20s %s\n", key, value, getDescription(key))
	}

	fmt.Println("\ncommon options:\n")
	fmt.Println("   Name       Current settings     Description                     ")
	fmt.Println("   ----       ----------------     -----------                     ")
	for key, value := range modules["Common"] {
		fmt.Printf("   %-14s %-20s %s\n", key, value, getDescription(key))
	}

	fmt.Printf("\n%s options:\n", activeModule)
	fmt.Println("   Name           Current settings     Description                           ")
	fmt.Println("   ----           ----------------     -----------                           ")

	for key, value := range modules[activeModule] {
		fmt.Printf("   %-14s %-20s %s\n", key, value, getDescription(key))
	}
	fmt.Println()
}

func getArgsVariable(param string, value interface{}) (result string) {
	typeHandlers := map[string]func(value interface{}) string{
		"InfoScanThreads":   handleInfoScanThreads,
		"ScanType":          handleScanType,
		"Timeout":           handleTimeout,
		"Command":           handleCommand,
		"Ports":             handlePorts,
		"PortsFile":         handlePortsFile,
		"URL":               handleURL,
		"URLFile":           handleURLFile,
		"AddPorts":          handleAddPorts,
		"AddUserNames":      handleAddUserNames,
		"AddPassWords":      handleAddPassWords,
		"Hash":              handleHash,
		"BruteThreads":      handleBruteThreads,
		"PocPath":           handlePocPath,
		"Cookie":            handleCookie,
		"WebTimeout":        handleWebTimeout,
		"Username":          handleUserName,
		"UsernameFile":      handleUserNameFile,
		"Password":          handlePassWord,
		"PasswordFile":      handlePassWordFile,
		"HashFile":          handleHashFile,
		"PocNum":            handlePocNum,
		"PocName":           handlePocName,
		"PocType":           handlePocType,
		"DnsLog":            handleDnsLog,
		"CeyeToken":         handleCeyeToken,
		"CeyeURL":           handleCeyeURL,
		"FullPOC":           handleFullPOC,
		"FullEXP":           handleFullEXP,
		"ExpNum":            handleExpNum,
		"ExpPath":           handleExpPath,
		"ExpType":           handleExpType,
		"ExpName":           handleExpName,
		"SaveResult":        handleSaveResult,
		"OutPutFile":        handleOutPutFile,
		"Hosts":             handleHosts,
		"HostFile":          handleHostFile,
		"Brute":             handleBrute,
		"FTPReadFile":       handleFTPReadFile,
		"FTPWriteFile":      handleFTPWriteFile,
		"Domain":            handleDomain,
		"SkipRedis":         handleSkipRedis,
		"RedisSshFile":      handleRedisSshFile,
		"RedisCronHost":     handleRedisCronHost,
		"RedisWebshellFile": handleRedisWebshellFile,
		"RemotePath":        handleRemotePath,
		"SshKey":            handleSshKey,
		"LogWaitTime":       handleLogWaitTime,
		"PrintLog":          handlePrintLog,
		"SaveLogToJSON":     handleSaveLogToJSON,
		"SaveLogToHTML":     handleSaveLogToHTML,
		"EnableContainer":   handleEnableInfoContainer,
	}

	// 查找处理函数
	if handler, exists := typeHandlers[param]; exists {
		// 执行处理函数
		return handler(value)
	}

	// 如果找不到对应的 param
	result = fmt.Sprintf("Invalid param: %s, no such config", param)
	return result
}

func completer(d prompt.Document) []prompt.Suggest {
	suggestions := []prompt.Suggest{
		{Text: "use infoscan", Description: "Enable InformationScan module"},
		{Text: "use pocscan", Description: "Enable PocScan module"},
		{Text: "use exploit", Description: "Enable ExploitScan module"},
		{Text: "show options", Description: "Show module options"},
		{Text: "exit", Description: "Exit the System"},
		{Text: "set", Description: "set config variables, for example: set InfoScanThreads 100"},
		{Text: "start", Description: "start task, for example: start infoscan"},
	}
	return prompt.FilterHasPrefix(suggestions, d.TextBeforeCursor(), true)
}

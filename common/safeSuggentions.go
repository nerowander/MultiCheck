package common

import (
	"github.com/fatih/color"
	"github.com/nerowander/MultiCheck/config"
	"strings"
)

var suggestions []string

// 扫描完后根据文件日志或输出日志获取安全修复建议
// 找一个变量存储所有的日志 -> 变量前缀去重 -> 匹配 -> 输出（建议的字符串进行去重）

func GetSugestions() {
	checkDatas := removeDuplicateString(config.ResLogs)
	for _, checkData := range checkDatas {
		suggestions = append(suggestions, getSuggestion(checkData))
	}
	resSuggestions := removeDuplicateString(suggestions)
	for _, resSuggestion := range resSuggestions {
		color.Yellow(resSuggestion)
	}
}

func removeDuplicateString(strs []string) []string {
	uniqueMap := make(map[string]struct{})
	var uniqueStrings []string
	for _, str := range strs {
		if _, exists := uniqueMap[str]; !exists {
			uniqueMap[str] = struct{}{} // 使用空结构体占位
			uniqueStrings = append(uniqueStrings, str)
		}
	}
	return uniqueStrings
}

// todo: 补充建议模块
// base+iot+poc模块的安全修复建议
func getSuggestion(log string) string {
	switch {
	case strings.HasPrefix(log, "[+] Redis"):
		return "[!] Redis未授权修复建议：1、启用身份验证机制 2、可限制Redis请求连接来源IP，并开启protect-mode，不使用root启动redis"
	case strings.HasPrefix(log, "[+] SSH"):
		return "[!] SSH弱密码修复建议：1、加强密码强度 2、最好禁用SSH密码登陆，改为仅可使用密钥验证，并禁用root登陆"
	case strings.HasPrefix(log, "[+] mysql"):
		return "[!] MySQL弱密码修复建议：1、加强密码强度 2、可限制MySQL请求连接来源IP，并限制普通用户权限，一般情况下也可禁用root登陆以及不使用root启动MySQL"
	case strings.Contains(log, "base-exp-command-injection") || strings.Contains(log, "base-poc-command-injection"):
		return "[!] 命令注入漏洞修复建议：1、黑白名单、WAF等防御策略过滤用户输入数据"
	default:
		return ""
	}
}

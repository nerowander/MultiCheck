package WebScan

import (
	"fmt"
	"github.com/nerowander/MultiCheck/WebScan/rules"
	"github.com/nerowander/MultiCheck/common"
	"regexp"
)

type CheckDatas struct {
	Body    []byte
	Headers string
}

// webinfo匹配
func InfoCheck(Url string, CheckData *[]CheckDatas) []string {
	var matched bool
	var infoname []string

	for _, data := range *CheckData {
		for _, rule := range rules.RuleDatas {
			if rule.Type == "code" {
				// body匹配
				matched, _ = regexp.MatchString(rule.Rule, string(data.Body))
			} else {
				// headers匹配
				matched, _ = regexp.MatchString(rule.Rule, data.Headers)
			}
			if matched == true {
				infoname = append(infoname, rule.Name)
			}
		}
	}

	infoname = removeDuplicateElement(infoname)

	if len(infoname) > 0 {
		result := fmt.Sprintf("[+] InfoScan %-25v %s ", Url, infoname)
		common.LogSuccess(result)
		return infoname
	}
	return []string{""}
}

func removeDuplicateElement(languages []string) []string {
	result := make([]string, 0, len(languages))
	temp := map[string]struct{}{}
	for _, item := range languages {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

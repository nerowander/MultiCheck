package poclib

import (
	"InformationScan/common"
	"InformationScan/config"
	"fmt"
	"github.com/google/cel-go/common/types/ref"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type Task struct {
	Req *http.Request
	Poc *Poc
}

var C CustomLib

func init() {
	C = NewEnvOption()
}
func CheckMultiPoc(req *http.Request, pocs []*Poc, workers int) {
	// tasks := make(chan Task, len(pocs))
	tasks := make(chan Task)
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		go func() {
			for task := range tasks {
				isVul, _, name := executePoc(task.Req, task.Poc)
				if isVul {
					result := fmt.Sprintf("[+] PocScan %s %s %s", task.Req.URL, task.Poc.Name, name)
					common.LogSuccess(result)
				}
				wg.Done()
			}
		}()
	}
	for _, poc := range pocs {
		task := Task{
			Req: req,
			Poc: poc,
		}
		wg.Add(1)
		tasks <- task
	}
	wg.Wait()
	close(tasks)
}

func executePoc(oReq *http.Request, p *Poc) (bool, error, string) {
	// 后续看要不要加锁
	//mu.Lock()         // 加锁，防止并发修改 c
	//defer mu.Unlock() // 任务完成后解锁
	//c = NewEnvOption()
	// set为poc yaml中的表达式set字段
	// sets为poc yaml当中的爆破字段
	C.UpdateCompileOptions(p.Set)
	if len(p.Sets) > 0 {
		var setMap StrMap
		for _, item := range p.Sets {
			if len(item.Value) > 0 {
				setMap = append(setMap, StrItem{item.Key, item.Value[0]})
			} else {
				setMap = append(setMap, StrItem{item.Key, ""})
			}
		}
		C.UpdateCompileOptions(setMap)
	}
	env, err := NewEnv(&C)
	if err != nil {
		fmt.Printf("[-] %s environment creation error: %s\n", p.Name, err)
		return false, err, ""
	}
	var req *Request
	req, err = ParseRequest(oReq)
	if err != nil {
		fmt.Printf("[-] %s ParseRequest error: %s\n", p.Name, err)
		return false, err, ""
	}

	variableMap := make(map[string]interface{})
	defer func() { variableMap = nil }()
	variableMap["request"] = req

	for _, item := range p.Set {
		k, expression := item.Key, item.Value
		// for example
		// set:
		// reverse: newReverse()
		// reverseURL: reverse.url
		if expression == "newReverse()" {
			// dnslog反连
			if !config.DnsLog {
				return false, nil, ""
			}
			variableMap[k] = NewReverse()
			continue
		}
		err, _ = EvalSet(env, variableMap, k, expression)
		if err != nil {
			fmt.Printf("[-] %s evalset error: %v\n", p.Name, err)
		}
	}

	success := false
	//爆破模式，比如弱口令爆破，敏感目录扫描等
	if len(p.Sets) > 0 {
		success, err = brutePoc(oReq, p, variableMap, req, env)
		return success, nil, ""
	}

	//非爆破模式
	DealWithRule := func(name string, rule Rules) (bool, error) {
		Headers := MapClone(rule.Headers)
		var (
			flag, ok bool
		)
		for k1, v1 := range variableMap {
			_, isMap := v1.(map[string]string)
			if isMap {
				continue
			}
			value := fmt.Sprintf("%v", v1)
			for k2, v2 := range Headers {
				if !strings.Contains(v2, "{{"+k1+"}}") {
					continue
				}
				Headers[k2] = strings.ReplaceAll(v2, "{{"+k1+"}}", value)
			}
			// 针对于web base的情况替换path和body
			// 目前支持替换path和body、expression
			// 考虑checkwebshellpath的检查
			if strings.HasPrefix(name, "base") {
				rule.Path = strings.ReplaceAll(rule.Path, "{{"+"path"+"}}", config.RequestPath)
				rule.Path = strings.ReplaceAll(rule.Path, "{{"+"checkwebshellpath"+"}}", config.CheckWebshellPath)
				rule.Body = strings.ReplaceAll(rule.Body, "{{"+"body"+"}}", config.RequestBody)
				rule.Body = strings.ReplaceAll(rule.Body, "{{"+"writewebshellbody"+"}}", config.WriteWebShellBody)
				rule.Body = strings.ReplaceAll(rule.Body, "{{"+"pocbody"+"}}", config.PocBody)
				rule.Body = strings.ReplaceAll(rule.Body, "{{"+"webshellcommand"+"}}", config.WebShellCommand)
				rule.Expression = strings.ReplaceAll(rule.Expression, "{{"+"checkpocres"+"}}", config.CheckPocResBody)
				rule.Expression = strings.ReplaceAll(rule.Expression, "{{"+"checkexpres"+"}}", config.CheckExpResBody)
				rule.Expression = strings.ReplaceAll(rule.Expression, "{{"+"checkwebshellcmdres"+"}}", config.CheckWebShellCmdBody)
			} else {
				rule.Path = strings.ReplaceAll(rule.Path, "{{"+k1+"}}", value)
				rule.Body = strings.ReplaceAll(rule.Body, "{{"+k1+"}}", value)
			}
		}

		if oReq.URL.Path != "" && oReq.URL.Path != "/" {
			req.Url.Path = fmt.Sprint(oReq.URL.Path, rule.Path)
		} else {
			req.Url.Path = rule.Path
		}
		req.Url.Path = strings.ReplaceAll(req.Url.Path, " ", "%20")

		var newRequest *http.Request
		newRequest, err = http.NewRequest(rule.Method, fmt.Sprintf("%s://%s%s", req.Url.Scheme, req.Url.Host, req.Url.Path), strings.NewReader(rule.Body))
		if err != nil {
			//fmt.Println("[-] newRequest error: ",err)
			return false, err
		}
		newRequest.Header = oReq.Header.Clone()
		for k, v := range Headers {
			newRequest.Header.Set(k, v)
		}
		Headers = nil
		var resp *Response
		resp, err = DoRequest(newRequest, rule.FollowRedirects)
		newRequest = nil
		if err != nil {
			return false, err
		}
		variableMap["response"] = resp

		if rule.Search != "" {
			result := ResSearch(rule.Search, FormatHeader(resp.Headers)+string(resp.Body))
			if len(result) > 0 { // 正则匹配成功
				for k, v := range result {
					variableMap[k] = v
				}
			} else {
				return false, nil
			}
		}
		var out ref.Val
		out, err = Evaluate(env, rule.Expression, variableMap)
		if err != nil {
			return false, err
		}
		//如果false不继续执行后续rule
		// 如果最后一步执行失败，就算前面成功了最终依旧是失败
		flag, ok = out.Value().(bool)
		if !ok {
			flag = false
		}
		return flag, nil
	}

	DealWithRules := func(name string, rules []Rules) bool {
		successFlag := false
		for _, rule := range rules {
			var flag bool
			flag, err = DealWithRule(name, rule)
			if err != nil || !flag { //如果false不继续执行后续rule
				successFlag = false // 如果其中一步为false，则直接break
				break
			}
			successFlag = true
		}
		return successFlag
	}

	if len(p.Rules) > 0 {
		success = DealWithRules(p.Name, p.Rules)
	} else {
		// groups是分组的rules
		for _, item := range p.Groups {
			name, rules := item.Key, item.Value
			success = DealWithRules(p.Name, rules)
			if success {
				return success, nil, name
			}
		}
	}

	return success, nil, ""
}

func NewReverse() *Reverse {
	//if !config.DnsLog {
	//	return &Reverse{}
	//}
	if !config.DnsLog || !strings.HasSuffix(config.CeyeURL, ".ceye.io") {
		return &Reverse{}
	}
	letters := "1234567890abcdefghijklmnopqrstuvwxyz"
	randSourceData := rand.New(rand.NewSource(time.Now().UnixNano()))
	sub := RandomStr(randSourceData, letters, 8)
	// 生成随机指定ceye域名的子域名
	urlStr := fmt.Sprintf("http://%s.%s", sub, config.CeyeURL)
	u, _ := url.Parse(urlStr)
	return &Reverse{
		Url:                urlStr,
		Domain:             u.Hostname(),
		Ip:                 u.Host,
		IsDomainNameServer: false,
	}
}

package poclib

import (
	"PocScan/common"
	"PocScan/config"
	"crypto/md5"
	"fmt"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types/ref"
	"net/http"
	"regexp"
	"strings"
)

func brutePoc(oReq *http.Request, p *Poc, variableMap map[string]interface{}, req *Request, env *cel.Env) (success bool, err error) {
	var strMap StrMap
	var tmpnum int
	// for example
	//sets:
	//path:
	//	- "1.sql"
	//	- "backup.sql"
	//	- "database.sql"
	//	- "data.sql"
	//	- "db_backup.sql"
	//	- "dbdump.sql"
	//	- "db.sql"
	//	- "dump.sql"
	//	- "{{host}}.sql"
	//	- "{{host}}_db.sql"
	//	- "localhost.sql"
	//	- "mysqldump.sql"
	//	- "mysql.sql"
	//	- "site.sql"
	//	- "sql.sql"
	//	- "temp.sql"
	//	- "translate.sql"
	//	- "users.sql"
	// rules:
	//  - method: GET
	//    path: /{{path}}
	//    follow_redirects: false
	//    continue: true
	//    expression: |
	//      "(?m)(?:DROP|CREATE|(?:UN)?LOCK) TABLE|INSERT INTO".bmatches(response.body)
	for i, rule := range p.Rules {
		// 寻找爆破的{{}}标记
		// 非fuzz情况
		if !isFuzz(rule, p.Sets) {
			success, err = bruteSend(oReq, variableMap, req, env, rule)
			if err != nil {
				return false, err
			}
			if success {
				continue
			} else {
				return false, err
			}
		}
		// fuzz情况
		// fuzz字典
		setsMap := makeSets(p.Sets)
		ruleHash := make(map[string]struct{})
		// 循环poc扫描 label语法
	look:
		for j, item := range setsMap {
			//shiro 10 keys default scan
			if p.Name == "poc-yaml-shiro-key" && config.FullPOC == false && j >= 10 {
				// todo: 看看这个cbc要不要单独处理
				if item[1] == "cbc" {
					continue
				} else {
					if tmpnum == 0 {
						tmpnum = j
					}
					if j-tmpnum >= 10 {
						break
					}
				}
			}

			rule1 := rulesClone(rule)

			var flag1 bool
			var tmpMap StrMap
			var payloads = make(map[string]interface{})
			var tmpExpression string

			for n, one := range p.Sets {
				// for example
				//sets:
				//username:
				//	- tomcat
				//	- admin
				//	- root
				//	- manager
				//password:
				//	- tomcat
				//	- ""
				//	- admin
				//	- 123456
				//	- root
				//payload:
				//	- base64(username+":"+password)
				//Username: tomcat, Password: tomcat, Payload: base64(username+":"+password)
				// ...

				// key = username, expression = tomcat
				// key = password, expression = tomcat
				// key = payload, expression = base64(username+":"+password)
				key, expression := one.Key, item[n]
				if key == "payload" {
					tmpExpression = expression
				}
				_, output := EvalSetAnother(env, variableMap, key, expression)
				// payloads['username'] = tomcat
				// payloads['password'] = tomcat
				// payloads['payload'] = base64编码结果
				payloads[key] = output
			}

			for _, one := range p.Sets {
				flag := false
				key := one.Key
				value := fmt.Sprintf("%v", payloads[key])

				for k2, v2 := range rule1.Headers {
					if strings.Contains(v2, "{{"+key+"}}") {
						rule1.Headers[k2] = strings.ReplaceAll(v2, "{{"+key+"}}", value)
						flag = true
					}
				}
				if strings.Contains(rule1.Path, "{{"+key+"}}") {
					rule1.Path = strings.ReplaceAll(rule1.Path, "{{"+key+"}}", value)
					flag = true
				}
				if strings.Contains(rule1.Body, "{{"+key+"}}") {
					rule1.Body = strings.ReplaceAll(rule1.Body, "{{"+key+"}}", value)
					flag = true
				}

				// 有替换标签{{}}的情况
				if flag {
					flag1 = true
					if key == "payload" {
						var flag2 bool
						for k, v := range variableMap {
							// tmpExpression = base64(username+":"+password)
							if strings.Contains(tmpExpression, k) {
								flag2 = true
								// 特殊情况，tmpExpression包含了k
								tmpMap = append(tmpMap, StrItem{k, fmt.Sprintf("%v", v)})
							}
						}
						if flag2 {
							continue
						}
					}
					// 将前面的payloads列表都放到tmpMap里
					tmpMap = append(tmpMap, StrItem{key, value})
				}

			}
			// 无替换{{}}则下一次for
			if !flag1 {
				continue
			}

			has := md5.Sum([]byte(fmt.Sprintf("%v", rule1)))
			md5str := fmt.Sprintf("%x", has)

			// 重复检测
			if _, ok := ruleHash[md5str]; ok {
				continue
			}
			ruleHash[md5str] = struct{}{}

			// 每一次look循环变的是variableMap
			// 每一次look外的for循环变的是rule
			success, err = bruteSend(oReq, variableMap, req, env, rule1)
			if err != nil {
				return false, err
			}

			// 若执行成功则跳出look循环
			if success {
				if rule.Continue {
					if p.Name == "poc-yaml-backup-file" || p.Name == "poc-yaml-sql-file" {
						common.LogSuccess(fmt.Sprintf("[+] PocScan %s://%s%s %s", req.Url.Scheme, req.Url.Host, req.Url.Path, p.Name))
					} else {
						common.LogSuccess(fmt.Sprintf("[+] PocScan %s://%s%s %s %v", req.Url.Scheme, req.Url.Host, req.Url.Path, p.Name, tmpMap))
					}
					// 继续下一次for j, item := range setsMap
					// 变换下一个排列组合：例如账号密码
					continue
				}
				// 尝试成功的组合
				strMap = append(strMap, tmpMap...)
				if i == len(p.Rules)-1 {
					common.LogSuccess(fmt.Sprintf("[+] PocScan %s://%s%s %s %v", req.Url.Scheme, req.Url.Host, req.Url.Path, p.Name, strMap))
					//遍历rules完成
					return true, nil
				}
				break look
			}
		}
		// 若尝试完未成功则跳出for j, item := range setsMap循环
		if !success {
			break
		}
		// 此时位于for i, rule := range p.Rules
		if rule.Continue {
			// 此时必定为执行失败状态，否则之前就执行完了
			return false, nil
		}
	}
	return success, nil
}

func isFuzz(rule Rules, Sets ListMap) bool {
	//sets:
	//path:
	//	- "1.sql"
	//	- "backup.sql"
	//	- "database.sql"
	//	- "data.sql"
	//	- "db_backup.sql"
	//	- "dbdump.sql"
	//	- "db.sql"
	//	- "dump.sql"
	//	- "{{host}}.sql"
	//	- "{{host}}_db.sql"
	//	- "localhost.sql"
	//	- "mysqldump.sql"
	//	- "mysql.sql"
	//	- "site.sql"
	//	- "sql.sql"
	//	- "temp.sql"
	//	- "translate.sql"
	//	- "users.sql"
	// rules:
	//  - method: GET
	//    path: /{{path}}
	//    follow_redirects: false
	//    continue: true
	//    expression: |
	//      "(?m)(?:DROP|CREATE|(?:UN)?LOCK) TABLE|INSERT INTO".bmatches(response.body)

	// headers、path、body三种fuzz，例如根据sets当中的path字段和{{path}}匹配，判断出path处fuzz
	for _, one := range Sets {
		key := one.Key
		for _, v := range rule.Headers {
			if strings.Contains(v, "{{"+key+"}}") {
				return true
			}
		}
		if strings.Contains(rule.Path, "{{"+key+"}}") {
			return true
		}
		if strings.Contains(rule.Body, "{{"+key+"}}") {
			return true
		}
	}
	return false
}

func bruteSend(oReq *http.Request, variableMap map[string]interface{}, req *Request, env *cel.Env, rule Rules) (bool, error) {
	// for example
	//set:
	//  host: request.url.domain
	// 经过前面的处理，variableMap['host'] = request.url.domain
	//rules:
	//  - method: GET
	//    path: /{{path}}
	//    follow_redirects: false
	//    continue: true
	//    expression: |
	//      "(?m)(?:DROP|CREATE|(?:UN)?LOCK) TABLE|INSERT INTO".bmatches(response.body)

	// 遍历variableMap里的set表达式
	// {{set.key}} -> set.value
	for k1, v1 := range variableMap {
		_, isMap := v1.(map[string]string)
		if isMap {
			continue
		}
		value := fmt.Sprintf("%v", v1)
		for k2, v2 := range rule.Headers {
			if strings.Contains(v2, "{{"+k1+"}}") {
				rule.Headers[k2] = strings.ReplaceAll(v2, "{{"+k1+"}}", value)
			}
		}
		rule.Path = strings.ReplaceAll(strings.TrimSpace(rule.Path), "{{"+k1+"}}", value)
		rule.Body = strings.ReplaceAll(strings.TrimSpace(rule.Body), "{{"+k1+"}}", value)
	}

	if oReq.URL.Path != "" && oReq.URL.Path != "/" {
		// 子目录遍历
		req.Url.Path = fmt.Sprint(oReq.URL.Path, rule.Path)
	} else {
		// 一级目录遍历
		req.Url.Path = rule.Path
	}

	// 某些poc没有区分path和query，需要处理
	// url编码%20
	req.Url.Path = strings.ReplaceAll(req.Url.Path, " ", "%20")

	// 创建http请求
	newRequest, err := http.NewRequest(rule.Method, fmt.Sprintf("%s://%s%s", req.Url.Scheme, req.Url.Host, req.Url.Path), strings.NewReader(rule.Body))
	if err != nil {
		return false, err
	}
	newRequest.Header = oReq.Header.Clone()
	for k, v := range rule.Headers {
		newRequest.Header.Set(k, v)
	}
	var resp *Response
	// 发送请求
	resp, err = DoRequest(newRequest, rule.FollowRedirects)
	newRequest = nil
	if err != nil {
		return false, err
	}
	variableMap["response"] = resp

	// 先判断响应页面是否匹配search规则，判断poc是否执行成功
	if rule.Search != "" {
		result := ResSearch(rule.Search, FormatHeader(resp.Headers)+string(resp.Body))
		if result != nil && len(result) > 0 { // 正则匹配成功
			for k, v := range result {
				// rule.search的key（匹配变量） = body当中匹配到的结果
				variableMap[k] = v
			}
			//return false, nil
		} else {
			return false, nil
		}
	}

	// 执行rule词条的expression，二次匹配
	// for example
	//expression: |
	//	response.status == 200
	var out ref.Val
	out, err = Evaluate(env, rule.Expression, variableMap)
	if err != nil {
		if strings.Contains(err.Error(), "Syntax error") {
			fmt.Println(rule.Expression, err)
		}
		return false, err
	}

	if fmt.Sprintf("%v", out) == "false" { //如果false不继续执行后续rule
		return false, err // 如果最后一步执行失败，就算前面成功了最终依旧是失败
	}
	return true, err
}

// poc yaml rules字段处理
func rulesClone(tags Rules) Rules {
	cloneTags := Rules{}
	cloneTags.Method = tags.Method
	cloneTags.Path = tags.Path
	cloneTags.Body = tags.Body
	cloneTags.Search = tags.Search
	cloneTags.FollowRedirects = tags.FollowRedirects
	cloneTags.Expression = tags.Expression
	cloneTags.Headers = MapClone(tags.Headers)
	return cloneTags
}

func MapClone(tags map[string]string) map[string]string {
	cloneTags := make(map[string]string)
	for k, v := range tags {
		cloneTags[k] = v
	}
	return cloneTags
}
func makeSets(input ListMap) (output [][]string) {
	// 多组字典
	if len(input) > 1 {
		// 递归
		output = makeSets(input[1:])
		// 排列组合
		output = MakeData(output, input[0].Value)
	} else {
		// 一组字典
		for _, i := range input[0].Value {
			output = append(output, []string{i})
		}
	}
	return
}

func MakeData(base [][]string, nextData []string) (output [][]string) {
	for i := range base {
		for _, j := range nextData {
			output = append(output, append([]string{j}, base[i]...))
		}
	}
	return
}
func ResSearch(re string, body string) map[string]string {
	// rule.Search字段匹配
	r, err := regexp.Compile(re)
	if err != nil {
		fmt.Println("[-] regexp.Compile error: ", err)
		return nil
	}
	result := r.FindStringSubmatch(body)
	//用于获取正则表达式中命名捕获组的名称
	names := r.SubexpNames()

	if len(result) > 1 && len(names) > 1 {
		paramsMap := make(map[string]string)
		for i, name := range names {
			if i > 0 && i <= len(result) {
				// 若匹配规则有Set-Cookie:
				if strings.HasPrefix(re, "Set-Cookie:") && strings.Contains(name, "cookie") {
					paramsMap[name] = formatCookies(result[i])
				} else {
					paramsMap[name] = result[i]
				}
			}
		}
		return paramsMap
	}
	return nil
}

func FormatHeader(header map[string]string) (output string) {
	for name, values := range header {
		line := fmt.Sprintf("%s: %s\n", name, values)
		output = output + line
	}
	output = output + "\r\n"
	return
}

func formatCookies(rawCookie string) (output string) {
	// format the cookies
	parsedCookie := strings.Split(rawCookie, "; ")
	for _, c := range parsedCookie {
		nameVal := strings.Split(c, "=")
		if len(nameVal) >= 2 {
			switch strings.ToLower(nameVal[0]) {
			case "expires", "max-age", "path", "domain", "version", "comment", "secure", "samesite", "httponly":
				continue
			}
			output += fmt.Sprintf("%s=%s; ", nameVal[0], strings.Join(nameVal[1:], "="))
		}
	}

	return
}

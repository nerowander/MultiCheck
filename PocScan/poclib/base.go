package poclib

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter/functions"
	"github.com/nerowander/MultiCheck/WebScan/lib"
	"github.com/nerowander/MultiCheck/common"
	"github.com/nerowander/MultiCheck/config"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type CustomLib struct {
	envOptions     []cel.EnvOption
	programOptions []cel.ProgramOption
}

func (c *CustomLib) CompileOptions() []cel.EnvOption {
	return c.envOptions
}

func (c *CustomLib) ProgramOptions() []cel.ProgramOption {
	return c.programOptions
}

func NewEnv(c *CustomLib) (*cel.Env, error) {
	return cel.NewEnv(cel.Lib(c))
}
func NewEnvOption() CustomLib {
	c := CustomLib{}

	c.envOptions = []cel.EnvOption{
		cel.Container("lib"),
		cel.Types(
			&UrlType{},
			&Request{},
			&Response{},
			&Reverse{},
		),
		cel.Declarations(
			decls.NewIdent("request", decls.NewObjectType("lib.Request"), nil),
			decls.NewIdent("response", decls.NewObjectType("lib.Response"), nil),
			decls.NewIdent("reverse", decls.NewObjectType("lib.Reverse"), nil),
		),
		cel.Declarations(
			// functions
			decls.NewFunction("bcontains",
				decls.NewInstanceOverload("bytes_bcontains_bytes",
					[]*exprpb.Type{decls.Bytes, decls.Bytes},
					decls.Bool)),
			decls.NewFunction("bmatches",
				decls.NewInstanceOverload("string_bmatches_bytes",
					[]*exprpb.Type{decls.String, decls.Bytes},
					decls.Bool)),
			decls.NewFunction("md5",
				decls.NewOverload("md5_string",
					[]*exprpb.Type{decls.String},
					decls.String)),
			decls.NewFunction("randomInt",
				decls.NewOverload("randomInt_int_int",
					[]*exprpb.Type{decls.Int, decls.Int},
					decls.Int)),
			decls.NewFunction("randomLowercase",
				decls.NewOverload("randomLowercase_int",
					[]*exprpb.Type{decls.Int},
					decls.String)),
			decls.NewFunction("randomUppercase",
				decls.NewOverload("randomUppercase_int",
					[]*exprpb.Type{decls.Int},
					decls.String)),
			decls.NewFunction("randomString",
				decls.NewOverload("randomString_int",
					[]*exprpb.Type{decls.Int},
					decls.String)),
			decls.NewFunction("base64",
				decls.NewOverload("base64_string",
					[]*exprpb.Type{decls.String},
					decls.String)),
			decls.NewFunction("base64",
				decls.NewOverload("base64_bytes",
					[]*exprpb.Type{decls.Bytes},
					decls.String)),
			decls.NewFunction("base64Decode",
				decls.NewOverload("base64Decode_string",
					[]*exprpb.Type{decls.String},
					decls.String)),
			decls.NewFunction("base64Decode",
				decls.NewOverload("base64Decode_bytes",
					[]*exprpb.Type{decls.Bytes},
					decls.String)),
			decls.NewFunction("urlencode",
				decls.NewOverload("urlencode_string",
					[]*exprpb.Type{decls.String},
					decls.String)),
			decls.NewFunction("urlencode",
				decls.NewOverload("urlencode_bytes",
					[]*exprpb.Type{decls.Bytes},
					decls.String)),
			decls.NewFunction("urldecode",
				decls.NewOverload("urldecode_string",
					[]*exprpb.Type{decls.String},
					decls.String)),
			decls.NewFunction("urldecode",
				decls.NewOverload("urldecode_bytes",
					[]*exprpb.Type{decls.Bytes},
					decls.String)),
			decls.NewFunction("substr",
				decls.NewOverload("substr_string_int_int",
					[]*exprpb.Type{decls.String, decls.Int, decls.Int},
					decls.String)),
			decls.NewFunction("wait",
				decls.NewInstanceOverload("reverse_wait_int",
					[]*exprpb.Type{decls.Any, decls.Int},
					decls.Bool)),
			decls.NewFunction("icontains",
				decls.NewInstanceOverload("icontains_string",
					[]*exprpb.Type{decls.String, decls.String},
					decls.Bool)),
			decls.NewFunction("TDdate",
				decls.NewOverload("tongda_date",
					[]*exprpb.Type{},
					decls.String)),
			decls.NewFunction("shirokey",
				decls.NewOverload("shiro_key",
					[]*exprpb.Type{decls.String, decls.String},
					decls.String)),
			decls.NewFunction("startsWith",
				decls.NewInstanceOverload("startsWith_bytes",
					[]*exprpb.Type{decls.Bytes, decls.Bytes},
					decls.Bool)),
			decls.NewFunction("istartsWith",
				decls.NewInstanceOverload("startsWith_string",
					[]*exprpb.Type{decls.String, decls.String},
					decls.Bool)),
			decls.NewFunction("hexdecode",
				decls.NewInstanceOverload("hexdecode",
					[]*exprpb.Type{decls.String},
					decls.Bytes)),
		),
	}
	c.programOptions = []cel.ProgramOption{
		cel.Functions(
			&functions.Overload{
				Operator: "bytes_bcontains_bytes",
				Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
					v1, ok := lhs.(types.Bytes)
					if !ok {
						return types.ValOrErr(lhs, "unexpected type '%v' passed to bcontains", lhs.Type())
					}
					v2, ok := rhs.(types.Bytes)
					if !ok {
						return types.ValOrErr(rhs, "unexpected type '%v' passed to bcontains", rhs.Type())
					}
					return types.Bool(bytes.Contains(v1, v2))
				},
			},
			&functions.Overload{
				Operator: "string_bmatches_bytes",
				Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
					v1, ok := lhs.(types.String)
					if !ok {
						return types.ValOrErr(lhs, "unexpected type '%v' passed to bmatch", lhs.Type())
					}
					v2, ok := rhs.(types.Bytes)
					if !ok {
						return types.ValOrErr(rhs, "unexpected type '%v' passed to bmatch", rhs.Type())
					}
					ok, err := regexp.Match(string(v1), v2)
					if err != nil {
						return types.NewErr("%v", err)
					}
					return types.Bool(ok)
				},
			},
			&functions.Overload{
				Operator: "md5_string",
				Unary: func(value ref.Val) ref.Val {
					v, ok := value.(types.String)
					if !ok {
						return types.ValOrErr(value, "unexpected type '%v' passed to md5_string", value.Type())
					}
					return types.String(fmt.Sprintf("%x", md5.Sum([]byte(v))))
				},
			},
			&functions.Overload{
				Operator: "randomInt_int_int",
				Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
					from, ok := lhs.(types.Int)
					if !ok {
						return types.ValOrErr(lhs, "unexpected type '%v' passed to randomInt", lhs.Type())
					}
					to, ok := rhs.(types.Int)
					if !ok {
						return types.ValOrErr(rhs, "unexpected type '%v' passed to randomInt", rhs.Type())
					}
					min, max := int(from), int(to)
					return types.Int(rand.Intn(max-min) + min)
				},
			},
			&functions.Overload{
				Operator: "randomLowercase_int",
				Unary: func(value ref.Val) ref.Val {
					n, ok := value.(types.Int)
					if !ok {
						return types.ValOrErr(value, "unexpected type '%v' passed to randomLowercase", value.Type())
					}
					return types.String(randomLowercase(int(n)))
				},
			},
			&functions.Overload{
				Operator: "randomUppercase_int",
				Unary: func(value ref.Val) ref.Val {
					n, ok := value.(types.Int)
					if !ok {
						return types.ValOrErr(value, "unexpected type '%v' passed to randomUppercase", value.Type())
					}
					return types.String(randomUppercase(int(n)))
				},
			},
			&functions.Overload{
				Operator: "randomString_int",
				Unary: func(value ref.Val) ref.Val {
					n, ok := value.(types.Int)
					if !ok {
						return types.ValOrErr(value, "unexpected type '%v' passed to randomString", value.Type())
					}
					return types.String(randomString(int(n)))
				},
			},
			&functions.Overload{
				Operator: "base64_string",
				Unary: func(value ref.Val) ref.Val {
					v, ok := value.(types.String)
					if !ok {
						return types.ValOrErr(value, "unexpected type '%v' passed to base64_string", value.Type())
					}
					return types.String(base64.StdEncoding.EncodeToString([]byte(v)))
				},
			},
			&functions.Overload{
				Operator: "base64_bytes",
				Unary: func(value ref.Val) ref.Val {
					v, ok := value.(types.Bytes)
					if !ok {
						return types.ValOrErr(value, "unexpected type '%v' passed to base64_bytes", value.Type())
					}
					return types.String(base64.StdEncoding.EncodeToString(v))
				},
			},
			&functions.Overload{
				Operator: "base64Decode_string",
				Unary: func(value ref.Val) ref.Val {
					v, ok := value.(types.String)
					if !ok {
						return types.ValOrErr(value, "unexpected type '%v' passed to base64Decode_string", value.Type())
					}
					decodeBytes, err := base64.StdEncoding.DecodeString(string(v))
					if err != nil {
						return types.NewErr("%v", err)
					}
					return types.String(decodeBytes)
				},
			},
			&functions.Overload{
				Operator: "base64Decode_bytes",
				Unary: func(value ref.Val) ref.Val {
					v, ok := value.(types.Bytes)
					if !ok {
						return types.ValOrErr(value, "unexpected type '%v' passed to base64Decode_bytes", value.Type())
					}
					decodeBytes, err := base64.StdEncoding.DecodeString(string(v))
					if err != nil {
						return types.NewErr("%v", err)
					}
					return types.String(decodeBytes)
				},
			},
			&functions.Overload{
				Operator: "urlencode_string",
				Unary: func(value ref.Val) ref.Val {
					v, ok := value.(types.String)
					if !ok {
						return types.ValOrErr(value, "unexpected type '%v' passed to urlencode_string", value.Type())
					}
					return types.String(url.QueryEscape(string(v)))
				},
			},
			&functions.Overload{
				Operator: "urlencode_bytes",
				Unary: func(value ref.Val) ref.Val {
					v, ok := value.(types.Bytes)
					if !ok {
						return types.ValOrErr(value, "unexpected type '%v' passed to urlencode_bytes", value.Type())
					}
					return types.String(url.QueryEscape(string(v)))
				},
			},
			&functions.Overload{
				Operator: "urldecode_string",
				Unary: func(value ref.Val) ref.Val {
					v, ok := value.(types.String)
					if !ok {
						return types.ValOrErr(value, "unexpected type '%v' passed to urldecode_string", value.Type())
					}
					decodeString, err := url.QueryUnescape(string(v))
					if err != nil {
						return types.NewErr("%v", err)
					}
					return types.String(decodeString)
				},
			},
			&functions.Overload{
				Operator: "urldecode_bytes",
				Unary: func(value ref.Val) ref.Val {
					v, ok := value.(types.Bytes)
					if !ok {
						return types.ValOrErr(value, "unexpected type '%v' passed to urldecode_bytes", value.Type())
					}
					decodeString, err := url.QueryUnescape(string(v))
					if err != nil {
						return types.NewErr("%v", err)
					}
					return types.String(decodeString)
				},
			},
			&functions.Overload{
				Operator: "substr_string_int_int",
				Function: func(values ...ref.Val) ref.Val {
					if len(values) == 3 {
						str, ok := values[0].(types.String)
						if !ok {
							return types.NewErr("invalid string to 'substr'")
						}
						start, ok := values[1].(types.Int)
						if !ok {
							return types.NewErr("invalid start to 'substr'")
						}
						length, ok := values[2].(types.Int)
						if !ok {
							return types.NewErr("invalid length to 'substr'")
						}
						runes := []rune(str)
						if start < 0 || length < 0 || int(start+length) > len(runes) {
							return types.NewErr("invalid start or length to 'substr'")
						}
						return types.String(runes[start : start+length])
					} else {
						return types.NewErr("too many arguments to 'substr'")
					}
				},
			},
			&functions.Overload{
				Operator: "reverse_wait_int",
				Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
					reverse, ok := lhs.Value().(*Reverse)
					if !ok {
						return types.ValOrErr(lhs, "unexpected type '%v' passed to 'wait'", lhs.Type())
					}
					timeout, ok := rhs.Value().(int64)
					if !ok {
						return types.ValOrErr(rhs, "unexpected type '%v' passed to 'wait'", rhs.Type())
					}
					// dnslog
					return types.Bool(reverseCheck(reverse, timeout))
				},
			},

			// dnslog的平台地址可由用户自己提供
			//&functions.Overload{
			//	Operator: "reverse_wait_int",
			//	Function: func(strVal ...ref.Val) ref.Val {
			//		if len(strVal) == 3 {
			//			reverse, ok := strVal[0].Value().(*Reverse)
			//			if !ok {
			//				return types.ValOrErr(strVal[0], "unexpected type '%v' passed to 'wait'", lhs.Type())
			//			}
			//			timeout, ok := strVal[1].Value().(int64)
			//			if !ok {
			//				return types.ValOrErr(strVal[1], "unexpected type '%v' passed to 'wait'", rhs.Type())
			//			}
			//			strParam, ok := strVal[2].Value().(string)
			//			// dnslog
			//			return types.Bool(reverseCheck(reverse, timeout))
			//		} else {
			//			return types.NewErr("too many arguments to 'reverse_wait_int'")
			//		}
			//	},
			//},
			&functions.Overload{
				Operator: "icontains_string",
				Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
					v1, ok := lhs.(types.String)
					if !ok {
						return types.ValOrErr(lhs, "unexpected type '%v' passed to bcontains", lhs.Type())
					}
					v2, ok := rhs.(types.String)
					if !ok {
						return types.ValOrErr(rhs, "unexpected type '%v' passed to bcontains", rhs.Type())
					}
					// 不区分大小写包含
					return types.Bool(strings.Contains(strings.ToLower(string(v1)), strings.ToLower(string(v2))))
				},
			},
			&functions.Overload{
				Operator: "tongda_date",
				Function: func(value ...ref.Val) ref.Val {
					return types.String(time.Now().Format("0601"))
				},
			},
			&functions.Overload{
				Operator: "shiro_key",
				Binary: func(key ref.Val, mode ref.Val) ref.Val {
					v1, ok := key.(types.String)
					if !ok {
						return types.ValOrErr(key, "unexpected type '%v' passed to shiro_key", key.Type())
					}
					v2, ok := mode.(types.String)
					if !ok {
						return types.ValOrErr(mode, "unexpected type '%v' passed to shiro_mode", mode.Type())
					}
					cookie := GetShiroCookie(string(v1), string(v2))
					if cookie == "" {
						return types.NewErr("%v", "key b64decode failed")
					}
					return types.String(cookie)
				},
			},
			&functions.Overload{
				Operator: "startsWith_bytes",
				Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
					v1, ok := lhs.(types.Bytes)
					if !ok {
						return types.ValOrErr(lhs, "unexpected type '%v' passed to startsWith_bytes", lhs.Type())
					}
					v2, ok := rhs.(types.Bytes)
					if !ok {
						return types.ValOrErr(rhs, "unexpected type '%v' passed to startsWith_bytes", rhs.Type())
					}
					// 不区分大小写包含
					return types.Bool(bytes.HasPrefix(v1, v2))
				},
			},
			&functions.Overload{
				Operator: "startsWith_string",
				Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
					v1, ok := lhs.(types.String)
					if !ok {
						return types.ValOrErr(lhs, "unexpected type '%v' passed to startsWith_string", lhs.Type())
					}
					v2, ok := rhs.(types.String)
					if !ok {
						return types.ValOrErr(rhs, "unexpected type '%v' passed to startsWith_string", rhs.Type())
					}
					// 不区分大小写包含
					return types.Bool(strings.HasPrefix(strings.ToLower(string(v1)), strings.ToLower(string(v2))))
				},
			},
			&functions.Overload{
				Operator: "hexdecode",
				Unary: func(lhs ref.Val) ref.Val {
					v1, ok := lhs.(types.String)
					if !ok {
						return types.ValOrErr(lhs, "unexpected type '%v' passed to hexdecode", lhs.Type())
					}
					out, err := hex.DecodeString(string(v1))
					if err != nil {
						return types.ValOrErr(lhs, "hexdecode error: %v", err)
					}
					// 不区分大小写包含
					return types.Bytes(out)
				},
			},
		),
	}
	return c
}

func (c *CustomLib) UpdateCompileOptions(args StrMap) {
	for _, item := range args {
		k, v := item.Key, item.Value
		var d *exprpb.Decl
		if strings.HasPrefix(v, "randomInt") {
			d = decls.NewIdent(k, decls.Int, nil)
		} else if strings.HasPrefix(v, "newReverse") {
			d = decls.NewIdent(k, decls.NewObjectType("lib.Reverse"), nil)
		} else {
			d = decls.NewIdent(k, decls.String, nil)
		}
		c.envOptions = append(c.envOptions, cel.Declarations(d))
	}
}

//func (c *CustomLib) UpdateCompileExpOptions(args explib.StrMap) {
//	for _, item := range args {
//		k, v := item.Key, item.Value
//		var d *exprpb.Decl
//		if strings.HasPrefix(v, "randomInt") {
//			d = decls.NewIdent(k, decls.Int, nil)
//		} else if strings.HasPrefix(v, "newReverse") {
//			d = decls.NewIdent(k, decls.NewObjectType("lib.Reverse"), nil)
//		} else {
//			d = decls.NewIdent(k, decls.String, nil)
//		}
//		c.envOptions = append(c.envOptions, cel.Declarations(d))
//	}
//}

var randSource = rand.New(rand.NewSource(time.Now().Unix()))

func RandomStr(randSource *rand.Rand, letterBytes string, n int) string {
	const (
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)
	randBytes := make([]byte, n)
	for i, cache, remain := n-1, randSource.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = randSource.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			randBytes[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(randBytes)
}
func randomLowercase(n int) string {
	lowercase := "abcdefghijklmnopqrstuvwxyz"
	return RandomStr(randSource, lowercase, n)
}

func randomUppercase(n int) string {
	uppercase := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	return RandomStr(randSource, uppercase, n)
}

func randomString(n int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	return RandomStr(randSource, charset, n)
}

// 默认使用ceye作为dnslog平台
func reverseCheck(r *Reverse, timeout int64) bool {
	if r.Domain == "" || !config.DnsLog {
		return false
	}
	if config.CeyeToken == "" {
		return false
	}
	time.Sleep(time.Second * time.Duration(timeout))
	sub := strings.Split(r.Domain, ".")[0]
	urlStr := fmt.Sprintf("http://api.ceye.io/v1/records?token=%s&type=dns&filter=%s", config.CeyeToken, sub)
	// 可能后续得加上token等数据
	//fmt.Println(urlStr)
	req, _ := http.NewRequest("GET", urlStr, nil)
	resp, err := DoRequest(req, false)
	if err != nil {
		return false
	}

	if !bytes.Contains(resp.Body, []byte(`"data": []`)) && bytes.Contains(resp.Body, []byte(`"message": "OK"`)) { // 返回结果不为空
		fmt.Println(urlStr)
		return true
	}
	return false
}

func DoRequest(req *http.Request, redirect bool) (*Response, error) {
	if req.Body == nil || req.Body == http.NoBody {
	} else {
		req.Header.Set("Content-Length", strconv.Itoa(int(req.ContentLength)))
		if req.Header.Get("Content-Type") == "" {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
	}
	var oResp *http.Response
	var err error
	if redirect {
		oResp, err = lib.Client.Do(req)
	} else {
		oResp, err = lib.ClientNoRedirect.Do(req)
	}
	if err != nil {
		//fmt.Println("[-]DoRequest error: ",err)
		return nil, err
	}
	defer oResp.Body.Close()
	var resp *Response
	resp, err = ParseResponse(oResp)
	if err != nil {
		common.LogError("[-] ParseResponse error: " + err.Error())
		//return nil, err
	}
	return resp, err
}

func ParseUrl(u *url.URL) *UrlType {
	nu := &UrlType{}
	nu.Scheme = u.Scheme
	nu.Domain = u.Hostname()
	nu.Host = u.Host
	nu.Port = u.Port()
	nu.Path = u.EscapedPath()
	nu.Query = u.RawQuery
	nu.Fragment = u.Fragment
	return nu
}

func ParseRequest(oReq *http.Request) (*Request, error) {
	req := &Request{}
	req.Method = oReq.Method
	req.Url = ParseUrl(oReq.URL)
	header := make(map[string]string)
	for k := range oReq.Header {
		header[k] = oReq.Header.Get(k)
	}
	req.Headers = header
	req.ContentType = oReq.Header.Get("Content-Type")
	if oReq.Body == nil || oReq.Body == http.NoBody {
	} else {
		data, err := io.ReadAll(oReq.Body)
		if err != nil {
			return nil, err
		}
		req.Body = data
		oReq.Body = io.NopCloser(bytes.NewBuffer(data))
	}
	return req, nil
}
func ParseResponse(oResp *http.Response) (*Response, error) {
	var resp Response
	header := make(map[string]string)
	resp.Status = int32(oResp.StatusCode)
	resp.Url = ParseUrl(oResp.Request.URL)
	for k := range oResp.Header {
		header[k] = strings.Join(oResp.Header.Values(k), ";")
	}
	resp.Headers = header
	resp.ContentType = oResp.Header.Get("Content-Type")
	body, _ := getRespBody(oResp)
	resp.Body = body
	return &resp, nil
}

func getRespBody(oResp *http.Response) (body []byte, err error) {
	body, err = io.ReadAll(oResp.Body)
	if strings.Contains(oResp.Header.Get("Content-Encoding"), "gzip") {
		reader, err1 := gzip.NewReader(bytes.NewReader(body))
		if err1 == nil {
			body, err = io.ReadAll(reader)
		}
	}
	if err == io.EOF {
		err = nil
	}
	return body, err
}

func EvalSet(env *cel.Env, variableMap map[string]interface{}, k string, expression string) (err error, output string) {
	var out ref.Val
	// for example
	// expression = newReverse()
	// params = map[string]interface{}{"reverse": newReverse()} 即为variablemap
	out, err = Evaluate(env, expression, variableMap)
	if err != nil {
		variableMap[k] = expression
	} else {
		switch value := out.Value().(type) {
		case *UrlType:
			variableMap[k] = UrlTypeToString(value)
		case int64:
			variableMap[k] = int(value)
		default:
			variableMap[k] = fmt.Sprintf("%v", out)
		}
	}
	return err, fmt.Sprintf("%v", variableMap[k])
}

func EvalSetAnother(env *cel.Env, variableMap map[string]interface{}, k string, expression string) (err error, output string) {
	var out ref.Val
	out, err = Evaluate(env, expression, variableMap)
	if err != nil {
		variableMap[k] = expression
	} else {
		variableMap[k] = fmt.Sprintf("%v", out)
	}
	return err, fmt.Sprintf("%v", variableMap[k])
}

func Evaluate(env *cel.Env, expression string, params map[string]interface{}) (ref.Val, error) {
	if expression == "" {
		return types.Bool(true), nil
	}
	//解析表达式
	// for example
	// expression = newReverse()
	// params = map[string]interface{}{"reverse": newReverse()} 即为variablemap
	expr, issues := env.Compile(expression)
	if issues.Err() != nil {
		return nil, issues.Err()
	}
	var prg cel.Program
	var err error
	// 编译表达式
	prg, err = env.Program(expr)
	if err != nil {
		return nil, err
	}
	var out ref.Val
	// 执行表达式
	out, _, err = prg.Eval(params)
	// prg为函数，params为参数
	if err != nil {
		//fmt.Printf("Evaluation error: %v", err)
		return nil, err
	}
	return out, nil
}

func UrlTypeToString(u *UrlType) string {
	var buf strings.Builder
	if u.Scheme != "" {
		buf.WriteString(u.Scheme)
		buf.WriteByte(':')
	}
	if u.Scheme != "" || u.Host != "" {
		if u.Host != "" || u.Path != "" {
			buf.WriteString("//")
		}
		if h := u.Host; h != "" {
			buf.WriteString(u.Host)
		}
	}
	path := u.Path
	if path != "" && path[0] != '/' && u.Host != "" {
		buf.WriteByte('/')
	}
	if buf.Len() == 0 {
		if i := strings.IndexByte(path, ':'); i > -1 && strings.IndexByte(path[:i], '/') == -1 {
			buf.WriteString("./")
		}
	}
	buf.WriteString(path)

	if u.Query != "" {
		buf.WriteByte('?')
		buf.WriteString(u.Query)
	}
	if u.Fragment != "" {
		buf.WriteByte('#')
		buf.WriteString(u.Fragment)
	}
	return buf.String()
}

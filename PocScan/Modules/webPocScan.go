package Modules

import (
	"embed"
	"fmt"
	"github.com/nerowander/MultiCheck/PocScan/poclib"
	"github.com/nerowander/MultiCheck/WebScan/rules"
	"github.com/nerowander/MultiCheck/common"
	"github.com/nerowander/MultiCheck/config"
	"gopkg.in/yaml.v2"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	once sync.Once
	//go:embed pocs
	Pocs    embed.FS
	AllPocs []*poclib.Poc
)

func WebPocScan(info *config.InfoScan) {
	once.Do(initPoc)
	var pocinfo = config.PocInfo
	buf := strings.Split(info.Url, "/")
	pocinfo.Target = strings.Join(buf[:3], "/")

	// 自己指定poc，例如只指定sql注入的检测
	// 1.基础Web漏洞 base
	// 2.各种Web软件漏洞：例如oa、cms software
	// 3.iot 模块
	// 4.全部都扫 all
	// 在获取poc的时候分类不同情况，base poc｜ software poc｜ iot poc｜ all
	// poctype pocname
	if config.PocType == "base" {
		// 指定其中一种poc：例如sql，若为空则后续execute返回所有base poc
		if config.PocName == "" {
			executePoc(pocinfo)
		} else {
			pocinfo.PocName = config.PocName
			executePoc(pocinfo)
		}
	} else if config.PocType == "software" {
		for _, infoStr := range info.WebInfo {
			// Web软件的poc由webinfo决定
			// 若software的pocname匹配不到，则后续execute会返回所有software poc
			pocinfo.PocName = checkInfoPoc(infoStr)
			//fmt.Println(pocinfo)
			executePoc(pocinfo)
		}
	} else if config.PocType == "iot" {
		// todo：补充iot模块
		// 应该是和web base是类似的
		if config.PocName == "" {
			executePoc(pocinfo)
		} else {
			pocinfo.PocName = config.PocName
			executePoc(pocinfo)
		}
	} else if config.PocType == "all" {
		// web+iot
		if config.PocName == "" {
			// web base and iot
			executePoc(pocinfo)
			for _, infoStr := range info.WebInfo {
				// web software
				pocinfo.PocName = checkInfoPoc(infoStr)
				executePoc(pocinfo)
			}
		} else {
			// web base and iot
			pocinfo.PocName = config.PocName
			executePoc(pocinfo)
			for _, infoStr := range info.WebInfo {
				// web software
				pocinfo.PocName = checkInfoPoc(infoStr)
				executePoc(pocinfo)
			}
		}
	} else {
		fmt.Println("[-] invalid PocType: " + config.PocType)
		return
	}

}

// software poc检查
func checkInfoPoc(infostr string) string {
	for _, poc := range rules.PocDatas {
		if strings.Contains(infostr, poc.Name) {
			return poc.Alias
		}
	}
	return ""
}
func executePoc(PocInfo config.Pocinfo) {
	req, err := http.NewRequest("GET", PocInfo.Target, nil)
	if err != nil {
		errlog := fmt.Sprintf("[-] webpocinit %v %v", PocInfo.Target, err)
		common.LogError(errlog)
		return
	}
	req.Header.Set("User-agent", config.UserAgent)
	req.Header.Set("Accept", config.Accept)
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	if config.Cookie != "" {
		req.Header.Set("Cookie", config.Cookie)
	}
	pocs := filterPoc(PocInfo.PocName)

	poclib.CheckMultiPoc(req, pocs, config.PocNum)
}
func initPoc() {
	// 内置poc
	var entries []fs.DirEntry
	var err error
	if config.PocPath == "" {
		if config.PocType == "base" {
			entries, err = Pocs.ReadDir("pocs/base")
		} else if config.PocType == "software" {
			entries, err = Pocs.ReadDir("pocs/software")
		} else if config.PocType == "iot" {
			entries, err = Pocs.ReadDir("pocs/iot")
		} else if config.PocType == "all" {
			entries, err = Pocs.ReadDir("pocs/all")
		} else {
			fmt.Println("[-] invalid PocType: " + config.PocType)
			return
		}
		if err != nil {
			fmt.Printf("[-] init poc error: %v", err)
			return
		}
		for _, one := range entries {
			path := one.Name()
			// 解析yaml或yml
			if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
				if poc, _ := loadPoc(path, Pocs); poc != nil {
					AllPocs = append(AllPocs, poc)
				}
			}
		}
	} else {
		// 自定义poc
		fmt.Println("[+] load poc from " + config.PocPath)
		err = filepath.Walk(config.PocPath,
			func(path string, info os.FileInfo, err error) error {
				if err != nil || info == nil {
					return err
				}
				if !info.IsDir() {
					if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
						poc, _ := loadPocByPath(path)
						if poc != nil {
							AllPocs = append(AllPocs, poc)
						}
					}
				}
				return nil
			})
		if err != nil {
			fmt.Printf("[-] init poc error: %v", err)
		}
	}
}

func loadPoc(fileName string, Pocs embed.FS) (*poclib.Poc, error) {
	p := &poclib.Poc{}
	var path string
	switch config.PocType {
	case "base":
		path = "pocs/base/" + fileName
	case "software":
		path = "pocs/software/" + fileName
	case "iot":
		path = "pocs/iot/" + fileName
	case "all":
		path = "pocs/all/" + fileName
	default:
		return nil, fmt.Errorf("invalid PocType: %s", config.PocType)
	}
	yamlFile, err := Pocs.ReadFile(path)
	if err != nil {
		fmt.Printf("[-] load poc %s error1: %v\n", fileName, err)
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, p)
	if err != nil {
		fmt.Printf("[-] load poc %s error2: %v\n", fileName, err)
		return nil, err
	}
	return p, err
}

// 自定义poc路径
func loadPocByPath(fileName string) (*poclib.Poc, error) {
	p := &poclib.Poc{}
	data, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Printf("[-] load poc %s error3: %v\n", fileName, err)
		return nil, err
	}
	err = yaml.Unmarshal(data, p)
	if err != nil {
		fmt.Printf("[-] load poc %s error4: %v\n", fileName, err)
		return nil, err
	}
	return p, err
}

func filterPoc(pocName string) (pocs []*poclib.Poc) {
	if pocName == "" {
		return AllPocs
	}
	for _, poc := range AllPocs {
		if strings.Contains(poc.Name, pocName) {
			pocs = append(pocs, poc)
		}
	}
	return
}

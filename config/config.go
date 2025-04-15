package config

import "sync"

var DefaultPorts = "21,22,80,81,135,139,443,445,1433,1521,3306,5432,6379,7001,8000,8080,8089,9000,9200,11211,27017"
var Webports = "80,81,82,83,84,85,86,87,88,89,90,91,92,98,99,443,800,801,808,880,888,889,1000,1010,1080,1081,1082,1099,1118,1888,2008,2020,2100,2375,2379,3000,3008,3128,3505,5555,6080,6648,6868,7000,7001,7002,7003,7004,7005,7007,7008,7070,7071,7074,7078,7080,7088,7200,7680,7687,7688,7777,7890,8000,8001,8002,8003,8004,8006,8008,8009,8010,8011,8012,8016,8018,8020,8028,8030,8038,8042,8044,8046,8048,8053,8060,8069,8070,8080,8081,8082,8083,8084,8085,8086,8087,8088,8089,8090,8091,8092,8093,8094,8095,8096,8097,8098,8099,8100,8101,8108,8118,8161,8172,8180,8181,8200,8222,8244,8258,8280,8288,8300,8360,8443,8448,8484,8800,8834,8838,8848,8858,8868,8879,8880,8881,8888,8899,8983,8989,9000,9001,9002,9008,9010,9043,9060,9080,9081,9082,9083,9084,9085,9086,9087,9088,9089,9090,9091,9092,9093,9094,9095,9096,9097,9098,9099,9100,9200,9443,9448,9800,9981,9986,9988,9998,9999,10000,10001,10002,10004,10008,10010,10250,12018,12443,14000,16080,18000,18001,18002,18004,18008,18080,18082,18088,18090,18098,19001,20000,20720,21000,21501,21502,28018,20880"

// todo: 写一下一些常见等url路径，例如配置文件等
var DefaultDirs = []string{}
var DefaultOutputFile = "output.txt"
var WG sync.WaitGroup
var ResLogs = []string{}

//var SaveResult = true

var PortInt = map[string]int{
	"ftp":     21,
	"ssh":     22,
	"smb":     445,
	"mssql":   1433,
	"oracle":  1521,
	"docker":  2375,
	"mysql":   3306,
	"rdp":     3389,
	"psql":    5432,
	"redis":   6379,
	"fastcgi": 9000,
	"mem":     11211,
	"mongodb": 27017,
	"web":     1000000,
	"pocscan": 1000000,
	"exploit": 1000000,
}

var PortMap = map[string]string{
	"ftp":     "21",
	"ssh":     "22",
	"http":    "80",
	"https":   "443",
	"smb":     "445",
	"mssql":   "1433",
	"oracle":  "1521",
	"docker":  "2375",
	"mysql":   "3306",
	"rdp":     "3389",
	"psql":    "5432",
	"redis":   "6379",
	"fastcgi": "9000",
	"mem":     "11211",
	"mongodb": "27017",
	"service": "21,22,135,139,445,1433,1521,2375,3306,3389,5432,6379,9000,11211,27017",
	"all":     "1-65535",
	"db":      "1433,1521,3306,5432,6379,11211,27017",
	"main":    "21,22,80,443,135,139,445,1433,1521,2375,3306,3389,5432,6379,9000,11211,27017",
}

// 配置信息
type InfoScan struct {
	Hosts    string
	Ports    string
	Url      string
	HostPort []string
	//HostFile string
	//PortsFile         string
	//DirFile           string
	//OutputFile        string
	Brute             bool
	FTPReadFile       string
	FTPWriteFile      string
	SshKey            string
	Domain            string
	SkipRedis         bool
	RedisSshFile      string
	RedisCronHost     string
	RedisWebshellFile string
	RemotePath        string
	WebInfo           []string
}

var UsernameDict = map[string][]string{
	"ftp":        {"ftp", "admin", "www", "web", "root", "db", "wwwroot", "data"},
	"mysql":      {"root", "mysql"},
	"mssql":      {"sa", "sql"},
	"smb":        {"administrator", "admin", "guest"},
	"rdp":        {"administrator", "admin", "guest"},
	"postgresql": {"postgres", "admin"},
	"ssh":        {"root", "admin"},
	"mongodb":    {"root", "admin"},
	"oracle":     {"sys", "system", "admin", "test", "web", "orcl"},
}

var PasswordDict = []string{"123456", "admin", "admin123", "root", "", "pass123", "pass@123", "password", "123123", "654321", "111111", "123", "1", "admin@123", "Admin@123", "admin123!@#", "{user}", "{user}1", "{user}111", "{user}123", "{user}@123", "{user}_123", "{user}#123", "{user}@111", "{user}@2019", "{user}@123#4", "P@ssw0rd!", "P@ssw0rd", "Passw0rd", "qwe123", "12345678", "test", "test123", "123qwe", "123qwe!@#", "123456789", "123321", "666666", "a123456.", "123456~a", "123456!a", "000000", "1234567890", "8888888", "!QAZ2wsx", "1qaz2wsx", "abc123", "abc123456", "1qaz@WSX", "a11111", "a12345", "Aa1234", "Aa1234.", "Aa12345", "a123456", "a123123", "Aa123123", "Aa123456", "Aa12345.", "sysadmin", "system", "1qaz!QAZ", "2wsx@WSX", "qwe123!@#", "Aa123456!", "A123456s!", "sa123456", "1q2w3e", "Charge123", "Aa123456789"}

//type PocScan struct {
//	Target string
//	// 内置POC
//	PocName string
//	// 自定义POC
//	PocFile string
//}

//	type Exploit struct {
//		Target string
//		// 内置EXP
//		ExpName string
//		// 自定义EXP
//		ExpFile string
//	}
type Pocinfo struct {
	Target  string
	PocName string
}

type Expinfo struct {
	Target  string
	ExpName string
}

var (
	RequestPath  string
	RequestBody  string
	Username     string
	UsernameFile string
	URL          string
	URLFile      string
	HostFile     string
	Ports        string = DefaultPorts
	PortsFile    string
	Hash         string
	HashFile     string
	Hashs        []string
	HashBytes    [][]byte
	AddPorts     string
	AddUserNames string
	AddPassWords string
	Password     string
	PasswordFile string
	Accept       = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"
	// random user agent
	UserAgent           = RandomUserAgent()
	Threads      int    = 500
	ScanType     string = "all"
	Timeout      int64  = 3
	Urls         []string
	Command      string
	BruteThreads int = 1
	Cookie       string
	PocNum       int = 20
	PocInfo      Pocinfo
	PocPath      string
	PocName      string
	PocType      string
	ExpNum       int = 3
	ExpPath      string
	ExpType      string
	ExpInfo      Expinfo
	ExpName      string
	WebTimeout   int64 = 5
	NoPOC        bool
	NoExploit    bool
	//NoIOT                  bool
	DnsLog                 bool = false
	CeyeToken              string
	CeyeURL                string
	FullPOC                bool   = false
	FullEXP                bool   = false
	SaveResult             bool   = true
	OutPutFile             string = DefaultOutputFile
	EnableInfoContainer    bool   = false
	EnablePocContainer     bool   = false
	EnableExploitContainer bool   = false
	PocBody                string
	CheckPocResBody        string
	CheckExpResBody        string
	WriteWebShellBody      string
	CheckWebshellPath      string
	WebShellCommand        string
	CheckWebShellCmdBody   string
)

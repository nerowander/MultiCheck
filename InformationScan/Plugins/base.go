package Plugins

var PluginList = map[string]interface{}{
	"21":   ftpScan,
	"22":   sshScan,
	"445":  smbScan,
	"1433": mssqlScan,
	"1521": oracleScan,
	"3306": mysqlScan,
	//"3389":    rdpScan,
	"5432":    postgresScan,
	"6379":    scanRedis,
	"9000":    fcgiScan,
	"11211":   memcachedScan,
	"27017":   mongodbScan,
	"1000000": webScan,
	//"1000001": webScan,
	// todo: 可能要添加路由器的某些特定端口的扫描
}

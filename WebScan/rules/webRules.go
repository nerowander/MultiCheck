package rules

type RuleData struct {
	Name string
	Type string
	Rule string
}

type PocData struct {
	Name  string
	Alias string
}

type ExpData struct {
	Name  string
	Alias string
}

var RuleDatas = []RuleData{
	{"宝塔", "code", "(app.bt.cn/static/app.png|安全入口校验失败|<title>入口校验失败</title>|href=\"http://www.bt.cn/bbs)"},
	{"深信服防火墙类产品", "code", "(SANGFOR FW)"},
	{"360网站卫士", "code", "(webscan.360.cn/status/pai/hash|wzws-waf-cgi|zhuji.360.cn/guard/firewall/stopattack.html)"},
	{"360网站卫士", "headers", "(360wzws|CWAP-waf|zhuji.360.cn|X-Safe-Firewall)"},
	{"绿盟防火墙", "code", "(NSFOCUS NF)"},
	{"绿盟防火墙", "headers", "(NSFocus)"},
	{"Topsec-Waf", "index", `(<META NAME="Copyright" CONTENT="Topsec Network Security Technology Co.,Ltd"/>","<META NAME="DESCRIPTION" CONTENT="Topsec web UI"/>)`},
	{"Anquanbao", "headers", "(Anquanbao)"},
	{"BaiduYunjiasu", "headers", "(yunjiasu)"},
	{"BigIP", "headers", "(BigIP|BIGipServer)"},
	{"BinarySEC", "headers", "(binarysec)"},
	{"BlockDoS", "headers", "(BlockDos.net)"},
	{"CloudFlare", "headers", "(cloudflare)"},
	{"Cloudfront", "headers", "(cloudfront)"},
	{"Comodo", "headers", "(Protected by COMODO)"},
	{"IBM-DataPower", "headers", "(X-Backside-Transport)"},
	{"DenyAll", "headers", "(sessioncookie=)"},
	{"dotDefender", "headers", "(dotDefender)"},
	{"Incapsula", "headers", "(X-CDN|Incapsula)"},
	{"Jiasule", "headers", "(jsluid=)"},
	{"KONA", "headers", "(AkamaiGHost)"},
	{"ModSecurity", "headers", "(Mod_Security|NOYB)"},
	{"NetContinuum", "headers", "(Cneonction|nnCoection|citrix_ns_id)"},
	{"Newdefend", "headers", "(newdefend)"},
	{"Safe3", "headers", "(Safe3WAF|Safe3 Web Firewall)"},
	{"Safedog", "code", "(404.safedog.cn/images/safedogsite/broswer_logo.jpg)"},
	{"Safedog", "headers", "(Safedog|WAF/2.0)"},
	{"SonicWALL", "headers", "(SonicWALL)"},
	{"Stingray", "headers", "(X-Mapping-)"},
	{"Sucuri", "headers", "(Sucuri/Cloudproxy)"},
	{"Usp-Sec", "headers", "(Secure Entry Server)"},
	{"Varnish", "headers", "(varnish)"},
	{"Wallarm", "headers", "(wallarm)"},
	{"阿里云", "code", "(errors.aliyun.com)"},
	{"WebKnight", "headers", "(WebKnight)"},
	{"Yundun", "headers", "(YUNDUN)"},
	{"Yunsuo", "headers", "(yunsuo)"},
	{"Coding pages", "header", "(Coding Pages)"},
	{"启明防火墙", "code", "(/cgi-bin/webui?op=get_product_model)"},
	{"Shiro", "headers", "(=deleteMe|rememberMe=)"},
	{"Portainer(Docker管理)", "code", "(portainer.updatePassword|portainer.init.admin)"},
	{"Gogs简易Git服务", "cookie", "(i_like_gogs)"},
	{"Gitea简易Git服务", "cookie", "(i_like_gitea)"},
	{"Nexus", "code", "(Nexus Repository Manager)"},
	{"Nexus", "cookie", "(NX-ANTI-CSRF-TOKEN)"},
	{"Harbor", "code", "(<title>Harbor</title>)"},
	{"Harbor", "cookie", "(harbor-lang)"},
	{"禅道", "code", "(/theme/default/images/main/zt-logo.png|/zentao/theme/zui/css/min.css)"},
	{"禅道", "cookie", "(zentaosid)"},
	{"协众OA", "code", "(Powered by 协众OA)"},
	{"协众OA", "cookie", "(CNOAOASESSID)"},
	{"xxl-job", "code", "(分布式任务调度平台XXL-JOB)"},
	{"atmail-WebMail", "cookie", "(atmail6)"},
	{"atmail-WebMail", "code", "(/index.php/mail/auth/processlogin|Powered by Atmail)"},
	{"weblogic", "code", "(/console/framework/skins/wlsconsole/images/login_WebLogic_branding.png|Welcome to Weblogic Application Server|<i>Hypertext Transfer Protocol -- HTTP/1.1</i>)"},
	{"致远OA", "code", "(/seeyon/common/|/seeyon/USER-DATA/IMAGES/LOGIN/login.gif)"},
	{"discuz", "code", "(content=\"Discuz! X\")"},
	{"Typecho", "code", "(Typecho</a>)"},
	{"金蝶EAS", "code", "(easSessionId)"},
	{"phpMyAdmin", "cookie", "(pma_lang|phpMyAdmin)"},
	{"phpMyAdmin", "code", "(/themes/pmahomme/img/logo_right.png)"},
	{"H3C-AM8000", "code", "(AM8000)"},
	{"360企业版", "code", "(360EntWebAdminMD5Secret)"},
	{"H3C公司产品", "code", "(service@h3c.com)"},
	{"H3C ICG 1000", "code", "(ICG 1000系统管理)"},
	{"Citrix-Metaframe", "code", "(window.location=\"/Citrix/MetaFrame)"},
	{"H3C ER5100", "code", "(ER5100系统管理)"},
	{"阿里云CDN", "code", "(cdn.aliyuncs.com)"},
	{"CISCO_EPC3925", "code", "(Docsis_system)"},
	{"CISCO ASR", "code", "(CISCO ASR)"},
	{"H3C ER3200", "code", "(ER3200系统管理)"},
	{"万户oa", "code", "(/defaultroot/templates/template_system/common/css/|/defaultroot/scripts/|css/css_whir.css)"},
	{"Spark_Master", "code", "(Spark Master at)"},
	{"华为_HUAWEI_SRG2220", "code", "(HUAWEI SRG2220)"},
	{"蓝凌OA", "code", "(/scripts/jquery.landray.common.js)"},
	{"深信服ssl-vpn", "code", "(login_psw.csp)"},
	{"华为 NetOpen", "code", "(/netopen/theme/css/inFrame.css)"},
	{"Citrix-Web-PN-Server", "code", "(Citrix Web PN Server)"},
	{"juniper_vpn", "code", "(welcome.cgi?p=logo|/images/logo_juniper_reversed.gif)"},
	{"360主机卫士", "headers", "(zhuji.360.cn)"},
	{"Nagios", "headers", "(Nagios Access)"},
	{"H3C ER8300", "code", "(ER8300系统管理)"},
	{"Citrix-Access-Gateway", "code", "(Citrix Access Gateway)"},
	{"华为 MCU", "code", "(McuR5-min.js)"},
	{"TP-LINK Wireless WDR3600", "code", "(TP-LINK Wireless WDR3600)"},
	{"泛微OA", "headers", "(ecology_JSessionid)"},
	{"泛微OA", "code", "(/spa/portal/public/index.js)"},
	{"华为_HUAWEI_ASG2050", "code", "(HUAWEI ASG2050)"},
	{"360网站卫士", "code", "(360wzb)"},
	{"Citrix-XenServer", "code", "(Citrix Systems, Inc. XenServer)"},
	{"H3C ER2100V2", "code", "(ER2100V2系统管理)"},
	{"zabbix", "cookie", "(zbx_sessionid)"},
	{"zabbix", "code", "(images/general/zabbix.ico|Zabbix SIA|zabbix-server: Zabbix)"},
	{"CISCO_VPN", "headers", "(webvpn)"},
	{"360站长平台", "code", "(360-site-verification)"},
	{"H3C ER3108GW", "code", "(ER3108GW系统管理)"},
	{"o2security_vpn", "headers", "(client_param=install_active)"},
	{"H3C ER3260G2", "code", "(ER3260G2系统管理)"},
	{"H3C ICG1000", "code", "(ICG1000系统管理)"},
	{"CISCO-CX20", "code", "(CISCO-CX20)"},
	{"H3C ER5200", "code", "(ER5200系统管理)"},
	{"linksys-vpn-bragap14-parintins", "code", "(linksys-vpn-bragap14-parintins)"},
	{"360网站卫士常用前端公共库", "code", "(libs.useso.com)"},
	{"H3C ER3100", "code", "(ER3100系统管理)"},
	{"H3C-SecBlade-FireWall", "code", "(js/MulPlatAPI.js)"},
	{"360webfacil_360WebManager", "code", "(publico/template/)"},
	{"Citrix_Netscaler", "code", "(ns_af)"},
	{"H3C ER6300G2", "code", "(ER6300G2系统管理)"},
	{"H3C ER3260", "code", "(ER3260系统管理)"},
	{"华为_HUAWEI_SRG3250", "code", "(HUAWEI SRG3250)"},
	{"exchange", "code", "(/owa/auth.owa|Exchange Admin Center)"},
	{"Spark_Worker", "code", "(Spark Worker at)"},
	{"H3C ER3108G", "code", "(ER3108G系统管理)"},
	{"Citrix-ConfProxy", "code", "(confproxy)"},
	{"360网站安全检测", "code", "(webscan.360.cn/status/pai/hash)"},
	{"H3C ER5200G2", "code", "(ER5200G2系统管理)"},
	{"华为（HUAWEI）安全设备", "code", "(sweb-poclib/resource/)"},
	{"华为（HUAWEI）USG", "code", "(UI_component/commonDefine/UI_regex_define.js)"},
	{"H3C ER6300", "code", "(ER6300系统管理)"},
	{"华为_HUAWEI_ASG2100", "code", "(HUAWEI ASG2100)"},
	{"TP-Link 3600 DD-WRT", "code", "(TP-Link 3600 DD-WRT)"},
	{"NETGEAR WNDR3600", "code", "(NETGEAR WNDR3600)"},
	{"H3C ER2100", "code", "(ER2100系统管理)"},
	{"jira", "code", "(jira.webresources)"},
	{"金和协同管理平台", "code", "(金和协同管理平台)"},
	{"Citrix-NetScaler", "code", "(NS-CACHE)"},
	{"linksys-vpn", "headers", "(linksys-vpn)"},
	{"通达OA", "code", "(/static/images/tongda.ico|http://www.tongda2000.com|通达OA移动版|Office Anywhere)"},
	{"华为（HUAWEI）Secoway设备", "code", "(Secoway)"},
	{"华为_HUAWEI_SRG1220", "code", "(HUAWEI SRG1220)"},
	{"H3C ER2100n", "code", "(ER2100n系统管理)"},
	{"H3C ER8300G2", "code", "(ER8300G2系统管理)"},
	{"金蝶政务GSiS", "code", "(/kdgs/script/kdgs.js)"},
	{"Jboss", "code", "(Welcome to JBoss|jboss.css)"},
	{"Jboss", "headers", "(JBoss)"},
	{"泛微E-mobile", "code", "(Weaver E-mobile|weaver,e-mobile)"},
	{"泛微E-mobile", "headers", "(EMobileServer)"},
	{"齐治堡垒机", "code", "(logo-icon-ico72.png|resources/themes/images/logo-login.png)"},
	//{"ThinkPHP", "headers", "(ThinkPHP)"},
	{"ThinkPHP", "code", "(ThinkPHP)"},
	{"ThinkPHP", "code", "(/Public/static/client.js/)"},
	{"ThinkPHP", "headers", "(think_lang)"},
	{"weaver-ebridge", "code", "(e-Bridge,http://wx.weaver)"},
	{"Laravel", "headers", "(laravel_session)"},
	{"DWR", "code", "(dwr/engine.js)"},
	{"swagger_ui", "code", "(swagger-ui/css|\"swagger\":|swagger-ui.min.js)"},
	{"大汉版通发布系统", "code", "(大汉版通发布系统|大汉网络)"},
	{"druid", "code", "(druid.index|DruidDrivers|DruidVersion|Druid Stat Index)"},
	{"Jenkins", "code", "(Jenkins)"},
	{"红帆OA", "code", "(iOffice)"},
	{"VMware vSphere", "code", "(VMware vSphere)"},
	{"打印机", "code", "(打印机|media/canon.gif)"},
	{"finereport", "code", "(isSupportForgetPwd|FineReport,Web Reporting Tool)"},
	{"蓝凌OA", "code", "(蓝凌软件|StylePath:\"/resource/style/default/\"|/resource/customization|sys/ui/extend/theme/default/style/profile.css|sys/ui/extend/theme/default/style/icon.css)"},
	{"GitLab", "code", "(href=\"https://about.gitlab.com/)"},
	{"Jquery-1.7.2", "code", "(/webui/js/jquerylib/jquery-1.7.2.min.js)"},
	{"Hadoop Applications", "code", "(/cluster/app/application)"},
	{"海昌OA", "code", "(/loginmain4/js/jquery.min.js)"},
	{"帆软报表", "code", "(WebReport/login.html|ReportServer)"},
	{"帆软报表", "headers", "(数据决策系统)"},
	{"华夏ERP", "headers", "(华夏ERP)"},
	{"金和OA", "cookie", "(ASPSESSIONIDSSCDTDBS)"},
	{"久其财务报表", "code", "(netrep/login.jsp|/netrep/intf)"},
	{"若依管理系统", "code", "(ruoyi/login.js|ruoyi/js/ry-ui.js)"},
	{"启莱OA", "code", "(js/jQselect.js|js/jquery-1.4.2.min.js)"},
	{"智慧校园管理系统", "code", "(DC_Login/QYSignUp)"},
	{"JQuery-1.7.2", "code", "(webui/js/jquerylib/jquery-1.7.2.min.js)"},
	{"浪潮 ClusterEngineV4.0", "code", "(0;url=module/login/login.html)"},
	{"会捷通云视讯平台", "code", "(him/api/rest/v1.0/node/role|him.app)"},
	{"源码泄露账号密码 F12查看", "code", "(get_dkey_passwd)"},
	{"Smartbi Insight", "code", "(smartbi.gcf.gcfutil)"},
	{"汉王人脸考勤管理系统", "code", "(汉王人脸考勤管理系统|/Content/image/hanvan.png|/Content/image/hvicon.ico)"},
	{"亿赛通-电子文档安全管理系统", "code", "(电子文档安全管理系统|/CDGServer3/index.jsp|/CDGServer3/SysConfig.jsp|/CDGServer3/help/getEditionInfo.jsp)"},
	{"天融信 TopApp-LB 负载均衡系统", "code", "(TopApp-LB 负载均衡系统)"},
	{"中新金盾信息安全管理系统", "code", "(中新金盾信息安全管理系统|中新网络信息安全股份有限公司)"},
	{"好视通", "code", "(深圳银澎云计算有限公司|itunes.apple.com/us/app/id549407870|hao-shi-tong-yun-hui-yi-yuan)"},
	{"蓝海卓越计费管理系统", "code", "(蓝海卓越计费管理系统|星锐蓝海网络科技有限公司)"},
	{"和信创天云桌面系统", "code", "(和信下一代云桌面VENGD|/vesystem/index.php)"},
	{"金山", "code", "(北京猎鹰安全科技有限公司|金山终端安全系统V9.0Web控制台|北京金山安全管理系统技术有限公司|金山V8)"},
	{"WIFISKY-7层流控路由器", "code", "(深圳市领空技术有限公司|WIFISKY 7层流控路由器)"},
	{"MetInfo-米拓建站", "code", "(MetInfo|/skin/style/metinfo.css|/skin/style/metinfo-v2.css)"},
	{"IBM-Lotus-Domino", "code", "(/mailjump.nsf|/domcfg.nsf|/names.nsf|/homepage.nsf)"},
	{"APACHE-kylin", "code", "(url=kylin)"},
	{"C-Lodop打印服务系统", "code", "(/CLodopfuncs.js|www.c-lodop.com)"},
	{"HFS", "code", "(href=\"http://www.rejetto.com/hfs/)"},
	{"Jellyfin", "code", "(content=\"http://jellyfin.org\")"},
	{"FIT2CLOUD-JumpServer-堡垒机", "code", "(<title>JumpServer</title>)"},
	{"Alibaba Nacos", "code", "(<title>Nacos</title>)"},
	{"Nagios", "headers", "(nagios admin)"},
	{"Pulse Connect Secure", "code", "(/dana-na/imgs/space.gif)"},
	{"h5ai", "code", "(powered by h5ai)"},
	{"jeesite", "cookie", "(jeesite.session.id)"},
	{"拓尔思SSO", "cookie", "(trsidsssosessionid)"},
	{"拓尔思WCMv7/6", "cookie", "(com.trs.idm.coSessionId)"},
	{"天融信脆弱性扫描与管理系统", "code", "(/js/report/horizontalReportPanel.js)"},
	{"天融信网络审计系统", "code", "(onclick=dlg_download())"},
	{"天融信日志收集与分析系统", "code", "(天融信日志收集与分析系统)"},
	{"URP教务系统", "code", "(北京清元优软科技有限公司)"},
	{"科来RAS", "code", "(科来软件 版权所有|i18ninit.min.js)"},
	{"正方OA", "code", "(zfoausername)"},
	{"希尔OA", "code", "(/heeroa/login.do)"},
	{"泛普建筑工程施工OA", "code", "(/dwr/interface/LoginService.js)"},
	{"中望OA", "code", "(/IMAGES/default/first/xtoa_logo.png|/app_qjuserinfo/qjuserinfoadd.jsp)"},
	{"海天OA", "code", "(HTVOS.js)"},
	{"信达OA", "code", "(http://www.xdoa.cn</a>)"},
	{"任我行CRM", "code", "(CRM_LASTLOGINUSERKEY)"},
	{"Spammark邮件信息安全网关", "code", "(/cgi-bin/spammark?empty=1)"},
	{"winwebmail", "code", "(WinWebMail Server|images/owin.css)"},
	{"浪潮政务系统", "code", "(LangChao.ECGAP.OutPortal|OnlineQuery/QueryList.aspx)"},
	{"天融信防火墙", "code", "(/cgi/maincgi.cgi)"},
	{"网神防火墙", "code", "(css/lsec/login.css)"},
	{"帕拉迪统一安全管理和综合审计系统", "code", "(module/image/pldsec.css)"},
	{"蓝盾BDWebGuard", "code", "(BACKGROUND: url(images/loginbg.jpg) #e5f1fc)"},
	{"Huawei SMC", "code", "(Script/SmcScript.js?version=)"},
	{"coremail", "code", "(/coremail/bundle/|contextRoot: \"/coremail\"|coremail/common)"},
	{"activemq", "code", "(activemq_logo|Manage ActiveMQ broker)"},
	{"锐捷网络", "code", "(static/img/title.ico|support.ruijie.com.cn|Ruijie - NBR|eg.login.loginBtn)"},
	{"禅道", "code", "(/theme/default/images/main/zt-logo.png|zentaosid)"},
	{"weblogic", "code", "(/console/framework/skins/wlsconsole/images/login_WebLogic_branding.png|Welcome to Weblogic Application Server|<i>Hypertext Transfer Protocol -- HTTP/1.1</i>|<TITLE>Error 404--Not Found</TITLE>|Welcome to Weblogic Application Server|<title>Oracle WebLogic Server 管理控制台</title>)"},
	{"weblogic", "headers", "(WebLogic)"},
	{"致远OA", "code", "(/seeyon/USER-DATA/IMAGES/LOGIN/login.gif|/seeyon/common/)"},
	{"蓝凌EIS智慧协同平台", "code", "(/scripts/jquery.landray.common.js)"},
	{"深信服ssl-vpn", "code", "(login_psw.csp|loginPageSP/loginPrivacy.js|/por/login_psw.csp)"},
	{"Struts2", "code", "(org.apache.struts2|Struts Problem Report|struts.devMode|struts-tags|There is no Action mapped for namespace)"},
	{"泛微OA", "code", "(/spa/portal/public/index.js|wui/theme/ecology8/page/images/login/username_wev8.png|/wui/index.html#/?logintype=1)"},
	{"Swagger UI", "code", "(/swagger-ui.css|swagger-ui-bundle.js|swagger-ui-standalone-preset.js)"},
	{"金蝶政务GSiS", "code", "(/kdgs/script/kdgs.js|HTML5/content/themes/kdcss.min.css|/ClientBin/Kingdee.BOS.XPF.App.xap)"},
	{"蓝凌OA", "code", "(蓝凌软件|StylePath:\"/resource/style/default/\"|/resource/customization|sys/ui/extend/theme/default/style/icon.css|sys/ui/extend/theme/default/style/profile.css)"},
	{"用友NC", "code", "(Yonyou UAP|YONYOU NC|/Client/Uclient/UClient.dmg|logo/images/ufida_nc.png|iufo/web/css/menu.css|/System/Login/Login.asp?AppID=|/nc/servlet/nc.ui.iufo.login.Index)"},
	{"用友IUFO", "code", "(iufo/web/css/menu.css)"},
	{"TELEPORT堡垒机", "code", "(/static/plugins/blur/background-blur.js)"},
	{"JEECMS", "code", "(/r/cms/www/red/js/common.js|/r/cms/www/red/js/indexshow.js|Powered by JEECMS|JEECMS|/jeeadmin/jeecms/index.do)"},
	{"CMS", "code", "(Powered by .*CMS)"},
	{"目录遍历", "code", "(Directory listing for /)"},
	{"ATLASSIAN-Confluence", "code", "(com.atlassian.confluence)"},
	{"ATLASSIAN-Confluence", "headers", "(X-Confluence)"},
	{"向日葵", "code", "({\"success\":false,\"msg\":\"Verification failure\"})"},
	{"Kubernetes", "code", "(Kubernetes Dashboard</title>|Kubernetes Enterprise Manager|Mirantis Kubernetes Engine|Kubernetes Resource Report)"},
	{"WordPress", "code", "(/wp-login.php?action=lostpassword|WordPress</title>)"},
	{"RabbitMQ", "code", "(RabbitMQ Management)"},
	{"dubbo", "headers", "(Basic realm=\"dubbo\")"},
	{"Spring env", "code", "(logback)"},
	{"ueditor", "code", "(ueditor.all.js|UE.getEditor)"},
	{"亿邮电子邮件系统", "code", "(亿邮电子邮件系统|亿邮邮件整体解决方案)"},
	{"apache", "headers", "Apache/"},
}

var PocDatas = []PocData{
	{"致远OA", "seeyon"},
	{"泛微OA", "weaver"},
	{"通达OA", "tongda"},
	{"蓝凌OA", "landray"},
	{"ThinkPHP", "thinkphp"},
	{"Nexus", "nexus"},
	{"齐治堡垒机", "qizhi"},
	{"weaver-ebridge", "weaver-ebridge"},
	{"weblogic", "weblogic"},
	{"zabbix", "zabbix"},
	{"VMware vSphere", "vmware"},
	{"Jboss", "jboss"},
	{"用友", "yongyou"},
	{"用友IUFO", "yongyou"},
	{"coremail", "coremail"},
	{"金山", "kingsoft"},
	{"apache", "apache"},
}

var ExpDatas = []ExpData{
	{"致远OA", "seeyon"},
	{"泛微OA", "weaver"},
	{"通达OA", "tongda"},
	{"蓝凌OA", "landray"},
	{"ThinkPHP", "thinkphp"},
	{"Nexus", "nexus"},
	{"齐治堡垒机", "qizhi"},
	{"weaver-ebridge", "weaver-ebridge"},
	{"weblogic", "weblogic"},
	{"zabbix", "zabbix"},
	{"VMware vSphere", "vmware"},
	{"Jboss", "jboss"},
	{"用友", "yongyou"},
	{"用友IUFO", "yongyou"},
	{"coremail", "coremail"},
	{"金山", "kingsoft"},
	{"apache", "apache"},
}

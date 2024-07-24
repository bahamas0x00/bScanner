package bScanner

import "strings"

var tmpSuffixFormat = []string{".zip", ".rar", ".tar.gz", ".tgz", ".tar.bz2", ".tar", ".jar", ".war", ".7z", ".bak", ".sql", ".gz", ".sql.gz", ".tar.tgz", ".backup"}

var tmpInfoDic = []string{"0", "00", "000", "012", "1", "111", "123", "127.0.0.1", "2", "2010", "2011", "2012", "2013", "2014", "2015", "2016", "2017", "2018", "2019", "2020", "2021", "2022", "2023", "2024", "2025", "234", "3", "333", "4", "444", "5", "555", "6", "666", "7", "777", "8", "888", "9", "999", "a", "about", "admin", "app", "application", "archive", "asp", "aspx", "auth", "b", "back", "backup", "backups", "bak", "bbs", "beifen", "bin", "cache", "clients", "code", "com", "config", "core", "customers", "dat", "data", "database", "db", "download", "dump", "engine", "error_log", "extend", "files", "forum", "ftp", "home", "html", "img", "include", "index", "install", "joomla", "js", "jsp", "local", "login", "localhost", "master", "media", "members", "my", "mysql", "new", "old", "orders", "output", "package", "php", "public", "root", "runtime", "sales", "server", "shujuku", "site", "sjk", "sql", "store", "tar", "template", "test", "upload", "user", "users", "vb", "vendor", "wangzhan", "web", "website", "wordpress", "wp", "www", "wwwroot", "wz", "log", "数据库", "数据库备份", "网站", "网站备份"}

// creat static wordlist
func generateStaticWordlist() []string {
	staticWordlist := make([]string, 0)
	for _, suffix := range tmpSuffixFormat {
		for _, infoDic := range tmpInfoDic {
			staticWordlist = append(staticWordlist, infoDic+suffix)
		}
	}
	return staticWordlist
}

// create domain wordlist
func generateDomainWordlist(url string) []string {
	domainWordlist := make([]string, 0)
	domainInfoDic := make([]string, 0)
	url = trimUrlPrefix(url)
	wordsInDomain := strings.Split(url, ".")
	sld := strings.SplitAfterN(url, ".", 2)
	sld_ := make([]string, 0)
	for _, v := range sld {
		sld_ = append(sld_, strings.Replace(v, ".", "_", -1))
	}

	domainInfoDic = append(domainInfoDic, url, strings.ReplaceAll(url, ".", ""), strings.ReplaceAll(url, ".", "_"), wordsInDomain[0], wordsInDomain[1])
	domainInfoDic = append(domainInfoDic, wordsInDomain...)
	domainInfoDic = append(domainInfoDic, sld_...)

	for _, suffix := range tmpSuffixFormat {
		for _, infoDic := range domainInfoDic {
			domainWordlist = append(domainWordlist, infoDic+suffix)
		}
	}

	return domainWordlist
}

func trimUrlPrefix(url string) string {
	prefixes := []string{"http://", "https://"}

	for _, prefix := range prefixes {
		if strings.HasPrefix(url, prefix) {
			return strings.TrimPrefix(url, prefix)
		}
	}
	return url
}

/*
 host , host替换. 为“” ， host替换. 为_  , 域名词组， 域名去除前缀 ， 域名去除前缀，后面的. 替换为 _ , 域名词组第一个， 域名词组第二个
*/

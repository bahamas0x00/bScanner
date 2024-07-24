/*
single url scan : ./backupFileScanner -u https://www.example.com -o result.txt
multiple urls scan : ./backupFileScanner -t 100 -f url.txt -o result.txt
options:

	-h --help 				check for usage help
	-f --url-file 			multiple scan , each line of url should contain prefix http:// or https:// , default http://
	-t --thread 			scanner threads
	-u --url 				single url scan
	-d --dict-file 			use custom wordlist
	-o --output-file 		output filename
	-p --proxy 				use proxy , eg : socks5://127.0.0.1:1080
*/
package main

import (
	"backupFileScanner/bScanner"
	"flag"
	"fmt"
)

var (
	urlFile    string
	thread     int
	url        string
	dictFile   string
	outputFile string
	proxy      string
)

func main() {

	flag.StringVar(&urlFile, "f", "", "multiple scan , each line of url should contain prefix http:// or https:// , default http://")
	flag.IntVar(&thread, "t", 50, "scanner threads")
	flag.StringVar(&url, "u", "", "single url scan")
	flag.StringVar(&dictFile, "d", "", "use custom wordlist")
	flag.StringVar(&outputFile, "o", "result.txt", "output filename")
	flag.StringVar(&proxy, "p", "", "use proxy , eg:socks5://127.0.0.1:1080")

	flag.Parse()

	if url == "" && urlFile == "" {
		fmt.Println("[!] Please specify a URL (-u) or URL file (-f)")
		flag.PrintDefaults()
		return
	}

	if url != "" && urlFile != "" {
		fmt.Println("[!] Can't use -u and -f at the same time,  Please specify a URL (-u) or URL file (-f)")
		flag.PrintDefaults()
		return
	}

	bScanner := bScanner.NewbScanner(urlFile, thread, url, dictFile, outputFile, proxy)

	bScanner.Scan()

}

package bScanner

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	client = &http.Client{Timeout: 10 * time.Second}

	randomSeed = rand.New(rand.NewSource(time.Now().UnixNano()))

	fileMutex sync.Mutex

	defaultPrefix = "http://"
)

type bScanner struct {
	urlFile    string
	thread     int
	url        string
	dictFile   string
	outputFile string
	proxy      string
}

func NewbScanner(urlFile string, thread int, url string, dictFile string, outputFile string, proxy string) *bScanner {
	return &bScanner{
		urlFile:    urlFile,
		thread:     thread,
		url:        url,
		dictFile:   dictFile,
		outputFile: outputFile,
		proxy:      proxy,
	}
}

func (b *bScanner) generateUrlCheckList() []string {
	urlCheckList := make([]string, 0)
	var url = b.url

	// single url
	if url != "" {
		// 检查前后缀，url -> http://url/ || https://url/
		url = urlCheck(url)
		urlCheckList = append(urlCheckList, url)

	} else {
		urlList, err := readLines(b.urlFile)
		if err != nil {
			log.Fatalf("Error opening urlFile: %v", err)
		}

		for _, url := range urlList {
			url = urlCheck(url)
			urlCheckList = append(urlCheckList, url)
		}

	}

	return urlCheckList
}

func (b *bScanner) generateWordlist() []string {
	fullWordlist := make([]string, 0)
	domainWordlist := make([]string, 0)

	// 生成静态字典
	staticWordlist := generateStaticWordlist()
	fullWordlist = append(fullWordlist, staticWordlist...)

	// 添加域名相关字典
	if b.url != "" {
		url := b.url
		domainWordlist = generateDomainWordlist(url)
	} else {
		urlList, err := readLines(b.urlFile)
		if err != nil {
			log.Fatalf("Error opening urlFile: %v", err)
		}

		for _, url := range urlList {
			wordlist := generateDomainWordlist(url)
			domainWordlist = append(domainWordlist, wordlist...)
		}
	}
	fullWordlist = append(fullWordlist, domainWordlist...)

	// 添加自定义字典
	customWordlist, err := readLines(b.dictFile)
	if err != nil {
		log.Fatalf("Error opening custom wordlist: %v", err)
	}
	fullWordlist = append(fullWordlist, customWordlist...)

	fmt.Println(fullWordlist)
	return fullWordlist
}

func (b *bScanner) Scan() {
	urlCheckList := b.generateUrlCheckList()
	fullWordlist := b.generateWordlist()

	checkQueue := make([]string, 0)

	for _, u := range urlCheckList {
		for _, suffix := range fullWordlist {
			checkQueue = append(checkQueue, u+suffix)
		}
	}

	if b.proxy != "" {
		proxyURL, err := url.Parse(b.proxy)
		if err != nil {
			log.Fatal("Error parsing proxy URL: ", err)
		}
		client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, b.thread)

	for _, checkUrl := range checkQueue {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(checkUrl string) {
			defer func() {
				<-semaphore
				wg.Done()
			}()

			err := vuln(checkUrl, b.outputFile)
			if err != nil {
				log.Printf("[ fail] %s: %v\n", checkUrl, err)
			}
		}(checkUrl)
	}

	wg.Wait()
	close(semaphore)

	fmt.Println("All URLs processed successfully.")

}

func vuln(url, outputFile string) error {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", generateRandomUserAgent())

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(contentType, "html") ||
		strings.Contains(contentType, "image") ||
		strings.Contains(contentType, "xml") ||
		strings.Contains(contentType, "text") ||
		strings.Contains(contentType, "json") ||
		strings.Contains(contentType, "javascript") {
		return fmt.Errorf("unexpected content type: %s", contentType)
	}

	contentLength := resp.Header.Get("Content-Length")
	if contentLength == "" {
		return fmt.Errorf("content length not provided")
	}

	// convert content length to integer
	tmpRarSize, err := strconv.Atoi(contentLength)
	if err != nil {
		return fmt.Errorf("error converting content length")
	}

	// check size condition
	rarSize := size(tmpRarSize)
	sizeParts := strings.Split(rarSize, " ")
	if len(sizeParts) != 2 {
		return fmt.Errorf("invalid content size format")
	}

	sizeValue, err := strconv.Atoi(sizeParts[0])
	if err != nil {
		return fmt.Errorf("error parsing content size")
	}

	if sizeValue <= 0 {
		return fmt.Errorf("invalid content size")
	}

	log.Printf("[ success ] %s size: %s\n", url, rarSize)
	if err := writeToFile(outputFile, fmt.Sprintf("%s  size: %s\n", url, rarSize)); err != nil {
		return fmt.Errorf("error writing to file")
	}

	return nil

}

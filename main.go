package main

import (
	"bufio"
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"net/url"
	"os"
	"sync"
)

func main() {
	filename := "links.txt"
	links, err := getLinksFromFile(filename)
	if err != nil {
		fmt.Println(err)
	}

	var waitGroup sync.WaitGroup
	waitGroup.Add(len(links))

	linkChan := make(chan string)
	for _, link := range links {
		go getLinksFromUrl(&waitGroup, linkChan, link)
	}

	for url := range linkChan {
		if IsUrlValid(url) {
			fmt.Println(url)
			go getLinksFromUrl(&waitGroup, linkChan, url)
		}
	}

	waitGroup.Wait()
}

func IsUrlValid(urlStr string) bool {
	u, err := url.Parse(urlStr)
	return err == nil && u.Scheme != "" && u.Host != ""
}

//func hitURL(waitGroup *sync.WaitGroup, url string) {
//	resp, err := http.Get(url)
//	if err != nil {
//		fmt.Println("Error fetching from URL : ", url)
//	}
//	defer resp.Body.Close()
//	links, err := getLinksFromBody(resp.Body)
//	for _, link := range links {
//		fmt.Println(link)
//	}
//	waitGroup.Done()
//}

func getLinksFromUrl(waitGroup *sync.WaitGroup, linkChan chan string, urlStr string) {

	resp, err := http.Get(urlStr)
	if err != nil {
		fmt.Println("Error fetching from URL : ", urlStr)
		return
	}
	body := resp.Body
	defer body.Close()

	z := html.NewTokenizer(body)
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			continue
		case html.StartTagToken:
			token := z.Token()
			if token.Data != "a" {
				continue
			}
			for _, attr := range token.Attr {
				if attr.Key == "href" {
					linkChan <- attr.Val
				}
			}
		}
	}

}

func getLinksFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var links []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		links = append(links, scanner.Text())
	}
	return links, nil
}

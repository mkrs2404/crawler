package main

import (
	"bufio"
	"fmt"
	"golang.org/x/net/html"
	"log"
	"net/http"
	_ "net/http/pprof"
	"net/url"
	"os"
	"sync"
)

const maxRoutines = 5

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	filename := "links.txt"
	links, err := getLinksFromFile(filename)
	if err != nil {
		fmt.Println(err)
	}

	var waitGroup sync.WaitGroup
	waitGroup.Add(maxRoutines)

	linkChan := make(chan string)

	for i := 1; i <= maxRoutines; i++ {
		go worker(linkChan)
	}

	for _, link := range links {
		if IsUrlValid(link) {
			linkChan <- link
		}
	}
	waitGroup.Wait()
}

func IsUrlValid(urlStr string) bool {
	u, err := url.Parse(urlStr)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func worker(linkChan chan string) {
	for link := range linkChan {
		links, err := getLinksFromUrl(link)
		if err != nil {
			fmt.Println("Error fetching from URL : ", link)
		} else {
			go queueLinks(linkChan, links)
		}
	}

}

func queueLinks(linkChan chan string, links []string) {
	for _, link := range links {
		if IsUrlValid(link) {
			fmt.Println(link)
			linkChan <- link
		}
	}
}

func getLinksFromUrl(urlStr string) ([]string, error) {

	var links []string
	resp, err := http.Get(urlStr)
	if err != nil {
		return links, err
	}
	body := resp.Body
	defer body.Close()

	z := html.NewTokenizer(body)
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			return links, err
		case html.StartTagToken:
			token := z.Token()
			if token.Data != "a" {
				continue
			}
			for _, attr := range token.Attr {
				if attr.Key == "href" {
					links = append(links, attr.Val)
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

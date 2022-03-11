package main

import (
	"bufio"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"os"
)

func main() {
	filename := "links.txt"
	links, err := getLinksFromFile(filename)
	if err != nil {
		fmt.Println(err)
	}
	for _, link := range links {
		hitURL(link)
	}
}

func hitURL(url string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching from URL : ", url)
	}
	defer resp.Body.Close()
	links, err := getLinksFromBody(resp.Body)
	for _, link := range links {
		fmt.Println(link)
	}
}

func getLinksFromBody(body io.Reader) ([]string, error) {

	var links []string
	z := html.NewTokenizer(body)
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			return links, nil
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

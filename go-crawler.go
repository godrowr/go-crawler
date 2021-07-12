package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/net/html"
)

func main() {
	flag.Parse()
	args := flag.Args()
	fmt.Println("Scraping", args[0])
	if len(args) < 1 {
		fmt.Println("No input arguement supplied")
		os.Exit(1)
	}
	map := make(map[string]string)
	ch := make(chan string)
	find_links(args[0])
}

func find_links(url string) {
	url, err := url.ParseRequestURI(args[0])
	if err != nil {
		fmt.Println("String inputted was not a URL", err)
		return
	}
	resp, _ := http.Get(url.String())
	


	// Channel here with a goroutine that finds links on a page. 

	z := html.NewTokenizer(resp.Body)
	for {
		tt := z.Next()
		switch {
		case tt == html.ErrorToken:
			return
		case tt == html.StartTagToken:
			if z.Token().Data == "a" {
				for _, a := range z.Token().Attr {
					if a.Key == "href" {
						fmt.Println("Found href:", a.Val)
						break
					}
				}
			}
		}
	}
}

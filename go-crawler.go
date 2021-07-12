package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/html"
)

//var wg *sync.WaitGroup
// var scraped map[string]bool
// var host []string
// var wg sync.WaitGroup

func main() {
	flag.Parse()
	args := flag.Args()
	fmt.Println("Scraping", args[0])
	if len(args) < 1 {
		fmt.Println("No input arguement supplied")
		os.Exit(1)
	}

	scraped = make(map[string]bool)
	ch_queue := make(chan string)
	ch_done := make(chan string)

	check, base := is_url(args[0])
	if check {
		host = strings.SplitN(base.String(), "//", 2)
		wg.Add(1)
		go find_links(ch_queue, ch_done, base.String())
	} else {
		fmt.Println("Input invalid - HTTP/s URL only")
		os.Exit(1)
	}

	// for {
	// 	select {
	// 	case link := <-ch_queue:
	// 		scraped[link] = true
	// 	case link := <-ch_done:
	// 		fmt.Println(link)
	// 	}
	// }

	for item := range ch_queue {
		wg.Add(1)
		go find_links(ch_queue, ch_done, item)
	}
	wg.Wait()
	close(ch_queue)

}

func find_links(ch_queue chan string, ch_done chan string, site string) {
	defer wg.Done()
	fmt.Println(site)
	// defer func() {
	// 	ch_done <- site
	// }()
	if scraped[site] {
		return
	}
	check, url := is_url(site) //check
	if strings.Index(site, host[1]) != 0 && !check {
		return
	}

	resp, err := http.Get(url.String())
	if err != nil {
		return
	}

	scraped[site] = true

	z := html.NewTokenizer(resp.Body)
	defer resp.Body.Close()
	for {
		tt := z.Next()
		switch tt {
		case html.StartTagToken:
			t := z.Token()
			if t.Data == "a" {
				for _, a := range t.Attr {
					if a.Key == "href" {
						if a.Val == "" {
							return
						}

						if strings.Index(a.Val, "http") != 0 {
							string := strings.Builder{}
							string.WriteString("http://")
							string.WriteString(host[1])
							string.WriteString(a.Val)
							a.Val = string.String()
						}
						defer fmt.Println("---", a.Val)
						ch_queue <- a.Val

						//defer fmt.Println("---", a.Val)
						// _, exists := scraped[a.Val]
						// if !exists {
						// 	scraped[a.Val] = false
						// 	//ch_queue <- a.Val
						// } else if !scraped[a.Val] {
						// 	break
						// }
					}
				}
			}
		}
	}
}

func is_url(site string) (bool, *url.URL) {
	url, err := url.ParseRequestURI(site)
	if err != nil { //&& strings.Index(site, "http") != 0
		fmt.Println("String inputted was not a HTTP/s URL", err)
		return false, url
	} else {
		return true, url
	}
}

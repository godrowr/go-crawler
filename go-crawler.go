package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

func main() {
	flag.Parse()
	args := flag.Args()
	fmt.Println("Scraping", args[0])
	if len(args) < 1 {
		fmt.Println("No input arguement supplied")
		os.Exit(1)
	}

	base, err := url.Parse(args[0])
	if err != nil {
		fmt.Println("URL has no base host", err)
		return
	}

	host := strings.SplitN(base.String(), "//", 2)

	c := colly.NewCollector(
		colly.AllowedDomains(host[1]),
		colly.Async(true),
	)

	c.URLFilters = []*regexp.Regexp{
		regexp.MustCompile(`https?://`),
	}

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		url, err := url.ParseRequestURI(link)
		if err == nil {
			URL := base.ResolveReference(url).String()
			if url.Scheme == "https" || url.Scheme == "http" {
				fmt.Println("  " + URL)
			}
			c.Visit(e.Request.AbsoluteURL(URL))
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(r.URL.String())
	})
	c.Visit(args[0])
	c.Wait()
}

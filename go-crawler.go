package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/html"
)

var links = make(map[string]string)

var host string

func main() {
	fmt.Println("Enter site to crawl (ex: http://www.rescale.com/):")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("An error occured while reading input. Please try again", err)
		return
	}
	input = strings.TrimSuffix(input, "\n")
	u, err := url.ParseRequestURI(input)
	if err != nil {
		fmt.Println("String inputted was not a URL", err)
		return
	}
	host = u.Hostname()
	loop(u.String())

}

func loop(input string) {
	crawl(input)
	for URL, visited := range links { //TODO
		if visited == "Not-Visited" && strings.Contains(URL, host) {
			links[URL] = "Visited"
			crawl(URL)
		} else {
			links[URL] = "Visited"
		}
	}
	for URL, visited := range links {
		if visited == "Not-Visited" {
			loop(URL)
		}
	}
	os.Exit(0)
}

func crawl(input string) {
	fmt.Println(input)
	links[input] = "Visited"
	resp, err := http.Get(input)
	if err != nil {
		fmt.Println("GET request could not be processed correctly", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("BODY could not be obtained from GET request", err)
		return
	}

	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		fmt.Println("DOC could not be obtained from BODY", err)
		return
	}
	f(doc, input)
}

func f(n *html.Node, input string) {
	base, err := url.Parse(input)
	if err != nil {
		log.Fatal(err)
	}
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" && len(a.Val) > 0 {
				url, err := url.ParseRequestURI(a.Val)
				if err == nil {
					URL := base.ResolveReference(url).String()
					if url.Scheme != "tel" && strings.Compare(URL, input) != 0 {
						if links[URL] == "" || links[URL] != "Visited" {
							links[URL] = "Not-Visited"
						}
						fmt.Println(" " + URL + " " + links[URL])
					}
				}
				break
			}

		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		f(c, input)
	}
}

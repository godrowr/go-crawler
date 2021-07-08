package main

import (
	"flag"
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

/*
Function that reads input string, verifies if it is a URL
then calls map_loop to begin crawling the sites and inserting the
results into a map.
*/
func main() {
	flag.Parse()
	args := flag.Args()
	fmt.Println("Scraping", args[0])
	if len(args) < 1 {
		fmt.Println("No input arguement supplied")
		os.Exit(1)
	}
	u, err := url.ParseRequestURI(args[0])
	if err != nil {
		fmt.Println("String inputted was not a URL", err)
		return
	}
	host = u.Hostname()
	map_loop(u.String())

}

/*
Function that takes a string and crawls a site, inserts URL into a map.
The map_loop handles updating the sites that are visited and only visits sites
with the same domain as the original input string.
*/
func map_loop(input string) {
	crawl_site(input)
	for URL, visited := range links { //TODO
		if visited == "Not-Visited" && strings.Contains(URL, host) {
			links[URL] = "Visited"
			crawl_site(URL)
		} else {
			links[URL] = "Visited"
		}
	}
	for URL, visited := range links {
		if visited == "Not-Visited" {
			map_loop(URL)
		}
	}
	os.Exit(0)
}

/* The crawl_site function fetches the HTML document, and
sends the body to search_links to be parsed.
*/
func crawl_site(input string) {
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
	search_links(doc, input)
}

/* The search_links function is recursively called to search for
<a> nodes in an html body segment of a url. It omits schemes that
are not http or https or are identical to the orginal input.

seach_links outputs the sites found on the webpage and adds new sites
to the map if it doesn't exist in the map currently.
*/
func search_links(n *html.Node, input string) {
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
		search_links(c, input)
	}
}

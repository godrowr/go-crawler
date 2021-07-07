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

func main() {
	fmt.Println("Enter site to crawl:")
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
	crawl(u)
}

func crawl(input *url.URL) {
	fmt.Println(input.String())
	resp, err := http.Get(input.String())
	if err != nil {
		fmt.Println("GET request could not be processed correctly", err)
		return
	}
	base, err := url.Parse(input.String())
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("BODY could not be obtained from GET request", err)
		return
	}

	doc, err := html.Parse(strings.NewReader(string(body)))
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" && len(a.Val) > 0 {
					url, err := url.ParseRequestURI(a.Val)
					if err != nil {
						log.Fatal(err)
					}
					URL := base.ResolveReference(url).String()
					if url.Scheme != "tel" && strings.Compare(URL, input.String()) != 0 {
						fmt.Println(" " + URL)
					}
					break
				}

			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	//http://www.rescale.com/
}

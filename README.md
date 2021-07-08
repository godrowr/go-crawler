# go-crawler
My first project in Go. A simple web crawler that takes in a url input arguement. 
```
go run go-crawler.go http://www.rescale.com
```
The script then visits and outputs all the links <a href> visible on that webpage. 
Only visits links with the same domain as the url input arguement to avoid an infinite web crawler scenario. 
  
## Implementations
  So originally the go-crawler was built using standard go libraries only, implementing
  the 4 main specs through several methods. The 4 main specifications in question are:
  
  - Fetch the HTML document from the input URL
  - Parse out URLs in the HTML document
  - Log/print the URl visited along with all the URLs on the page
  - Loop back to step 1 for each of these new URLS
  
The old solution 	
#### Fetch the HTML document from the input URL

```go
func crawl_site(input string) {
	fmt.Println(input)
	links[input] = "Visited"
	resp, err := http.Get(input)
	if err != nil {
		fmt.Println("GET request could not be processed correctly", err)
		return
	}
	...
	
	search_links(doc, input)
}
```
  
#### Parse out URLs in the HTML document, Log/print the URl visited
  ```go
  func search_links(n *html.Node, input string) {
	...
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
	...
}
  ```
  
#### Loop back to step 1 for each of these new URLS
  ```go
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
	...
	os.Exit(0)
}
  ```

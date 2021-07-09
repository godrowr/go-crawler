# go-crawler
My first project in Go. A simple web crawler that takes in a url input arguement. 
```
go run go-crawler.go http://www.rescale.com
```
The script then visits and outputs all the links 'a href' visible on that webpage. 
Only visits links with the same domain as the url input arguement to avoid an infinite web crawler scenario. 

Though there is the improved-solution and the old-solution, main is currently the improved-solution, until I improve it even more. 

## Tests
```
git checkout old-solution
go run go-crawler.go http://www.rescale.com >> old-solution.txt
```
```
git checkout improved-solution
go run go-crawler.go http://www.rescale.com >> improved-solution.txt
```
  
## Specifications
  - Fetch the HTML document from the input URL
  - Parse out URLs in the HTML document
  - Log/print the URl visited along with all the URLs on the page
  - Loop back to step 1 for each of these new URLS

## Implementation
Typically I start out with some sort of vague idea in the steps of building the tool. Usually I try to do anything to get it to work and most of this code is poorly written and not suitable for production. Once it works, I restart the development process again, this time to build the tool optimally. 
  
  
### The Old solution 	
  So originally the go-crawler was built using standard go libraries only, implementing
  the 4 main specs through several methods. It used a map to store all the links and record if they were visited or not. It's time complexity was O(n^3) as it had a double recursive function call and a for loop. Very slow but this was my first web scraper in Go so I didn't expect too much.

  
### The Improved solution
Uses the open source module 'colly' which is build for Go web scapers. The performance improved significantly and the amount of code was vastly reduced. 

```go
c := colly.NewCollector(
	colly.AllowedDomains(host[1]),
	colly.Async(true),
)
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
```



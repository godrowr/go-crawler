package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
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
	fmt.Printf(": err=%+v url=%+v\n", err, u)
	if err != nil {
		fmt.Println("String inputted was not a URL", err)
		return
	}
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
	fmt.Println(string(body))
	//http://www.rescale.com
}

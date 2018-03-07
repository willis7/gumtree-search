// gumtree searcher
//
// usage:
// 		appname "<location>" "<search-category>" "<search-query>"...
//
package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

var baseURL = "http://gumtree.co.uk"
var location string
var category string
var search []string

// Helper function to pull the href attribute from a Token
func getAttr(t html.Token, attr string) (ok bool, val string) {
	// Iterate over all of the Token's attributes until we find an "href"
	for _, a := range t.Attr {
		if a.Key == attr {
			val = a.Val
			ok = true
		}
	}

	// "bare" return will return the variables (ok, href) as defined in
	// the function definition
	return
}

func crawl(query string, products chan string, done chan bool) {
	defer func() {
		//  notify that we're done after this function
		done <- true
	}()
	client := &http.Client{}
	req, _ := http.NewRequest("GET", baseURL+"/search", nil)
	q := req.URL.Query()
	q.Add("search_category", category)
	q.Add("search_location", location)
	q.Add("q", query)
	req.URL.RawQuery = q.Encode()

	log.Println(req.URL.String())
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("search failed: %v", err)
	}
	defer resp.Body.Close()

	// write a list of "hits" to a file "hits.txt"
	z := html.NewTokenizer(resp.Body)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			return
		case tt == html.StartTagToken:
			t := z.Token()

			// Check if the token is an <a> tag
			isAnchor := t.Data == "a"
			if !isAnchor {
				continue
			}

			ok, class := getAttr(t, "class")
			if class != "listing-link" {
				continue
			}

			// Extract the href value, if there is one
			ok, uri := getAttr(t, "href")
			if !ok {
				continue
			}

			// Make sure the uri begines in /p/
			isProduct := strings.Index(uri, "/p/") == 0
			if isProduct {
				products <- uri
			}

		}
	}
}

// newItems returns a slice of items which has not been seen before.
// This is similar to a SQL Right Join
func newItems(old []string, new []string) []string {
	var joinSlice []string
	oldMap := make(map[string]interface{}, len(old))

	for _, url := range old {
		oldMap[url] = nil
	}

	for _, url := range new {
		_, ok := oldMap[url]
		if ok {
			continue
		}
		joinSlice = append(joinSlice, url)
	}

	return joinSlice
}

// writeOutput takes a file handle and a map of lines and writes them
// to a file.
func writeOutput(f *os.File, lines []string) {
	for _, item := range lines {
		n3, err := f.WriteString(baseURL + item + "\n")
		if err != nil {
			log.Println("failed to write")
		}
		log.Println(baseURL + item)
		log.Printf("wrote %d bytes", n3)
	}
	f.Sync()
}

func main() {

	argLen := len(os.Args)
	if argLen < 4 {
		log.Fatalln("missing args")
	} else {
		location = os.Args[1]
		category = os.Args[2]
		search = os.Args[3:]
	}

	products := make(chan string)
	done := make(chan bool)

	for _, query := range search {
		go crawl(query, products, done)
	}

	foundItems := []string{}
	for c := 0; c < len(search); {
		select {
		case url := <-products:
			foundItems = append(foundItems, url)
		case <-done:
			c++
		}
	}
	close(products)

	// output the results to a file
	f, err := os.Create("output.txt")
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	writeOutput(f, foundItems)
}

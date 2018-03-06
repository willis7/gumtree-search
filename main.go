// gumtree searcher
//
// usage:
// 		appname "<location>" "<category>" "<search-query>"
//
package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

var baseURL = "http://gumtree.co.uk/"

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

func main() {
	var location string
	var category string
	var search []string

	if len(os.Args) != 4 {
		log.Fatalln("Missing args")
	} else {
		location = os.Args[1]
		category = os.Args[2]
		search = os.Args[3:]
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", baseURL+"search", nil)
	q := req.URL.Query()
	q.Add("search_category", category)
	q.Add("search_location", location)
	q.Add("q", search[0])
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
				log.Println(baseURL + uri)
			}

		}
	}
}

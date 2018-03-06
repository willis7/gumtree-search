package main

import (
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

var baseUrl = "http://gumtree.co.uk/"

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
	// hit the gumtree service with a url and search query
	// example request: https://www.gumtree.com/search?search_category=all&q=le+creuset&search_location=Exmouth
	var category = "all"
	var search = "le creuset"
	var location = "London"

	client := &http.Client{}
	req, _ := http.NewRequest("GET", baseUrl+"search", nil)
	q := req.URL.Query()
	q.Add("search_category", category)
	q.Add("search_location", location)
	q.Add("q", search)
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
			ok, url := getAttr(t, "href")
			if !ok {
				continue
			}

			// Make sure the url begines in /p/
			isProduct := strings.Index(url, "/p/") == 0
			if isProduct {
				log.Println(baseUrl + url)
			}

		}
	}
}

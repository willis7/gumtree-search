# Gumtree searcher

Takes search terms, and location and returns a list of URLs for that given search.

## Usage


```
$ go run main.go "<location>" "<search-category>" "<search-query>"...
e.g. go run main.go "London" "all" "some cool item1" "some cool item2"...
```

### Limitations

The application doesn't support multiple search categories. Therefore, it's advised to use `all` unless you're confident all of the search queries belong to the same supplied search category.

## TODO

* [x] os args
* [x] multiple searches
* [x] concurrent searches
* [ ] write output to file
* [ ] check for existing output file and ignore duplicates
* [ ] remove expired
* [ ] add time since post
* [ ] add email/sms support
* [ ] deploy serverless
* [ ] docker support
* [ ] makefile
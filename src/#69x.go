/* #69x.go
Exercise: Web Crawler

In this exercise you'll use Go's concurrency features to parallelize a web crawler.

Modify the Crawl function to fetch URLs in parallel without fetching the same URL 
twice.
================================= my summary =================================
An overview of existing pgm, in great detail since I'm new to golang
- main does 'Crawl("http://golang.org/", 4, fetcher)'
 - Crawl sig: func Crawl(url string, depth int, fetcher Fetcher) {
  - note that Crawl recursively calls itself with decremented depth; at each level calls:
    body, urls, err := fetcher.Fetch(url)
   - Fetch is a method of fetcher, see next 2 lines...
 - fetcher is a func of type Fetcher
  - type defn: type Fetcher interface { Fetch(url string) (body string, urls []string, err error)}
  - fetcher defn, first part: var fetcher = &fakeFetcher{"http://golang.org/": &fakeResult{
   - fakeFetcher defn: type fakeFetcher map[string]*fakeResult
   - fakeResult defn: type fakeResult struct {  body string urls     []string}
   - var fetcher defn also has more strings to simulate page content
Note that 'fake' part might be to allow demo of failure, so some can report 'not found' (in Fetch if 'ok' false, vs
'found' in Crawl after Fetch rtns if err is nil)
==============================================================================
Now I might understand well enough to begin to plan the 'TODO' items
- 'Fetch URLs in parallel': probably submit 
  'body, urls, err := fetcher.Fetch(url)' to a 'go' (maybe add channel too)
- 'Don't fetch the same URL twice': create a map, skip URLs alrdy in map
==============================================================================
== Sat  2012.04.14
Copied code as-is to goclipse-enabled Eclipse (from LiteIDE) [//` original cmts, //' mine]
trivial: chgd 'covered' to 'old'
realized map stg w 'u' is too soon, need to do only at 'url' (next recursion level): seems correct
//  */
package main

import (
    "fmt"
)

type Fetcher interface {
    //` Fetch returns the body of URL and
    //` a slice of URLs found on that page.
    Fetch(url string) (body string, urls []string, err error)
}

//` Crawl uses fetcher to recursively crawl
//` pages starting with url, to a maximum of depth.
//func Crawl(url string, depth int, fetcher Fetcher) {
func Crawl(url string, depth int, fetcher Fetcher, old map[string] bool) {
    //` TODO: Fetch URLs in parallel.
    //` TODO: Don't fetch the same URL twice.
    //` This implementation doesn't do either:
    
    //' for uniqueness
    //old := map[string] bool {  }

    if depth <= 0 {
        return
    }
    //fmt.Println("<diag> url: ", url, "; old: ", old[url], "; at depth: ", depth)
    if !old[url] { //'skips if already covered
        old[url] = true
        body, urls, err := fetcher.Fetch(url)
        //fmt.Println("<diag> err: [", err, "]; body: [", body, "]; urls: ", urls)
        if err != nil {
            fmt.Println(err)
            return
        }
        fmt.Printf("found: %s %q\n", url, body)
        for _, u := range urls {
            ////fmt.Println("<diag> u: %v; old: %v; at depth: %v", u, old[u], depth)
            //fmt.Println("<diag> u: ", u, "; old: ", old[u], "; at depth: ", depth)
            if !old[u] { //' skips if already covered
                //old[u] = true //' this might be too soon, try handle at next recursion level
                //Crawl(u, depth-1, fetcher)
                Crawl(u, depth-1, fetcher, old)
            }
        }
    }
    return
}

func main() {
    //' for uniqueness; must be added as param of Crawl
    old := map[string] bool {  }

    //Crawl("http://golang.org/", 4, fetcher)
    Crawl("http://golang.org/", 5, fetcher, old)
}


//` fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
    body string
    urls     []string
}

func (f *fakeFetcher) Fetch(url string) (string, []string, error) {
    if res, ok := (*f)[url]; ok {
        return res.body, res.urls, nil
    }
    return "", nil, fmt.Errorf("not found: %s", url)
}

//` fetcher is a populated fakeFetcher.
var fetcher = &fakeFetcher{
    "http://golang.org/": &fakeResult{
        "The Go Programming Language",
        []string{
            "http://golang.org/pkg/",
            "http://golang.org/cmd/",
        },
    },
    "http://golang.org/pkg/": &fakeResult{
        "Packages",
        []string{
            "http://golang.org/",
            "http://golang.org/cmd/",
            "http://golang.org/pkg/fmt/",
            "http://golang.org/pkg/os/",
        },
    },
    "http://golang.org/pkg/fmt/": &fakeResult{
        "Package fmt",
        []string{
            "http://golang.org/",
            "http://golang.org/pkg/",
        },
    },
    "http://golang.org/pkg/os/": &fakeResult{
        "Package os",
        []string{
            "http://golang.org/",
            "http://golang.org/pkg/",
        },
    },
}

package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/jackdanger/collectlinks"
	"net/http"
	"net/url"
	"os"
)

func main() {
	flag.Parse()

	args := flag.Args()
	fmt.Println(args)
	if len(args) < 1 {
		fmt.Println("Please specify start page")
		os.Exit(1)
	}

	queue := make(chan string)

	go func() { queue <- args[0] }()

	for uri := range queue {
		enqueue(uri, queue)
	}
}

func enqueue(uri string, queue chan string) {
	fmt.Println("fetching", uri)
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := http.Client{Transport: transport}
	resp, err := client.Get(uri)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	links := collectlinks.All(resp.Body)

	for _, link := range links {
		// absolute := fixUrl(link, uri) // Don't enqueue the raw thing we find,
		// fix it first.
		if uri != "" { // We'll set invalid URLs to blank strings
			// so let's never send those to the channel.
			go func() { queue <- link }()
		}
	}
}

func fixUrl(href, base string) string { // given a relative link and the page on
	uri, err := url.Parse(href) // which it's found we can parse them
	if err != nil {             // both and use the url package's
		return "" // ResolveReference function to figure
	} // out where the link really points.
	baseUrl, err := url.Parse(base) // If it's not a relative link this
	if err != nil {                 // is a no-op.
		return ""
	}
	uri = baseUrl.ResolveReference(uri)
	return uri.String() // We work with parsed url objects in this
}

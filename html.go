package main

import (
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// htmlBody is an HTML string
// rawBaseURL is the root URL of the website we're crawling. This will allow us to rewrite relative URLs into absolute URLs.
// It returns an un-normalized array of all the URLs found within the HTML, and an error if one occurs.
//
// You may find the reflect.DeepEqual function package to be particularly useful for testing.
//
// Test that relative URLs are converted to absolute URLs
// Test to be sure you find all the <a> tags in a body of HTML
// strings.NewReader(htmlBody) creates a io.Reader
// html.Parse(htmlReader) creates a tree of html.Nodes
// Use recursion to traverse the node tree and find the <a> tag "anchor" elements
func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {

	tree, err := html.Parse(strings.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}

	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, err
	}

	var urls []string

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {

					// Parse the URL
					u, err := url.Parse(a.Val)
					if err != nil {
						continue
					}

					// Resolve the URL
					u = baseURL.ResolveReference(u)

					// Add the URL to the list
					urls = append(urls, u.String())
				}
			}
		}
	}

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		f(n)
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(tree)

	return urls, nil
}

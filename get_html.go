package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func getHTML(rawURL string) (string, error) {
	// fetch the webpage with http.Get
	resp, err := http.Get(rawURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	// return error is status code is error-level
	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		return "", fmt.Errorf("status code %d", resp.StatusCode)
	}

	if !strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
		return "", fmt.Errorf("content type is %s", resp.Header.Get("Content-Type"))
	}

	// read the body of the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// return the body as a string
	return string(body), nil
}

func (cfg *config) crawlPage(rawCurrentURL string) {
	cfg.concurrencyControl <- struct{}{}
	defer func() {
		<-cfg.concurrencyControl
		cfg.wg.Done()
	}()
	if len(cfg.pages) >= cfg.maxPages {
		return
	}

	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Printf("error parsing currentURL: %s\n", err)
		return
	}

	if cfg.baseURL.Hostname() != currentURL.Hostname() {
		return
	}

	normalizedRawCurrentURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Printf("error normalizing currentURL: %s\n", err)
		return
	}

	if !cfg.addPageVisit(normalizedRawCurrentURL) {
		return
	}

	fmt.Printf("Crawling %s\n", rawCurrentURL)
	html, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("error getting HTML: %s\n", err)
		return
	}

	links, err := parseLinks(html, cfg.baseURL)
	if err != nil {
		fmt.Printf("error parsing links: %s\n", err)
		return
	}

	for _, link := range links {
		fmt.Printf("Found link: %s\n", link)
		cfg.wg.Add(1)
		go cfg.crawlPage(link)
	}
}

func parseLinks(raw string, baseURL *url.URL) ([]string, error) {
	var links []string

	doc, err := html.Parse(strings.NewReader(raw))
	if err != nil {
		return nil, err
	}

	for n := range doc.Descendants() {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					href, err := url.Parse(a.Val)
					if err != nil {
						fmt.Printf("error parsing href: %s\n", err)
						continue
					}
					resolvedURL := baseURL.ResolveReference(href)
					links = append(links, resolvedURL.String())
				}
			}
		}
	}

	return links, nil
}

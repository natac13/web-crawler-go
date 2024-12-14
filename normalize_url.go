package main

import (
	"net/url"
	"strings"
)

func normalizeURL(input string) (string, error) {
	parsedURL, err := url.Parse(input)
	if err != nil {
		return "", err
	}

	host := parsedURL.Host
	path := parsedURL.Path

	if strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}

	if strings.HasPrefix(host, "www.") {
		host = host[4:]
	}

	normalizedURL := host + path

	return normalizedURL, nil
}

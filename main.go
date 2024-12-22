package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("not enough arguments provided")
		os.Exit(1)
	}
	if len(os.Args) > 4 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	baseUrl := os.Args[1]
	fmt.Printf("starting crawl of: %s\n", baseUrl)
	mcs := os.Args[2]
	mps := os.Args[3]

	maxCurrency, err := strconv.Atoi(mcs)
	if err != nil {
		fmt.Printf("error converting maxConcurrency to int: %s\n", err)
		return
	}
	maxPages, err := strconv.Atoi(mps)
	if err != nil {
		fmt.Printf("error converting maxPages to int: %s\n", err)
		return
	}

	cfg, err := configure(baseUrl, maxCurrency, maxPages)
	if err != nil {
		fmt.Printf("error configuring: %s\n", err)
		return
	}
	cfg.wg.Add(1)
	go cfg.crawlPage(baseUrl)
	cfg.wg.Wait()

	// for page, count := range cfg.pages {
	// 	fmt.Printf("%s: %d\n", page, count)
	// }
	printReport(cfg.pages, baseUrl)
}

type ReportPage struct {
	page  string
	count int
}

func printReport(pages map[string]int, baseURL string) {
	fmt.Println("=============================")
	fmt.Printf("REPORT for %s\n", baseURL)
	fmt.Println("=============================")

	temp := make([]ReportPage, len(pages))
	for page, count := range pages {
		temp = append(temp, ReportPage{page: page, count: count})
	}

	slices.SortFunc(temp, func(a, b ReportPage) int {
		return a.count - b.count
	})

	for _, rp := range temp {
		if rp.count == 0 || rp.page == "" {
			continue
		}
		fmt.Printf("Fount %d internal links to %s\n", rp.count, rp.page)
	}
}

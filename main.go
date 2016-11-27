package main

import (
	"flag"
	"log"
	"time"

	"github.com/jordanpotter/site-analyzer/browser"
)

var (
	url              string
	chromedriverPath string
)

func init() {
	flag.StringVar(&url, "url", "", "url of the website")
	flag.StringVar(&chromedriverPath, "chromedriver", "/usr/bin/chromedriver", "path to chromedriver binary")
	flag.Parse()
}

func main() {
	verifyFlags()

	config := browser.Config{
		URL:                 url,
		ChromedriverPath:    chromedriverPath,
		PostNavigationSleep: 1 * time.Second,
	}
	analysis, err := browser.Analyze(config)
	if err != nil {
		log.Fatalln("Unexpected error while analyzing %q: %v", url, err)
	}

	log.Printf("Received %d console log entries", len(analysis.ConsoleLog))
	log.Printf("Received %d performance log entries", len(analysis.PerformanceLog))
}

func verifyFlags() {
	if url == "" {
		log.Fatalln("Must specify url")
	}
}

package main

import (
	"flag"
	"log"
	"time"

	"github.com/jordanpotter/site-analyzer/browser"
)

var (
	url              string
	chromeDriverPath string
)

func init() {
	flag.StringVar(&url, "url", "", "url of the website")
	flag.StringVar(&chromeDriverPath, "chromedriver", "/usr/bin/chromedriver", "path to chromedriver binary")
	flag.Parse()
}

func main() {
	verifyFlags()

	b, err := browser.New(chromeDriverPath, 0)
	if err != nil {
		log.Fatalln("Unexpected error while creating browser: %v", err)
	}

	analysis, err := b.Analyze(url, 1*time.Second)
	if err != nil {
		log.Fatalln("Unexpected error while analyzing %q: %v", url, err)
	}

	if err = b.Kill(); err != nil {
		log.Fatalln("Unexpected error while killing browser: %v", err)
	}

	log.Printf("Page took %f seconds to load", analysis.PageLoadTime.Seconds())
	log.Printf("Received %d console log entries", len(analysis.ConsoleLog))
	log.Printf("Received %d performance log entries", len(analysis.PerformanceLog))
}

func verifyFlags() {
	if url == "" {
		log.Fatalln("Must specify url")
	} else if chromeDriverPath == "" {
		log.Fatalln("Must specify chromedriver path")
	}
}

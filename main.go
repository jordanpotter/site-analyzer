package main

import (
	"flag"
	"log"
	"time"

	"github.com/jordanpotter/site-analyzer/browser"
	"github.com/jordanpotter/site-analyzer/video"
)

var (
	url              string
	chromeDriverPath string
	videoFPS         int
)

func init() {
	flag.StringVar(&url, "url", "", "url of the website")
	flag.StringVar(&chromeDriverPath, "chromedriver", "/usr/bin/chromedriver", "path to chromedriver binary")
	flag.IntVar(&videoFPS, "fps", 20, "fps of the captured video")
	flag.Parse()
}

func main() {
	verifyFlags()

	b, err := browser.New(chromeDriverPath, 0)
	if err != nil {
		log.Fatalln("Unexpected error while creating browser: %v", err)
	}

	capture, err := video.StartCapture(0, 1600, 1200, videoFPS, "output.mp4")
	if err != nil {
		log.Fatalln("Unexpected error while starting video capture: %v", err)
	}

	analysis, err := b.Analyze(url, 1*time.Second)
	if err != nil {
		log.Fatalln("Unexpected error while analyzing %q: %v", url, err)
	}

	if err = capture.Stop(); err != nil {
		log.Fatalln("Unexpected error while stopping video capture: %v", err)
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
	} else if videoFPS <= 0 {
		log.Fatalf("Invalid video fps %d", videoFPS)
	}
}

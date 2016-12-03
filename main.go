package main

import (
	"flag"
	"log"
	"time"

	"github.com/jordanpotter/site-analyzer/browser"
	"github.com/jordanpotter/site-analyzer/display"
	"github.com/jordanpotter/site-analyzer/video"
)

var (
	url              string
	width            int
	height           int
	fps              int
	videoDir         string
	chromeDriverPath string
)

func init() {
	flag.StringVar(&url, "url", "", "url of the website")
	flag.IntVar(&width, "width", 1600, "width of the captured video")
	flag.IntVar(&height, "height", 1200, "height of the captured video")
	flag.IntVar(&fps, "fps", 20, "fps of the captured video")
	flag.StringVar(&videoDir, "videodir", ".", "directory to save the captured video")
	flag.StringVar(&chromeDriverPath, "chromedriver", "/usr/bin/chromedriver", "path to chromedriver binary")
	flag.Parse()
}

func main() {
	verifyFlags()

	d, err := display.New(width, height)
	if err != nil {
		log.Fatalf("Unexpected error while creating display: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	b, err := browser.New(chromeDriverPath, d.Num)
	if err != nil {
		log.Fatalf("Unexpected error while creating browser: %v", err)
	}

	capture, err := video.StartCapture(d.Num, width, height, fps, videoDir)
	if err != nil {
		log.Fatalf("Unexpected error while starting video capture: %v", err)
	}

	analysis, err := b.Analyze(url, 1*time.Second)
	if err != nil {
		log.Fatalf("Unexpected error while analyzing %q: %v", url, err)
	}

	if err = capture.Stop(); err != nil {
		log.Fatalf("Unexpected error while stopping video capture: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	if err = b.Kill(); err != nil {
		log.Fatalf("Unexpected error while killing browser: %v", err)
	}

	if err = d.Kill(); err != nil {
		log.Fatalf("Unexpected error while killing display: %v", err)
	}

	log.Printf("Page took %f seconds to load", analysis.PageLoadTime.Seconds())
	log.Printf("Received %d console log entries", len(analysis.ConsoleLog))
	log.Printf("Received %d performance log entries", len(analysis.PerformanceLog))
}

func verifyFlags() {
	if url == "" {
		log.Fatalln("Must specify url")
	} else if width <= 0 {
		log.Fatalf("Invalid video width %d", width)
	} else if height <= 0 {
		log.Fatalf("Invalid video height %d", height)
	} else if fps <= 0 {
		log.Fatalf("Invalid video fps %d", fps)
	} else if videoDir == "" {
		log.Fatalln("Must specify video directory")
	} else if chromeDriverPath == "" {
		log.Fatalln("Must specify chromedriver path")
	}
}

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
	flag.IntVar(&fps, "fps", 30, "fps of the captured video")
	flag.StringVar(&videoDir, "videodir", ".", "directory to save the captured video")
	flag.StringVar(&chromeDriverPath, "chromedriver", "/usr/bin/chromedriver", "path to chromedriver binary")
	flag.Parse()
}

func main() {
	verifyFlags()

	log.Println("Creating display...")
	d, err := display.New(width, height)
	if err != nil {
		log.Fatalf("Unexpected error while creating display: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	log.Println("Opening browser...")
	b, err := browser.New(chromeDriverPath, width, height, d.Num)
	if err != nil {
		log.Fatalf("Unexpected error while creating browser: %v", err)
	}

	log.Println("Capturing video...")
	capture, err := video.StartCapture(d.Num, width, height, fps)
	if err != nil {
		log.Fatalf("Unexpected error while starting video capture: %v", err)
	}

	log.Printf("Analyzing %q...", url)
	analysis, err := b.Analyze(url, &browser.LoadedSpec{Operand: "and", Elements: []string{".thing"}}, 1*time.Second)
	if err != nil {
		log.Fatalf("Unexpected error while analyzing %q: %v", url, err)
	}

	log.Println("Stopping video capture...")
	if err = capture.Stop(); err != nil {
		log.Fatalf("Unexpected error while stopping video capture: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	log.Println("Closing browser...")
	if err = b.Kill(); err != nil {
		log.Fatalf("Unexpected error while killing browser: %v", err)
	}

	log.Println("Closing display...")
	if err = d.Kill(); err != nil {
		log.Fatalf("Unexpected error while killing display: %v", err)
	}

	log.Println("Outputting video capture...")
	if err = capture.Output(videoDir); err != nil {
		log.Fatalf("Unexpected error while outputting capture: %v", err)
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

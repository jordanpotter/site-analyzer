package main

import (
	"flag"
	"log"
	"time"

	"github.com/pkg/errors"

	"github.com/jordanpotter/site-analyzer/browser"
	"github.com/jordanpotter/site-analyzer/display"
	"github.com/jordanpotter/site-analyzer/utils"
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

	analysis, err := analyzeAndCapture()
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}

	log.Printf("Page took %f seconds to load", analysis.PageLoadTime.Seconds())
	log.Printf("Received %d console log entries", len(analysis.ConsoleLog))
	log.Printf("Received %d performance log entries", len(analysis.PerformanceLog))
}

func analyzeAndCapture() (*browser.Analysis, error) {
	log.Println("Creating display...")
	d, err := display.New(width, height)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create display")
	}
	defer utils.MustFunc(d.Close)

	log.Println("Creating Chrome browser...")
	b, err := browser.NewChrome(chromeDriverPath, width, height, d.Num)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create browser")
	}
	defer utils.MustFunc(b.Close)

	log.Println("Starting video capture...")
	c, err := video.StartCapture(d.Num, width, height, fps)
	if err != nil {
		return nil, errors.Wrap(err, "failed to start video capture")
	}
	defer utils.MustFunc(c.Stop)

	log.Printf("Analyzing %q...", url)
	analysis, err := b.Analyze(url, &browser.LoadedSpec{Operand: "and", Elements: []string{".thing"}}, 2*time.Second)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to analyze %q", url)
	}

	log.Println("Stopping video capture...")
	if err = c.Stop(); err != nil {
		return nil, errors.Wrap(err, "failed to stop video capture")
	}

	log.Println("Outputting video...")
	if err = c.Output(videoDir); err != nil {
		return nil, errors.Wrap(err, "failed to output video capture")
	}

	return analysis, nil
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

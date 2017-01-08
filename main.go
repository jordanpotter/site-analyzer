package main

import (
	"context"
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
	dataDir          string
	chromeDriverPath string
	deadline         string
)

func init() {
	flag.StringVar(&url, "url", "", "url of the website")
	flag.IntVar(&width, "width", 1600, "width of the captured video")
	flag.IntVar(&height, "height", 1200, "height of the captured video")
	flag.IntVar(&fps, "fps", 30, "fps of the captured video")
	flag.StringVar(&dataDir, "data", ".", "directory to save output")
	flag.StringVar(&deadline, "deadline", "30s", "cancel if have not completed within this duration")
	flag.StringVar(&chromeDriverPath, "chromedriver", "/usr/bin/chromedriver", "path to chromedriver binary")
	flag.Parse()
}

func main() {
	verifyFlags()

	timeout, err := time.ParseDuration(deadline)
	if err != nil {
		log.Fatalf("Unexpected error while parsing deadline: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	log.Printf("Analyzing %q...", url)
	analysis, capture, err := analyzeAndCapture(ctx)
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}

	log.Printf("Saving console logs...")
	consoleLogPath, err := analysis.ConsoleLog.Save(ctx, dataDir)
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}

	log.Println("Saving video...")
	videoPath, err := capture.SaveVideo(ctx, dataDir)
	if err != nil {
		log.Fatalf("Unexpected error while saving video: %v", err)
	}

	log.Println("Saving thumbnail...")
	thumbnailPath, err := capture.SaveThumbnail(ctx, analysis.PageLoadTime, dataDir)
	if err != nil {
		log.Fatalf("Unexpected error while saving thumbnail: %v", err)
	}

	log.Printf("Page took %f seconds to load", analysis.PageLoadTime.Seconds())
	log.Printf("Received %d console log entries", len(analysis.ConsoleLog.Entries))
	log.Printf("Received %d performance log entries", len(analysis.PerformanceLog))
	log.Printf("Console log saved to %s", consoleLogPath)
	log.Printf("Video saved to %s", videoPath)
	log.Printf("Thumbnail saved to %s", thumbnailPath)
}

func analyzeAndCapture(ctx context.Context) (*browser.Analysis, *video.Capture, error) {
	log.Println("Creating the display...")
	d, err := display.New(ctx, width, height)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create display")
	}
	defer utils.MustFunc(d.Close)

	log.Println("Opening Chrome...")
	b, err := browser.NewChrome(ctx, chromeDriverPath, width, height, d.Num, dataDir)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create browser")
	}
	defer utils.MustFunc(b.Close)

	log.Println("Starting video capture...")
	capture, err := video.StartCapture(ctx, d.Num, width, height, fps)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to start video capture")
	}
	defer utils.MustFunc(capture.Stop)

	log.Println("Performing analysis...")
	analysis, err := b.Analyze(ctx, url, &browser.LoadedSpec{}, 10*time.Second)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to analyze %q", url)
	}

	return analysis, capture, nil
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
	} else if dataDir == "" {
		log.Fatalln("Must specify data directory")
	} else if chromeDriverPath == "" {
		log.Fatalln("Must specify chromedriver path")
	}
}

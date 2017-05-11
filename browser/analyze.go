package browser

import (
	"context"
	"time"

	"github.com/pkg/errors"
)

type Analysis struct {
	PageLoadTime   time.Duration
	ConsoleLog     *ConsoleLog
	PerformanceLog *PerformanceLog
}

func (b *Browser) Analyze(ctx context.Context, url string, loadedSpec *LoadedSpec, postPageLoadSleep time.Duration) (*Analysis, error) {
	var analysis *Analysis
	var err error

	c := make(chan bool, 1)
	go func() {
		analysis, err = b.doAnalysis(url, loadedSpec, postPageLoadSleep)
		c <- true
	}()

	select {
	case <-c:
		return analysis, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (b *Browser) doAnalysis(url string, loadedSpec *LoadedSpec, postPageLoadSleep time.Duration) (*Analysis, error) {
	pageLoadTime, err := b.load(url, loadedSpec)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load %q", url)
	}

	time.Sleep(postPageLoadSleep)

	consoleLog, err := b.consoleLog()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get console log")
	}

	performanceLog, err := b.performanceLog()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get performance log")
	}

	return &Analysis{
		PageLoadTime:   pageLoadTime,
		ConsoleLog:     consoleLog,
		PerformanceLog: performanceLog,
	}, nil
}

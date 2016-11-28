package browser

import (
	"time"

	"github.com/pkg/errors"
)

type Analysis struct {
	PageLoadTime   time.Duration
	ConsoleLog     []ConsoleLogEntry
	PerformanceLog []PerformanceLogEntry
}

func (b *Browser) Analyze(url string, postPageLoadSleep time.Duration) (*Analysis, error) {
	start := time.Now()

	if err := b.session.Url(url); err != nil {
		return nil, errors.Wrapf(err, "failed to navigate to %q", url)
	}

	pageLoadTime := time.Since(start)

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

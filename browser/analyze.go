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

func (b *Browser) Analyze(url string, loadedSpec *LoadedSpec, postPageLoadSleep time.Duration) (*Analysis, error) {
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

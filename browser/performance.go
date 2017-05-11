package browser

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fedesog/webdriver"
	"github.com/jordanpotter/site-analyzer/utils"
	"github.com/pkg/errors"
)

const (
	performanceLogName     = "performance"
	performanceLogFilename = "performance.log"
)

type PerformanceLog struct {
	Entries []PerformanceLogEntry
}

type PerformanceLogEntry struct {
	Level   string
	Message string
	Time    time.Time
}

func (b *Browser) performanceLog() (*PerformanceLog, error) {
	logEntries, err := b.session.Log(performanceLogName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve performance log")
	}

	performanceLogEntries := make([]PerformanceLogEntry, 0, len(logEntries))
	for _, logEntry := range logEntries {
		performanceLogEntry := performanceLogEntry(logEntry)
		performanceLogEntries = append(performanceLogEntries, performanceLogEntry)
	}
	return &PerformanceLog{performanceLogEntries}, nil
}

func performanceLogEntry(logEntry webdriver.LogEntry) PerformanceLogEntry {
	return PerformanceLogEntry{
		Level:   logEntry.Level,
		Message: logEntry.Message,
		Time:    time.Unix(int64(logEntry.TimeStamp/1000), 0),
	}
}

func (pl *PerformanceLog) Save(ctx context.Context, dir string) (string, error) {
	var path string
	var err error

	c := make(chan bool, 1)
	go func() {
		path, err = pl.doSave(dir)
		c <- true
	}()

	select {
	case <-c:
		return path, err
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func (pl *PerformanceLog) doSave(dir string) (string, error) {
	path := filepath.Join(dir, performanceLogFilename)
	f, err := os.Create(path)
	if err != nil {
		return "", errors.Wrapf(err, "failed to create file %s", path)
	}
	defer utils.MustFunc(f.Close)

	for _, entry := range pl.Entries {
		time := entry.Time.Format(time.RFC3339)
		level := strings.ToUpper(entry.Level)
		str := fmt.Sprintf("%s %-7s %s\n", time, level, entry.Message)
		_, err = f.WriteString(str)
		if err != nil {
			return "", errors.Wrapf(err, "failed to write string to file %s", path)
		}
	}

	return path, nil
}

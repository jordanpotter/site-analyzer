package browser

import (
	"time"

	"github.com/fedesog/webdriver"
	"github.com/pkg/errors"
)

const performanceLogName = "performance"

type PerformanceLogEntry struct {
	Level   string
	Message string
	Time    time.Time
}

func performanceLog(session *webdriver.Session) ([]PerformanceLogEntry, error) {
	logEntries, err := session.Log(performanceLogName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve performance log")
	}

	performanceLogEntries := make([]PerformanceLogEntry, 0, len(logEntries))
	for _, logEntry := range logEntries {
		performanceLogEntry, err := performanceLogEntry(logEntry)
		if err != nil {
			return nil, errors.Wrap(err, "failed to process log entry")
		}
		performanceLogEntries = append(performanceLogEntries, performanceLogEntry)
	}
	return performanceLogEntries, nil
}

func performanceLogEntry(logEntry webdriver.LogEntry) (PerformanceLogEntry, error) {
	return PerformanceLogEntry{
		Level:   logEntry.Level,
		Message: logEntry.Message,
		Time:    time.Unix(int64(logEntry.TimeStamp/1000), 0),
	}, nil
}

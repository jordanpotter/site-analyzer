package browser

import (
	"time"

	"github.com/fedesog/webdriver"
	"github.com/pkg/errors"
)

const consoleLogName = "browser"

type ConsoleLogEntry struct {
	Level   string
	Message string
	Time    time.Time
}

func consoleLog(session *webdriver.Session) ([]ConsoleLogEntry, error) {
	logEntries, err := session.Log(consoleLogName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve console log")
	}

	consoleLogEntries := make([]ConsoleLogEntry, 0, len(logEntries))
	for _, logEntry := range logEntries {
		consoleLogEntry, err := consoleLogEntry(logEntry)
		if err != nil {
			return nil, errors.Wrap(err, "failed to process log entry")
		}
		consoleLogEntries = append(consoleLogEntries, consoleLogEntry)
	}
	return consoleLogEntries, nil
}

func consoleLogEntry(logEntry webdriver.LogEntry) (ConsoleLogEntry, error) {
	return ConsoleLogEntry{
		Level:   logEntry.Level,
		Message: logEntry.Message,
		Time:    time.Unix(int64(logEntry.TimeStamp/1000), 0),
	}, nil
}

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
	consoleLogName     = "browser"
	consoleLogFilename = "console.log"
)

type ConsoleLog struct {
	Entries []ConsoleLogEntry
}

type ConsoleLogEntry struct {
	Level   string
	Message string
	Time    time.Time
}

func (b *Browser) consoleLog() (*ConsoleLog, error) {
	logEntries, err := b.session.Log(consoleLogName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve console log")
	}

	consoleLogEntries := make([]ConsoleLogEntry, 0, len(logEntries))
	for _, logEntry := range logEntries {
		consoleLogEntry := consoleLogEntry(logEntry)
		consoleLogEntries = append(consoleLogEntries, consoleLogEntry)
	}
	return &ConsoleLog{consoleLogEntries}, nil
}

func consoleLogEntry(logEntry webdriver.LogEntry) ConsoleLogEntry {
	return ConsoleLogEntry{
		Level:   logEntry.Level,
		Message: logEntry.Message,
		Time:    time.Unix(int64(logEntry.TimeStamp/1000), 0),
	}
}

func (cl *ConsoleLog) Save(ctx context.Context, dir string) (string, error) {
	var path string
	var err error

	c := make(chan bool, 1)
	go func() {
		path, err = cl.doSave(dir)
		c <- true
	}()

	select {
	case <-c:
		return path, err
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func (cl *ConsoleLog) doSave(dir string) (string, error) {
	path := filepath.Join(dir, consoleLogFilename)
	f, err := os.Create(path)
	if err != nil {
		return "", errors.Wrapf(err, "failed to create file %s", path)
	}
	defer utils.MustFunc(f.Close)

	for _, entry := range cl.Entries {
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

package browser

import (
	"time"

	"github.com/fedesog/webdriver"
	"github.com/pkg/errors"
)

type Config struct {
	URL                 string
	ChromedriverPath    string
	PostNavigationSleep time.Duration
}

type Analysis struct {
	ConsoleLog     []ConsoleLogEntry
	PerformanceLog []PerformanceLogEntry
}

func Analyze(config Config) (*Analysis, error) {
	chromeDriver := webdriver.NewChromeDriver(config.ChromedriverPath)

	if err := chromeDriver.Start(); err != nil {
		return nil, errors.Wrap(err, "failed to start chromedriver")
	}

	session, err := chromeDriver.NewSession(desiredCapabilities(), requiredCapabilities())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create a new session")
	}

	if err = session.Url(config.URL); err != nil {
		return nil, errors.Wrapf(err, "failed to navigate to %q", config.URL)
	}

	time.Sleep(config.PostNavigationSleep)

	consoleLog, err := consoleLog(session)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get console log")
	}

	performanceLog, err := performanceLog(session)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get performance log")
	}

	if err = session.Delete(); err != nil {
		return nil, errors.Wrap(err, "failed to delete session")
	}

	if err = chromeDriver.Stop(); err != nil {
		return nil, errors.Wrap(err, "failed to stop chromedriver")
	}

	return &Analysis{
		ConsoleLog:     consoleLog,
		PerformanceLog: performanceLog,
	}, nil
}

func desiredCapabilities() webdriver.Capabilities {
	return webdriver.Capabilities{
		"loggingPrefs": map[string]interface{}{
			consoleLogName:     webdriver.LogAll,
			performanceLogName: webdriver.LogAll,
		},
		"chromeOptions": map[string]interface{}{
			"perfLoggingPrefs": map[string]interface{}{
				"enableNetwork": true,
				"enablePage":    true,
			},
		},
	}
}

func requiredCapabilities() webdriver.Capabilities {
	return webdriver.Capabilities{}
}

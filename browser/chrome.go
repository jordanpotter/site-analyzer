package browser

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/fedesog/webdriver"
	"github.com/pkg/errors"
)

const (
	chromedriverLogName       = "chromedriver.log"
	chromedriverOutputLogName = "chromedriver_output.log"
)

func NewChrome(ctx context.Context, chromeDriverPath string, width, height, displayNum int, logsDir string) (*Browser, error) {
	chromeDriver := webdriver.NewChromeDriver(chromeDriverPath)
	chromeDriver.LogPath = filepath.Join(logsDir, chromedriverLogName)
	chromeDriver.LogFile = filepath.Join(logsDir, chromedriverOutputLogName)

	if deadline, ok := ctx.Deadline(); ok {
		if deadline.Before(time.Now()) {
			return nil, errors.New("context deadline exceeded")
		}
		chromeDriver.StartTimeout = deadline.Sub(time.Now())
	}

	if err := chromeDriver.Start(); err != nil {
		return nil, errors.Wrap(err, "failed to start chromedriver")
	}

	session, err := chromeDriver.NewSession(chromeDesiredCapabilities(displayNum), chromeRequiredCapabilities())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create a new session")
	}

	if err = session.SetTimeoutsImplicitWait(implicitWaitTimeoutMs); err != nil {
		return nil, errors.Wrap(err, "failed to set implicit wait timeout")
	}

	if err = session.SetTimeoutsAsyncScript(asyncScriptTimeoutMs); err != nil {
		return nil, errors.Wrap(err, "failed to set async script timeout")
	}

	window, err := session.WindowHandle()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get window handle")
	}

	if err := window.SetPosition(webdriver.Position{X: 0, Y: 0}); err != nil {
		return nil, errors.Wrap(err, "failed to set window position")
	}

	if err := window.SetSize(webdriver.Size{Width: width, Height: height}); err != nil {
		return nil, errors.Wrap(err, "failed to set window size")
	}

	return &Browser{chromeDriver, session}, nil
}

func chromeDesiredCapabilities(displayNum int) webdriver.Capabilities {
	return webdriver.Capabilities{
		"pageLoadStrategy": "none",
		"loggingPrefs": map[string]interface{}{
			consoleLogName:     webdriver.LogAll,
			performanceLogName: webdriver.LogAll,
		},
		"chromeOptions": map[string]interface{}{
			"args": []string{
				"no-sandbox",
				"disable-infobars",
				"start-maximized",
				fmt.Sprintf("display=:%d", displayNum),
			},
			"perfLoggingPrefs": map[string]interface{}{
				"enableNetwork": true,
				"enablePage":    true,
			},
		},
	}
}

func chromeRequiredCapabilities() webdriver.Capabilities {
	return webdriver.Capabilities{}
}

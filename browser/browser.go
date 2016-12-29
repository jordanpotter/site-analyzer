package browser

import (
	"fmt"

	"github.com/fedesog/webdriver"
	"github.com/pkg/errors"
)

const (
	asyncScriptTimeoutMs = 30000
)

type Browser struct {
	chromeDriver *webdriver.ChromeDriver
	session      *webdriver.Session
}

func New(chromeDriverPath string, width, height, displayNum int) (*Browser, error) {
	chromeDriver := webdriver.NewChromeDriver(chromeDriverPath)
	if err := chromeDriver.Start(); err != nil {
		return nil, errors.Wrap(err, "failed to start chromedriver")
	}

	session, err := chromeDriver.NewSession(desiredCapabilities(displayNum), requiredCapabilities())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create a new session")
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

func (b *Browser) Kill() error {
	if err := b.session.Delete(); err != nil {
		return errors.Wrap(err, "failed to delete b.session")
	}

	if err := b.chromeDriver.Stop(); err != nil {
		return errors.Wrap(err, "failed to stop chromedriver")
	}

	return nil
}

func desiredCapabilities(displayNum int) webdriver.Capabilities {
	return webdriver.Capabilities{
		"pageLoadStrategy": "none",
		"loggingPrefs": map[string]interface{}{
			consoleLogName:     webdriver.LogAll,
			performanceLogName: webdriver.LogAll,
		},
		"chromeOptions": map[string]interface{}{
			"args": []string{
				"no-sandbox",
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

func requiredCapabilities() webdriver.Capabilities {
	return webdriver.Capabilities{}
}

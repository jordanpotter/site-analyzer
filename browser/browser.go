package browser

import (
	"github.com/fedesog/webdriver"
	"github.com/pkg/errors"
)

const (
	implicitWaitTimeoutMs = 60000
	asyncScriptTimeoutMs  = 60000
)

type Browser struct {
	webDriver webdriver.WebDriver
	session   *webdriver.Session
}

func (b *Browser) Close() error {
	if err := b.session.Delete(); err != nil {
		return errors.Wrap(err, "failed to delete b.session")
	}

	if err := b.webDriver.Stop(); err != nil {
		return errors.Wrap(err, "failed to stop chromedriver")
	}

	return nil
}

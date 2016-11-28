package video

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/pkg/errors"
)

type Capture struct {
	cmd *exec.Cmd
}

func StartCapture(displayNum, width, height, fps int, path string) (*Capture, error) {
	cmd := exec.Command("ffmpeg",
		"-f", "x11grab",
		"-video_size", fmt.Sprintf("%dx%d", width, height),
		"-framerate", strconv.Itoa(fps),
		"-i", fmt.Sprintf(":%d.0", displayNum), path,
	)

	if err := cmd.Start(); err != nil {
		return nil, errors.Wrap(err, "failed to start process")
	}

	return &Capture{cmd}, nil
}

func (c *Capture) Stop() error {
	err := c.cmd.Process.Signal(os.Interrupt)
	return errors.Wrap(err, "failed to interrupt process")
}

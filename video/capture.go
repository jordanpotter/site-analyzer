package video

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/pkg/errors"
)

const (
	videoName = "output.mp4"
)

type Capture struct {
	cmd *exec.Cmd
}

func StartCapture(displayNum, width, height, fps int, dir string) (*Capture, error) {
	path := filepath.Join(dir, videoName)
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

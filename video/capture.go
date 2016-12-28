package video

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/jordanpotter/site-analyzer/utils"
	"github.com/pkg/errors"
)

const (
	captureName  = "capture.mp4"
	videoName    = "video.mp4"
	videoQuality = 18
)

type Capture struct {
	cmd         *exec.Cmd
	capturePath string
}

func StartCapture(displayNum, width, height, fps int) (*Capture, error) {
	dir, err := ioutil.TempDir("", "capture")
	if err != nil {
		return nil, err
	}

	path := filepath.Join(dir, captureName)
	cmd := exec.Command("ffmpeg",
		"-probesize", "20M", // Analyze 20MB to retrieve stream information
		"-video_size", fmt.Sprintf("%dx%d", width, height), // Video dimensions
		"-framerate", strconv.Itoa(fps), // Framerate
		"-f", "x11grab", // Process input from X11
		"-draw_mouse", "0", // Hide mouse
		"-i", fmt.Sprintf(":%d.0+nomouse", displayNum), // Display to capture
		"-c:v", "libx264", // Video codec
		"-preset", "ultrafast", "-crf", "0", // Video quality
		"-pix_fmt", "yuv420p", // Require yuva420p pixel format
		"-an",                  // Disable audio
		"-loglevel", "warning", // Log level
		path,
	)

	if err := utils.ProcessCmdOutput(cmd); err != nil {
		return nil, errors.Wrap(err, "failed to process command output")
	}

	if err := cmd.Start(); err != nil {
		return nil, errors.Wrap(err, "failed to start process")
	}

	return &Capture{cmd, path}, nil
}

func (c *Capture) Stop() error {
	err := c.cmd.Process.Signal(os.Interrupt)
	return errors.Wrap(err, "failed to interrupt process")
}

func (c *Capture) Output(dir string) error {
	path := filepath.Join(dir, videoName)
	cmd := exec.Command("ffmpeg",
		"-i", c.capturePath, // Input file
		"-c:v", "libx264", // Video codec
		"-preset", "veryslow", "-crf", strconv.Itoa(videoQuality), // Video quality
		"-an",                  // Disable audio
		"-loglevel", "warning", // Log level
		path,
	)

	if err := utils.ProcessCmdOutput(cmd); err != nil {
		return errors.Wrap(err, "failed to process command output")
	}

	err := cmd.Run()
	return errors.Wrap(err, "failed to run process")
}

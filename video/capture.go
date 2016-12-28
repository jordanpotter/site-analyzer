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
	args := captureArgs(displayNum, width, height, fps, path)
	cmd := exec.Command("ffmpeg", args...)

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
	args := outputArgs(c.capturePath, path)
	cmd := exec.Command("ffmpeg", args...)

	if err := utils.ProcessCmdOutput(cmd); err != nil {
		return errors.Wrap(err, "failed to process command output")
	}

	err := cmd.Run()
	return errors.Wrap(err, "failed to run process")
}

func captureArgs(displayNum, width, height, fps int, dst string) []string {
	// WARNING: the order of arguments is very delicate
	var args []string

	// Analyze 20MB to retrieve stream information
	args = append(args, "-probesize", "20M")

	// Source dimensions and framerate
	args = append(args, "-video_size", fmt.Sprintf("%dx%d", width, height), "-framerate", strconv.Itoa(fps))

	// Source from X11, hide mouse
	args = append(args, "-f", "x11grab", "-draw_mouse", "0", "-i", fmt.Sprintf(":%d.0", displayNum))

	// Output codec with parameters
	args = append(args, "-c:v", "libx264", "-preset", "ultrafast", "-crf", "0", "-pix_fmt", "yuv420p")

	// Disable audio
	args = append(args, "-an")

	// Log level
	args = append(args, "-loglevel", "warning")

	// Destination
	args = append(args, dst)

	return args
}

func outputArgs(src, dst string) []string {
	// WARNING: the order of arguments is very delicate
	var args []string

	// Source
	args = append(args, "-i", src)

	// Output codec with parameters
	args = append(args, "-c:v", "libx264", "-preset", "slow", "-crf", strconv.Itoa(videoQuality))

	// Disable audio
	args = append(args, "-an")

	// Log level
	args = append(args, "-loglevel", "warning")

	// Destination
	args = append(args, dst)

	return args
}

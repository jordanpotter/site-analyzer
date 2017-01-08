package video

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/jordanpotter/site-analyzer/utils"
	"github.com/pkg/errors"
)

const (
	captureName       = "capture.mp4"
	captureStopDelay  = 100 * time.Millisecond
	videoFilename     = "video.mp4"
	videoQuality      = 18
	thumbnailFilename = "thumbnail.png"
)

type Capture struct {
	cmd         *exec.Cmd
	capturePath string
}

func StartCapture(ctx context.Context, displayNum, width, height, fps int) (*Capture, error) {
	dir, err := ioutil.TempDir("", "capture")
	if err != nil {
		return nil, err
	}

	path := filepath.Join(dir, captureName)
	args := captureArgs(displayNum, width, height, fps, path)
	cmd := exec.CommandContext(ctx, "ffmpeg", args...)

	if err := utils.ProcessCmdOutput(cmd); err != nil {
		return nil, errors.Wrap(err, "failed to process command output")
	}

	if err := cmd.Start(); err != nil {
		return nil, errors.Wrap(err, "failed to start process")
	}

	return &Capture{cmd, path}, nil
}

func (c *Capture) Stop() error {
	if err := c.cmd.Process.Signal(os.Interrupt); err != nil {
		return errors.Wrap(err, "failed to interrupt process")
	}

	time.Sleep(captureStopDelay)
	return nil
}

func (c *Capture) SaveVideo(ctx context.Context, dir string) (string, error) {
	path := filepath.Join(dir, videoFilename)
	args := videoArgs(c.capturePath, path)
	cmd := exec.CommandContext(ctx, "ffmpeg", args...)

	if err := utils.ProcessCmdOutput(cmd); err != nil {
		return "", errors.Wrap(err, "failed to process command output")
	}

	err := cmd.Run()
	return path, errors.Wrap(err, "failed to run process")
}

func (c *Capture) SaveThumbnail(ctx context.Context, loc time.Duration, dir string) (string, error) {
	path := filepath.Join(dir, thumbnailFilename)
	args := thumbnailArgs(loc, c.capturePath, path)
	cmd := exec.CommandContext(ctx, "ffmpeg", args...)

	if err := utils.ProcessCmdOutput(cmd); err != nil {
		return "", errors.Wrap(err, "failed to process command output")
	}

	err := cmd.Run()
	return path, errors.Wrap(err, "failed to run process")
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

func videoArgs(src, dst string) []string {
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

func thumbnailArgs(loc time.Duration, src, dst string) []string {
	// WARNING: the order of arguments is very delicate
	var args []string

	// Source
	args = append(args, "-i", src)

	// Seek to location
	args = append(args, "-ss", strconv.FormatFloat(loc.Seconds(), 'f', -1, 64))

	// Only capture one image
	args = append(args, "-vframes", "1")

	// Log level
	args = append(args, "-loglevel", "warning")

	// Destination
	args = append(args, dst)

	return args

}

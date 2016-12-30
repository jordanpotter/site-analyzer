package display

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"

	"github.com/jordanpotter/site-analyzer/utils"
	"github.com/pkg/errors"
)

const (
	depth         = 24
	maxDisplayNum = 50000
	startDelay    = 100 * time.Millisecond
)

type Display struct {
	Num int
	cmd *exec.Cmd
}

func New(ctx context.Context, width, height int) (*Display, error) {
	displayNum := randomDisplayNum()
	args := displayArgs(displayNum, width, height)
	cmd := exec.CommandContext(ctx, "Xvfb", args...)

	if err := utils.ProcessCmdOutput(cmd); err != nil {
		return nil, errors.Wrap(err, "failed to process command output")
	}

	if err := cmd.Start(); err != nil {
		return nil, errors.Wrap(err, "failed to start process")
	}

	time.Sleep(startDelay)
	return &Display{displayNum, cmd}, nil
}

func (d *Display) Close() error {
	err := d.cmd.Process.Signal(os.Interrupt)
	return errors.Wrap(err, "failed to interrupt process")
}

func displayArgs(displayNum, width, height int) []string {
	var args []string

	// Display number
	args = append(args, fmt.Sprintf(":%d", displayNum))

	// Screen settings
	args = append(args, "-screen", "0", fmt.Sprintf("%dx%dx%d", width, height, depth))

	return args
}

func randomDisplayNum() int {
	return 1 + rand.Intn(maxDisplayNum)
}

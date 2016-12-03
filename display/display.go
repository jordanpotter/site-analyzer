package display

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

const (
	depth         = 24
	maxDisplayNum = 50000
)

type Display struct {
	Num int
	cmd *exec.Cmd
}

func New(width, height int) (*Display, error) {
	displayNum := randomDisplayNum()
	cmd := exec.Command("Xvfb",
		fmt.Sprintf(":%d", displayNum),
		"-screen", "0", fmt.Sprintf("%dx%dx%d", width, height, depth),
	)

	if err := cmd.Start(); err != nil {
		return nil, errors.Wrap(err, "failed to start process")
	}

	return &Display{displayNum, cmd}, nil
}

func (d *Display) Kill() error {
	err := d.cmd.Process.Signal(os.Interrupt)
	return errors.Wrap(err, "failed to interrupt process")
}

func randomDisplayNum() int {
	return 1 + rand.Intn(maxDisplayNum)
}

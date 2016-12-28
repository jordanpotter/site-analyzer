package utils

import (
	"bufio"
	"io"
	"log"
	"os/exec"

	"github.com/pkg/errors"
)

func ProcessCmdOutput(cmd *exec.Cmd) error {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return errors.Wrap(err, "failed to receive stdout pipe")
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return errors.Wrap(err, "failed to receive stderr pipe")
	}

	go displayOutput(stdout)
	go displayOutput(stderr)

	return nil
}

func displayOutput(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		log.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("unexpected error scanning output: %v", err)
	}
}

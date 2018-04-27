// +build !linux,!darwin

package cmd

import (
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
)

func execPlatform(command string, args []string, env []string) error {
	ecmd := exec.Command(command, args...)
	ecmd.Stdin = os.Stdin
	ecmd.Stdout = os.Stdout
	ecmd.Stderr = os.Stderr
	ecmd.Env = env

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan)

	if err := ecmd.Start(); err != nil {
		return errors.Wrap(err, "Failed to start command")
	}

	go func() {
		for {
			sig := <-sigChan
			ecmd.Process.Signal(sig)
		}
	}()

	if err := ecmd.Wait(); err != nil {
		ecmd.Process.Signal(os.Kill)
		return errors.Wrap(err, "Failed to wait for command termination")
	}

	waitStatus := ecmd.ProcessState.Sys().(syscall.WaitStatus)
	os.Exit(waitStatus.ExitStatus())
	return nil // unreachable but Go doesn't know about it
}

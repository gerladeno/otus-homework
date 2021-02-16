package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
)

func RunCmd(cmd []string, env Environment) (returnCode int) {
	for key, value := range env {
		_ = os.Unsetenv(key)
		if !value.NeedRemove {
			if err := os.Setenv(key, value.Value); err != nil {
				log.Fatalf("unable to set env %s with value %s", key, value.Value)
			}
		}
	}
	if len(cmd) == 0 {
		return 0
	}
	command := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	if err := command.Run(); err != nil {
		var ee *exec.ExitError
		var e *exec.Error
		if errors.As(err, &ee) {
			return ee.ExitCode()
		}
		if errors.As(err, &e) {
			return 127
		}
	}
	return 0
}

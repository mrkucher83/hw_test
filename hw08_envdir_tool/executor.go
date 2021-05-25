package main

import (
	"fmt"
	"os"
	"os/exec"
)

func RunCmd(cmd []string, env Environment) (returnCode int) {
	cmnd := exec.Command(cmd[0], cmd[1:]...)
	for key, value := range env {
		if value.NeedRemove {
			if err := os.Unsetenv(key); err != nil {
				return 1
			}
			delete(env, key)
		}
	}
	cmnd.Env = os.Environ()
	for key, value := range env {
		elem := fmt.Sprintf("%s=%s", key, value.Value)
		cmnd.Env = append(cmnd.Env, elem)
	}
	cmnd.Stdout = os.Stdout
	cmnd.Stderr = os.Stderr
	cmnd.Stdin = os.Stdin

	if err := cmnd.Run(); err != nil {
		return 1
	}
	return 0
}

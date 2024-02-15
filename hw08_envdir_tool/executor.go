package main

import (
	"fmt"
	"os"
	"os/exec"
)

func createEnv(env Environment) ([]string, error) {
	for key, value := range env {
		err := os.Unsetenv(key)
		if err != nil {
			return nil, err
		}
		if !value.NeedRemove {
			err = os.Setenv(key, value.Value)
			if err != nil {
				return nil, err
			}
		}
	}
	return os.Environ(), nil
}

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, environment Environment) (returnCode int) {
	if len(cmd) == 0 {
		return 0
	}

	env, err := createEnv(environment)
	if err != nil {
		fmt.Printf("error on env setup: %v\n", err)
		return -1
	}

	command := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Env = env

	err = command.Run()
	if err != nil {
		fmt.Println(err)
	}

	return command.ProcessState.ExitCode()
}

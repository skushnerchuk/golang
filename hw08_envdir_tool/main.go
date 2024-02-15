package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Println("Expected at least 2 arguments")
		return
	}

	dir := args[1]
	command := args[2:]

	envDir, err := ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(RunCmd(command, envDir))
}

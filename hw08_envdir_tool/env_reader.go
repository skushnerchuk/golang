package main

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

func isCorrectFile(fileInfo os.DirEntry) bool {
	return fileInfo.Type().IsRegular() && !strings.Contains(fileInfo.Name(), "=")
}

func clearValue(value []byte) string {
	value = bytes.ReplaceAll(value, []byte{0x00}, []byte("\n"))
	return strings.TrimRightFunc(string(value), func(r rune) bool {
		return r == '\t' || r == ' '
	})
}

func readLine(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	scanner := bufio.NewScanner(f)
	scanner.Scan()
	return scanner.Bytes(), nil
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	env := make(Environment)

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if !isCorrectFile(f) {
			continue
		}
		fullName := filepath.Join(dir, f.Name())
		value, readErr := readLine(fullName)
		if readErr == nil {
			info, _ := f.Info()
			env[f.Name()] = EnvValue{clearValue(value), info.Size() == 0}
		}
	}

	return env, nil
}

package main

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

func isCorrectFile(fileInfo os.FileInfo) bool {
	return !fileInfo.IsDir() &&
		fileInfo.Mode().IsRegular() &&
		!strings.Contains(fileInfo.Name(), "=")
}

func clearValue(value []byte) string {
	vStr := string(value)
	vStr = strings.TrimRightFunc(vStr, unicode.IsSpace)
	parts := strings.Split(vStr, "\n")
	if len(parts) > 0 {
		vBytes := []byte(parts[0])
		vBytes = bytes.ReplaceAll(vBytes, []byte{0x00}, []byte("\n"))
		return string(vBytes)
	}
	return ""
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

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err == nil && isCorrectFile(info) {
			value, readErr := readLine(path)
			if readErr == nil {
				name := filepath.Base(path)
				env[name] = EnvValue{clearValue(value), info.Size() == 0}
				return nil
			}
			return err
		}
		return err
	})
	if err != nil {
		return nil, err
	}

	return env, nil
}

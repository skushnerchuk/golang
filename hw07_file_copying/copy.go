package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/schollz/progressbar/v3" //nolint:depguard
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	info, err := os.Stat(fromPath)
	if err != nil {
		return err
	}

	if !info.Mode().IsRegular() {
		return fmt.Errorf("%w: non-regular file %v", ErrUnsupportedFile, fromPath)
	}

	size := info.Size()
	if offset > size {
		return fmt.Errorf("%w: offset %d, file size: %d", ErrOffsetExceedsFileSize, offset, size)
	}

	fileFrom, err := os.OpenFile(fromPath, os.O_RDONLY, 0o666)
	if err != nil {
		return err
	}
	defer func() { _ = fileFrom.Close() }()

	if limit == 0 {
		limit = size
	}
	if limit > (size - offset) {
		limit = size - offset
	}

	fileTo, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer func() { _ = fileTo.Close() }()

	if _, err = fileFrom.Seek(offset, io.SeekStart); err != nil {
		return err
	}

	bar := progressbar.DefaultBytes(limit, "copy")
	defer func() { _ = bar.Close() }()

	_, err = io.CopyN(io.MultiWriter(fileTo, bar), fileFrom, limit)
	if err != nil {
		_ = os.Remove(toPath)
		return err
	}
	return nil
}

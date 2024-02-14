package main

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require" //nolint:depguard
)

func TestCopy(t *testing.T) {
	in := "in.txt"
	out := "out.txt"
	data := []byte(
		"Lorem ipsum dolor sit amet, consectetur adipiscing elit. In dapibus nibh nec metus lobortis ullamcorper.",
	)

	t.Run("full copy", func(t *testing.T) {
		prepare(t, in, data)
		defer purge(in, out)

		err := Copy(in, out, 0, 0)

		require.NoError(t, err)
		assert(t, out, data)
	})

	t.Run("offset 10 limit 0", func(t *testing.T) {
		prepare(t, in, data)
		defer purge(in, out)

		err := Copy(in, out, 10, 0)

		require.NoError(t, err)
		assert(t, out, data[10:])
	})

	t.Run("offset 0 limit 10", func(t *testing.T) {
		prepare(t, in, data)
		defer purge(in, out)

		err := Copy(in, out, 0, 10)

		require.NoError(t, err)
		assert(t, out, data[:10])
	})

	t.Run("offset 10 limit 10", func(t *testing.T) {
		prepare(t, in, data)
		defer purge(in, out)

		err := Copy(in, out, 10, 10)

		require.NoError(t, err)
		assert(t, out, data[10:20])
	})

	t.Run("offset exceeds file size", func(t *testing.T) {
		prepare(t, in, data)
		defer purge(in, out)

		err := Copy(in, out, int64(len(data)+10), 0)
		require.True(
			t,
			errors.Is(err, ErrOffsetExceedsFileSize),
			"offset exceeds file size: offset 100, file size: 13",
		)
	})

	t.Run("offset with limit exceeds file size", func(t *testing.T) {
		prepare(t, in, data)
		defer purge(in, out)

		err := Copy(in, out, int64(len(data)-10), 20)
		require.NoError(t, err)
		assert(t, out, data[len(data)-10:])
	})

	t.Run("unsupported file as source", func(t *testing.T) {
		err := Copy("/dev/random", out, 0, 0)

		require.True(
			t,
			errors.Is(err, ErrUnsupportedFile),
			"unsupported file: non-regular file /dev/random",
		)
	})

	t.Run("directory as source", func(t *testing.T) {
		err := Copy("/tmp", out, 0, 0)

		require.True(
			t,
			errors.Is(err, ErrUnsupportedFile),
			"unsupported file: non-regular file /tmp",
		)
	})

	t.Run("source does not exist", func(t *testing.T) {
		err := Copy("/tmp/does-not-exist", out, 0, 0)

		require.Equal(t, "stat /tmp/does-not-exist: no such file or directory", err.Error())
	})
}

func prepare(t *testing.T, inFilePath string, data []byte) {
	t.Helper()
	err := os.WriteFile(inFilePath, data, os.ModePerm)
	require.NoError(t, err)
}

func purge(files ...string) {
	for _, file := range files {
		_ = os.Remove(file)
	}
}

func assert(t *testing.T, outFilePath string, expected []byte) {
	t.Helper()
	data, err := os.ReadFile(outFilePath)
	require.NoError(t, err)
	require.Equal(t, expected, data)
}

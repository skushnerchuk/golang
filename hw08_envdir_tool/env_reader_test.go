package main

import (
	//nolint:depguard
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestReadDir(t *testing.T) {

	t.Run("env values", func(t *testing.T) {
		valid := Environment{
			"BAR":   {Value: "bar", NeedRemove: false},
			"FOO":   {Value: "   foo\nwith new line", NeedRemove: false},
			"UNSET": {Value: "", NeedRemove: true},
			"HELLO": {Value: "\"hello\"", NeedRemove: false},
			"EMPTY": {Value: "", NeedRemove: false},
		}

		env, err := ReadDir("./testdata/env")

		require.NoError(t, err)
		require.Equal(t, len(env), len(valid))
		for k, v := range env {
			require.Equal(t, valid[k], v)
		}
	})

	t.Run("incorrect dir", func(t *testing.T) {
		env, err := ReadDir("./dir-not-exist")

		require.NotNil(t, err)
		require.Nil(t, env)
	})

	t.Run("incorrect file name", func(t *testing.T) {
		d, err := os.MkdirTemp("", "env_dir*")
		require.NoError(t, err)
		_, err = os.CreateTemp(d, "bad=file*")
		require.NoError(t, err)
		defer func() { _ = os.RemoveAll(d) }()
		env, err := ReadDir(d)
		require.Len(t, env, 0)
		require.Nil(t, err)
	})
}

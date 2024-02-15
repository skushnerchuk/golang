package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require" //nolint:depguard
)

func TestRunCmd(t *testing.T) {
	t.Run("check unset", func(t *testing.T) {
		err := os.Setenv("ENV_VAL", "value")
		require.NoError(t, err)

		env := make(Environment)
		env["ENV_VAL"] = EnvValue{"", true}
		retCode := RunCmd([]string{"ls"}, env)

		require.Equal(t, retCode, 0)
		_, ok := os.LookupEnv("ENV_VAL")
		require.Equal(t, ok, false)
	})

	t.Run("check set", func(t *testing.T) {
		env := make(Environment)
		env["ENV_VAL"] = EnvValue{"value", false}
		_ = RunCmd([]string{"ls"}, env)
		val, ok := os.LookupEnv("ENV_VAL")
		require.Equal(t, ok, true)
		require.Equal(t, val, "value")
	})

	t.Run("check change", func(t *testing.T) {
		err := os.Setenv("ENV_VAL", "value1")
		require.NoError(t, err)
		env := make(Environment)
		env["ENV_VAL"] = EnvValue{"value2", false}
		retCode := RunCmd([]string{"ls"}, env)
		require.Equal(t, retCode, 0)
		val, ok := os.LookupEnv("ENV_VAL")
		require.Equal(t, ok, true)
		require.Equal(t, val, "value2")
	})

	t.Run("check ret code", func(t *testing.T) {
		env := make(Environment)
		retCode := RunCmd([]string{"unknown_command"}, env)
		require.Equal(t, retCode, -1)
		retCode = RunCmd([]string{"ls"}, env)
		require.Equal(t, retCode, 0)
	})
}

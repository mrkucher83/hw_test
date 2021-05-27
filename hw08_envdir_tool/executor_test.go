package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	args := []string{"testdata/echo.sh", "arg1=1", "arg2=2"}
	envs, _ := ReadDir("testdata/env")

	t.Run("valid data", func(t *testing.T) {
		n := RunCmd(args, envs)

		require.Equal(t, 0, n)
	})

	t.Run("invalid directory", func(t *testing.T) {
		n := RunCmd([]string{"testdata111/echo.sh"}, envs)

		require.Equal(t, 1, n)
	})
}

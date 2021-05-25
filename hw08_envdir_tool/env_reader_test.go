package main

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	envDir := "testdata/env"

	t.Run("valid data", func(t *testing.T) {
		err := os.WriteFile("testdata/env/NEW=", []byte("new"), 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer os.Remove("testdata/env/NEW=")

		envs, err := ReadDir(envDir)

		require.NoError(t, err)
		require.Equal(t, "bar", envs["BAR"].Value)
		require.Equal(t, "", envs["EMPTY"].Value)
		require.Equal(t, "new", envs["NEW"].Value)
	})
}

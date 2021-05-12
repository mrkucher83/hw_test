package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	fromPath := "testdata/input.txt"
	toPath := "outFile.txt"

	t.Run("offset0_limit0", func(t *testing.T) {
		_ = Copy(fromPath, toPath, 0, 0)

		b, _ := ioutil.ReadFile(toPath)
		srcFile, _ := os.Open(fromPath)
		defer srcFile.Close()
		finfo, _ := srcFile.Stat()

		require.Equal(t, finfo.Size(), int64(len(b)))
	})

	t.Run("limit exceeds file size", func(t *testing.T) {
		_ = Copy(fromPath, toPath, 10000, 0)

		b, _ := ioutil.ReadFile(toPath)
		srcFile, _ := os.Open(fromPath)
		defer srcFile.Close()
		finfo, _ := srcFile.Stat()

		require.Equal(t, finfo.Size(), int64(len(b)))
	})

	t.Run("offset0_limit10", func(t *testing.T) {
		_ = Copy(fromPath, toPath, 10, 0)

		b, _ := ioutil.ReadFile(toPath)
		require.Equal(t, 10, len(b))
	})

	t.Run("offset3_limit9", func(t *testing.T) {
		_ = Copy(fromPath, toPath, 9, 3)

		b, _ := ioutil.ReadFile(toPath)
		require.EqualValues(t, "Documents", string(b))
		require.Equal(t, 9, len(b))
	})

	t.Run("offset exceeds file size", func(t *testing.T) {
		err := Copy(fromPath, toPath, 10, 7000)

		require.ErrorIs(t, err, ErrOffsetExceedsFileSize)
	})

	t.Run("unsupported file", func(t *testing.T) {
		err := Copy("/dev/urandom", toPath, 0, 0)

		require.ErrorIs(t, err, ErrUnsupportedFile)
	})
}

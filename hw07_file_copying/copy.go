package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrNegativeLimit         = errors.New("limit should be a zero or positive number")
)

func Copy(fromPath, toPath string, limit, offset int64) error {
	newLimit := limit

	srcFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	finfo, err := srcFile.Stat()
	if err != nil {
		fmt.Println("err from Stat(): ", err)
		return err
	}
	if !finfo.Mode().IsRegular() {
		return ErrUnsupportedFile
	}
	if offset > finfo.Size() {
		return ErrOffsetExceedsFileSize
	}
	if newLimit == 0 || newLimit+offset > finfo.Size() {
		newLimit = finfo.Size() - offset
	}
	if newLimit < 0 {
		return ErrNegativeLimit
	}

	newFile, err := os.Create(toPath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer newFile.Close()

	bar := pb.Full.Start64(newLimit)
	if _, err := srcFile.Seek(offset, 0); err != nil {
		return fmt.Errorf("file seek error: %w", err)
	}
	barReader := bar.NewProxyReader(srcFile)

	_, err = io.CopyN(newFile, barReader, newLimit)
	if err != nil {
		return err
	}
	bar.Finish()

	return nil
}

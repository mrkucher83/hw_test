package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrNegativeLimit         = errors.New("limit should be a zero or positive number")
)

func Copy(fromPath, toPath string, limit, offset int64) error {
	newLimit := limit

	srcFile, err := os.OpenFile(fromPath, os.O_RDWR|os.O_CREATE, 0755)
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
	if newLimit == 0 || newLimit > finfo.Size() {
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

	buf, err := readOffset(srcFile, offset, newLimit)
	if err != nil {
		return err
	}

	r := strings.NewReader(string(buf))

	_, err = io.CopyN(newFile, r, newLimit)
	if err != nil {
		return err
	}

	return nil
}

func readOffset(file *os.File, offset, limit int64) ([]byte, error) {
	buffer := make([]byte, limit+offset)
	n, err := file.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("file read error: %w", err)
	}
	return buffer[offset:n], nil
}

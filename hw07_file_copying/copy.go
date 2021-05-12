package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
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
	if finfo.Size() == 0 && finfo.Name() == "urandom" {
		return ErrUnsupportedFile
	}
	if offset > finfo.Size() {
		return ErrOffsetExceedsFileSize
	}
	if newLimit == 0 {
		newLimit = finfo.Size() - offset
	}

	tmpfile, err := os.CreateTemp("testdata", "tmpFile.txt")
	if err != nil {
		return err
	}
	defer os.Remove(tmpfile.Name())

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

	_, err = tmpfile.Write(buf)
	if err != nil {
		return err
	}

	tf, err := os.OpenFile(tmpfile.Name(), os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer tf.Close()

	_, err = io.CopyN(newFile, tf, newLimit)
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

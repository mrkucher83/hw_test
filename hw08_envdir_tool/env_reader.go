package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type Environment map[string]EnvValue

type EnvValue struct {
	Value      string
	NeedRemove bool
}

func ReadDir(dir string) (Environment, error) {
	envs := make(map[string]EnvValue)
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	for _, file := range files {
		// открываем файл
		f, err := os.Open(fmt.Sprintf("%s/%s", dir, file.Name()))
		if err != nil {
			return nil, fmt.Errorf("failed to open file: %w", err)
		}

		// создаем io.Reader и читаем первую строку
		reader := bufio.NewReader(f)
		line, err := reader.ReadBytes('\n')
		if errors.Is(err, io.EOF) {
			err = nil
		}
		if err != nil {
			return nil, fmt.Errorf("failed to ReadBytes: %w", err)
		}
		if err = f.Close(); err != nil {
			return nil, fmt.Errorf("failed to close file: %w", err)
		}

		line = bytes.ReplaceAll(line, []byte{0x00}, []byte("\n"))
		fileName := strings.ReplaceAll(file.Name(), "=", "")

		// удаляем лишние пробелы у строки
		str := bytes.TrimRight(line, "\n\t ")

		// создаем экземпляр EnvValue, заполняем значения и добавляем в map
		var envVal EnvValue
		envVal.Value = string(str)
		if len(str) == 0 {
			envVal.NeedRemove = true
		} else {
			envVal.NeedRemove = false
		}
		envs[fileName] = envVal
	}

	return envs, nil
}

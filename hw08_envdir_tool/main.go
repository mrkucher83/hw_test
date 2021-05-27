package main

import (
	"log"
	"os"
)

func main() {
	args := os.Args

	envs, err := ReadDir(args[1])
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	n := RunCmd(args[2:], envs)
	os.Exit(n)
}

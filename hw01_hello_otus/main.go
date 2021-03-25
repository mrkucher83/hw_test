package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	phrase := "Hello, OTUS!"

	fmt.Println(stringutil.Reverse(phrase))
}

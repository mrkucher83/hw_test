package hw02unpackstring

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(val string) (string, error) {
	var exp strings.Builder
	var temp rune

	for _, el := range val {
		if unicode.IsDigit(el) && temp == 0 {
			return "", ErrInvalidString
		}

		if !unicode.IsDigit(el) {
			temp = el
			exp.WriteRune(el)
		} else {
			count, err := strconv.Atoi(string(el))
			if err != nil {
				return "", fmt.Errorf("%w: %v", ErrInvalidString, err)
			}
			if count > 0 {
				exp.WriteString(strings.Repeat(string(temp), count-1))
				temp = 0
			} else {
				arr := []rune(exp.String())
				arr = arr[:len(arr)-1]
				temp = 0
				exp.Reset()
				exp.WriteString(string(arr))
			}
		}
	}

	return exp.String(), nil
}

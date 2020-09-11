package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	var result string

	for i, v := range []rune(s) {

		var g string
		if unicode.IsDigit(v) {

			b, _ := strconv.Atoi(string(v))

			if b == 0 {
				size := len(result)
				result = result[:size-1]
			} else {
				if i == 0 {
					return "", ErrInvalidString
				}
				_, err := strconv.Atoi(string([]rune(s)[i-1]))
				if err == nil {
					return "", ErrInvalidString
				}
				g = strings.Repeat(string([]rune(s)[i-1]), b-1)

				result += g
			}
		} else {
			result += string(v)
		}
	}
	return result, ErrInvalidString
}

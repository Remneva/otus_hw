package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	var result strings.Builder
	for i, v := range []rune(s) {
		if unicode.IsDigit(v) {
			var err error
			_, err = workWithDigit(i, s, v, &result)
			if err != nil {
				return "", err
			}
		} else {
			result.WriteRune(v)
		}
	}
	return result.String(), nil
}

func workWithDigit(i int, s string, v rune, result *strings.Builder) (*strings.Builder, error) {
	var g string
	b, _ := strconv.Atoi(string(v))
	if i == 0 {
		return nil, ErrInvalidString
	}
	_, err := strconv.Atoi(string(s[i-1]))
	if err == nil {
		return nil, ErrInvalidString
	}
	if b == 0 {
		size := result.Len()
		result.Reset()
		result.WriteString(s[:size-1])
	} else {
		g = strings.Repeat(string(s[i-1]), b-1)
		result.WriteString(g)
	}
	return result, nil
}

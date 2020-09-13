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

			if i == 0 { // проверка первого символа - первый символ не должен быть интом
				return "", ErrInvalidString
			}
			_, err := strconv.Atoi(string([]rune(s)[i-1])) // проверка предыдущего символа - не должно ыть два инта подряд

			if err == nil {
				return "", ErrInvalidString
			}
			if b == 0 { // если инт = 0, не повторять букву ни одного раза
				size := len(result)
				result = result[:size-1]

			} else {
				g = strings.Repeat(string([]rune(s)[i-1]), b-1)
				result += g
			}

		} else {
			result += string(v)
		}
	}
	return result, nil
}

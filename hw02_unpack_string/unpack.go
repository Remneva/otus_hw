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
		if unicode.IsDigit(v) {
			var err error
			result, err = workWithDigit(i, s, v, result)
			if err != nil {
				return "", err
			}
		} else {
			result += string(v)
		}
	}
	return result, nil
}

func workWithDigit(i int, s string, v rune, result string) (string, error) {
	b, _ := strconv.Atoi(string(v))
	var g string
	if i == 0 { // проверка первого символа - первый символ не должен быть интом
		return result, ErrInvalidString
	}
	_, err := strconv.Atoi(string([]rune(s)[i-1])) // проверка предыдущего символа - не должно быть два инта подряд
	if err == nil {
		return result, ErrInvalidString
	}
	if b == 0 { // если инт = 0, не повторять букву ни одного раза
		size := len(result)
		result = result[:size-1]
	} else {
		g = strings.Repeat(string([]rune(s)[i-1]), b-1) // повторить букву i раз
		result += g
	}
	return result, nil
}

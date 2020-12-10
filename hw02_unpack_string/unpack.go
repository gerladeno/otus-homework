package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func processTwoRunes(first, second *rune, builder *strings.Builder) (bool, error) {
	switch *first {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return false, ErrInvalidString
	case '\\':
		if second == nil {
			return false, ErrInvalidString
		}
		if *second == '\\' || unicode.IsDigit(*second) {
			builder.WriteRune(*second)
			return false, nil
		}
		return false, ErrInvalidString
	default:
		if second == nil {
			builder.WriteRune(*first)
			return false, nil
		}
		if unicode.IsDigit(*second) {
			i, _ := strconv.Atoi(string(*second))
			builder.WriteString(strings.Repeat(string(*first), i))
			return true, nil
		}

		builder.WriteRune(*first)
		return false, nil
	}
}

func Unpack(s string) (string, error) {
	result := new(strings.Builder)
	runes := []rune(s)
	for i := 0; i <= len(runes)-1; i++ {
		var next *rune
		switch i {
		case len(runes) - 1:
		case len(runes):
			break
		default:
			next = &runes[i+1]
		}
		skipSecond, err := processTwoRunes(&runes[i], next, result)
		if err != nil {
			return "", err
		}
		if skipSecond {
			i++
		}
	}
	return result.String(), nil
}

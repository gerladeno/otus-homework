package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

const ESCAPE = '\\'

type state int

const (
	unaffected state = iota
	escaped
	symbol
)

type holder struct {
	r rune
	s state
}

func processRune(r rune, h *holder, result *strings.Builder) error {
	switch {
	case h.s == unaffected:
		if unicode.IsDigit(r) {
			return ErrInvalidString
		}
		if r == ESCAPE {
			h.s = escaped
			return nil
		}
		h.r = r
		h.s = symbol
	case h.s == escaped:
		if !unicode.IsDigit(r) && r != ESCAPE {
			return ErrInvalidString
		}
		h.r = r
		h.s = symbol
		return nil
	case h.s == symbol:
		if unicode.IsDigit(r) {
			result.WriteString(strings.Repeat(string(h.r), int(r-'0')))
			h.s = unaffected
			return nil
		}
		if r == ESCAPE {
			result.WriteRune(h.r)
			h.r = ESCAPE
			h.s = escaped
			return nil
		}
		result.WriteRune(h.r)
		h.r = r
		h.s = symbol
	default:
		return fmt.Errorf("wtf error")
	}
	return nil
}

func Unpack(s string) (string, error) {
	result := &strings.Builder{}
	h := holder{}
	for _, symbol := range s {
		if err := processRune(symbol, &h, result); err != nil {
			return "", err
		}
	}
	if err := processRune('a', &h, result); err != nil {
		return "", err
	}
	return result.String(), nil
}

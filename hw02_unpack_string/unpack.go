package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"strconv"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

type packedString struct {
	raw      string
	bytes    []byte
	unpacked string
}

func newPackedString(s string) *packedString {
	return &packedString{
		raw:   s,
		bytes: []byte(s),
	}
}

func (p *packedString) isValid() error {
	if len(p.bytes) == 0 {
		return nil
	}
	if unicode.IsDigit(rune(p.bytes[0])) {
		return ErrInvalidString
	}
	for i := 1; i <= len(p.bytes)-1; i++ {
		if unicode.IsDigit(rune(p.bytes[i-1])) && unicode.IsDigit(rune(p.bytes[i])) {
			return ErrInvalidString
		}
	}
	return nil
}

func (p *packedString) unpack() error{
	tmp := make([]byte, 0)
	for i, v := range p.bytes {
		switch {
		case unicode.IsDigit(rune(v)):
		case i == len(p.bytes)-1:
			tmp = append(tmp, v)
		case !unicode.IsDigit(rune(v)) && !unicode.IsDigit(rune(p.bytes[i+1])):
			tmp = append(tmp, v)
		case !unicode.IsDigit(rune(v)) && unicode.IsDigit(rune(p.bytes[i+1])):
			n, err := strconv.Atoi(string(p.bytes[i+1]))
			if err != nil {
				return err
			}
			for j:=0; j < n; j++ {
				tmp = append(tmp, v)
			}
		}
	}
	p.unpacked = string(tmp)
	return nil
}

func Unpack(s string) (string, error) {
	str := newPackedString(s)
	if err := str.isValid(); err != nil {
		return str.unpacked, err
	}
	err := str.unpack()
	if err != nil {
		return str.unpacked, err
	}
	return str.unpacked, nil
}

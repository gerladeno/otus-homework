package hw02_unpack_string //nolint:golint,stylecheck

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type test struct {
	input    string
	expected string
	err      error
}

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		err      error
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "3abc", expected: "", err: ErrInvalidString},
		{input: "45", expected: "", err: ErrInvalidString},
		{input: "aaa10b", expected: "", err: ErrInvalidString},
		{input: "aaa10b", expected: "", err: ErrInvalidString},
		{input: "8", expected: "", err: ErrInvalidString},
		{input: "", expected: ""},
		{input: "aaa0b", expected: "aab"},
		{input: " 6", expected: "      "},
		{input: "a+3", expected: "a+++"},
		{input: "d\n4", expected: "d\n\n\n\n"},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.Equal(t, tc.err, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackWithEscape(t *testing.T) {
	t.Skip() // NeedRemove if task with asterisk completed

	tests := []struct {
		input    string
		expected string
		err      error
	}{
		{input: `qwe\4\5`, expected: `qwe45`},
		{input: `qwe\45`, expected: `qwe44444`},
		{input: `qwe\\5`, expected: `qwe\\\\\`},
		{input: `qwe\\\3`, expected: `qwe\3`},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.Equal(t, tc.err, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

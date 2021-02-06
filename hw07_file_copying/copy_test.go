package main

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestCopyErr(t *testing.T) {
	tests := []struct {
		name   string
		from   string
		to     string
		offset int64
		limit  int64
		err    error
	}{
		{"offset too big", "test.sh", "", 800, 0, ErrOffsetExceedsFileSize},
		{"file no length", "/dev/urandom", "", 0, 0, ErrUnsupportedFile},
		{"ok", "test.sh", "/tmp/out", 0, 0, nil},
		{"max offset", "test.sh", "/tmp/out", 734, 0, nil},
		{"limit exceeds size", "testdata/out_offset0_limit10000.txt", "/dev/null", 0, 0, nil},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := Copy(test.from, test.to, test.offset, test.limit)
			require.Equal(t, err, test.err)
		})
	}
}

func TestCopy(t *testing.T) {
	text := "string_of_twenty_nine_symbols"
	nonLatin := "ролыиукфщ пйукщнгп п9768ТЕ:%Р?ШН(*Н?О"
	binaryData := make([]byte, 100)
	fileIn := "tmp_in"
	fileOut := "tmp_out"
	tests := []struct {
		name   string
		offset int64
		limit  int64
		input  string
		result string
	}{
		{"equals", 0, 0, text, text},
		{"offset equals size", 29, 0, text, ""},
		{"offset 10", 10, 0, text, "twenty_nine_symbols"},
		{"limit 10", 0, 10, text, "string_of_"},
		{"limit and offset 10", 10, 10, text, "twenty_nin"},
		{"limit exceeds size", 0, 100, text, text},
		{"non latin equals", 0, 0, nonLatin, nonLatin},
		{"non latin offset", 1, 0, nonLatin, nonLatin[1:]},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			from, err := os.Create(fileIn)
			if err != nil {
				log.Fatal(err)
			}
			defer os.Remove(fileIn)
			_, err = from.WriteString(test.input)
			if err != nil {
				log.Fatal(err)
			}
			defer from.Close()
			err = Copy(fileIn, fileOut, test.offset, test.limit)
			if err != nil {
				log.Fatal(err)
			}
			defer os.Remove(fileOut)
			result, err := ioutil.ReadFile(fileOut)
			if err != nil {
				log.Fatal(err)
			}
			require.Equal(t, string(result), test.result)
		})
	}
	t.Run("binary", func(t *testing.T) {
		from, err := os.Create(fileIn)
		if err != nil {
			log.Fatal(err)
		}
		defer os.Remove(fileIn)
		_, err = from.Write(binaryData)
		if err != nil {
			log.Fatal(err)
		}
		defer from.Close()
		err = Copy(fileIn, fileOut, 0, 0)
		if err != nil {
			log.Fatal(err)
		}
		defer os.Remove(fileOut)
		result, err := ioutil.ReadFile(fileOut)
		if err != nil {
			log.Fatal(err)
		}
		require.Equal(t, result, binaryData)
	})
}

package main

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestReadTestDir(t *testing.T) {
	testdata := make(Environment)
	testdata["BAR"] = EnvValue{"bar", false}
	testdata["EMPTY"] = EnvValue{"", false}
	testdata["FOO"] = EnvValue{"   foo\nwith new line", false}
	testdata["HELLO"] = EnvValue{"\"hello\"", false}
	testdata["UNSET"] = EnvValue{"", true}
	result, err := ReadDir("testdata/env")
	require.Nil(t, err)
	require.Equal(t, result, testdata)
}

func TestEmptyDir(t *testing.T) {
	tmp, err := ioutil.TempDir("", "temp")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := os.Remove(tmp)
		if err != nil {
			log.Fatal(err)
		}
	}()
	testdata := make(Environment)
	result, err := ReadDir(tmp)
	require.Nil(t, err)
	require.Equal(t, result, testdata)
}

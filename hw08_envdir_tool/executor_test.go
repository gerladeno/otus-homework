package main

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestRunCmd(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		code  int
	}{
		{"empty call", []string{""}, 0},
		{"empty echo", []string{"echo"}, 0},
		{"command not found", []string{"ecco"}, 127},
		{"command error", []string{"cp", "?"}, 1},
		{"sh", []string{"./testdata/echo.sh"}, 0},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			code := RunCmd(test.input, nil)
			require.Equal(t, code, test.code)
		})
	}
}

func TestRunCmdErr2(t *testing.T) {
	err := ioutil.WriteFile("tmp.sh", []byte("#!/bin/bash\nexit 2\n"), 0777)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := os.Remove("tmp.sh")
		if err != nil {
			log.Fatal(err)
		}
	}()

	code := RunCmd([]string{"./tmp.sh"}, nil)
	require.Equal(t, code, 2)
}

func TestRunMvCmd(t *testing.T) {
	tmp, err := ioutil.TempDir("", "temp")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := os.RemoveAll(tmp)
		if err != nil {
			log.Fatal(err)
		}
	}()

	nameBefore := filepath.Join(tmp, "tmp")
	nameAfter := filepath.Join(tmp, "tmp_renamed")
	err = ioutil.WriteFile(nameBefore, []byte("just some file"), 0777)
	if err != nil {
		log.Fatal(err)
	}

	code := RunCmd([]string{"mv", nameBefore, nameAfter}, nil)

	contents, err := ioutil.ReadDir(tmp)
	if err != nil {
		log.Fatal(err)
	}

	require.Equal(t, code, 0)
	require.Len(t, contents, 1)
	require.Equal(t, contents[0].Name(), "tmp_renamed")
}


func TestWithEnvRunCmd(t *testing.T) {
	err := ioutil.WriteFile("tmp.sh", []byte("#!/bin/bash\necho $FOO > tmp.out"), 0777)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := os.Remove("tmp.sh")
		if err != nil {
			log.Fatal(err)
		}
	}()

	env := make(Environment)
	env["FOO"] = EnvValue{"text", false}
	code := RunCmd([]string{"./tmp.sh"}, env)
	text, err := ioutil.ReadFile("tmp.out")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := os.Remove("tmp.out")
		if err != nil {
			log.Fatal(err)
		}
	}()

	require.Equal(t, string(text), "text\n")
	require.Equal(t, code, 0)
}
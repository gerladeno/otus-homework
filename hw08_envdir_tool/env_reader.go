package main

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	contents, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	envs := make(Environment)
	for _, file := range contents {
		if !file.Mode().IsRegular() {
			continue
		}
		name := file.Name()
		if strings.Contains(name, "=") {
			log.Printf("filename \"%s\" contains \"=\", skipped...", name)
			continue
		}
		value, err := readEnv(filepath.Join(dir, name))
		if err != nil {
			log.Printf("can't read file \"%s\", %s, skipped...", name, err.Error())
			continue
		}
		envs[name] = *value
	}
	return envs, nil
}

func readEnv(file string) (*EnvValue, error) {
	valueBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	if len(valueBytes) == 0 {
		return &EnvValue{"", true}, nil
	}

	str := string(valueBytes)
	if strings.Contains(str, "\n") {
		str = strings.Split(str, "\n")[0]
	}
	str = strings.ReplaceAll(str, "\x00", "\n")
	value := strings.TrimRight(strings.TrimRight(str, " "), "\t")
	env := EnvValue{
		Value:      value,
		NeedRemove: false,
	}
	return &env, nil
}

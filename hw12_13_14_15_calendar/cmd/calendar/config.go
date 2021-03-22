package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
)

type Config struct {
	Logger  LoggerConf
	Storage StorageConf
	Http    HttpConf
}

type LoggerConf struct {
	Level string `json:"level"`
	Path  string `json:"path"`
}

type StorageConf struct {
	Remote   bool   `json:"remote"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	Ssl      string `json:"ssl"`
}

type HttpConf struct {
	Port int `json:"port"`
}

func NewConfig(path string) Config {
	if path == "" {
		path = filepath.Join("configs", "config.json")
	}
	configJSON, err := ioutil.ReadFile(path)
	if err != nil {
		return defaultConfig()
	}
	config := Config{}
	err = json.Unmarshal(configJSON, &config)
	if err != nil {
		return defaultConfig()
	}
	return config
}

func defaultConfig() Config {
	log.Print("failed to config properly, using default settings...")
	return Config{
		Logger:  LoggerConf{"Debug", "stdout"},
		Storage: StorageConf{Remote: false},
		Http: HttpConf{Port: 3000},
	}
}

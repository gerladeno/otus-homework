package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/cmd"
)

type Config struct {
	Logger cmd.LoggerConf
	Rabbit cmd.RabbitConf
	Sender SenderConfig
}

type SenderConfig struct {
	SenderParam1 string `json:"sender_param_1"`
	SenderParam2 string `json:"sender_param_2"`
	SenderParam3 string `json:"sender_param_3"`
	SenderParam4 string `json:"sender_param_4"`
	SenderParam5 string `json:"sender_param_5"`
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
	var dsn string
	if dsn = os.Getenv("RABBIT_DSN"); dsn == "" {
		dsn = "amqp://guest:guest@localhost:5672/"
	}
	return Config{
		Logger: cmd.LoggerConf{Level: "Debug", Path: "stdout"},
		Rabbit: cmd.RabbitConf{Dsn: dsn},
	}
}

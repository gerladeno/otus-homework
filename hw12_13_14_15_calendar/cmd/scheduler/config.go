package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/cmd"
)

type Config struct {
	Scheduler SchedulerConf
	Storage   cmd.StorageConf
	Logger    cmd.LoggerConf
	Rabbit    cmd.RabbitConf
}

type SchedulerConf struct {
	Period Duration `json:"period"`
}

type Duration time.Duration

func (d *Duration) UnmarshalJSON(b []byte) error {
	durStr := strings.Trim(string(b), "\"")
	duration, err := time.ParseDuration(durStr)
	if err != nil {
		return fmt.Errorf("error parsing duration %s %w", durStr, err)
	}
	*d = Duration(duration)
	return nil
}

func (d *Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d)
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
		Scheduler: SchedulerConf{Period: Duration(time.Minute)},
		Storage:   cmd.StorageConf{Remote: false},
		Logger:    cmd.LoggerConf{Level: "Debug", Path: "stdout"},
		Rabbit:    cmd.RabbitConf{Dsn: dsn},
	}
}

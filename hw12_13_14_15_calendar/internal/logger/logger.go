package logger

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

func New(levelStr, path string) *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	switch strings.ToLower(path) {
	case "", "stderr":
		log.Out = os.Stderr
	case "stdout":
		log.Out = os.Stdout
	default:
		if _, err := os.Stat(path); os.IsNotExist(err) {
			if err = os.MkdirAll(path, 0644); err != nil {
				log.Warn("failed to find or create log directory, using stderr")
				log.Out = os.Stderr
			}
		}
	}

	level, err := logrus.ParseLevel(levelStr)
	if err != nil {
		log.Warn("failed to parse loglevel, using warn")
	}
	log.SetLevel(level)
	return log
}

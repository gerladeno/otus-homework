package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/common"
	"github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/logger"
	"github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/rmq"
	"github.com/sirupsen/logrus"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/sender_config.json", "Path to configuration file")
}

func main() {
	flag.Parse()
	config := NewConfig(configFile)
	log := logger.New(config.Logger.Level, config.Logger.Path)

	if rabbitDsn := os.Getenv("RABBIT_DSN"); rabbitDsn != "" {
		config.Rabbit.Dsn = rabbitDsn
	}
	rabbit, err := rmq.GetRMQConnectionAndDeclare(log, config.Rabbit.Dsn, config.Rabbit.TTL)
	if err != nil {
		log.Fatalf("failed to connect to rmq and declare topic: %s", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	sigCh := make(chan os.Signal, 1)
	go func() {
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGHUP)
		select {
		case <-ctx.Done():
			return
		case <-sigCh:
		}
		log.Info("terminated by syscall...")
		signal.Stop(sigCh)
		cancel()
	}()

	if err = rabbit.ConsumeAndSend(PrepareSender(log, config.Sender)); err != nil {
		log.Fatal("failed to init consumer: ", err)
	}
	if err = rabbit.Close(); err != nil {
		log.Warn("failed to disconnect from rabbit properly: ", err)
	}
}

func PrepareSender(log *logrus.Logger, conf SenderConfig) func([]byte) {
	return func(body []byte) {
		n := common.Notification{}
		if err := json.Unmarshal(body, &n); err != nil {
			log.Warnf("failed to decode a message: %s", string(body))
		}
		if conf.SenderParam1 == "TEST" {
			log.Info("NOTIFICATION: ", n.String())
			host := os.Getenv("CALENDAR_HOST")
			if host == "" {
				host = "http://172.17.0.1:3002"
			} else {
				host = "http://" + host + ":3002"
			}
			c := &http.Client{Transport: http.DefaultTransport}
			c.Transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
			_, err := c.Post(host+"/notify", "application/json", bytes.NewReader(body))
			if err != nil {
				log.Warnf("err notifying: %s. err: %s", n.String(), err)
				return
			}
		} else if conf.SenderParam1 == "INFO" {
			log.Info("NOTIFICATION: ", n.String())
		} else {
			log.Debug("NOTIFICATION: ", n.String())
		}
	}
}

package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/cmd"
	"github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/logger"
	"github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/rmq"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/scheduler_config.json", "Path to configuration file")
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
	scheduler := time.NewTicker(time.Duration(config.Scheduler.Period))
	defer scheduler.Stop()
	ctx, cancel := context.WithCancel(context.Background())

	storage, err := cmd.GetStorage(ctx, log, config.Storage)
	if err != nil {
		log.Fatalf("failed to connect to database: %s", err)
	}

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
		scheduler.Stop()
		if err = rabbit.Close(); err != nil {
			log.Warn("failed to disconnect from rabbit properly: ", err)
		}
	}()

	func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-scheduler.C:
				events, err := storage.ListEventsToNotify(ctx)
				if err != nil {
					log.Warn("failed to retrieve events for notification: ", err)
				}
				rabbit.Notify(events)
			}
		}
	}()
}

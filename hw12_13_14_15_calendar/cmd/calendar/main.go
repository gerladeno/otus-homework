package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/app"
	"github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/server/http"
	"github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/storage/common"
	memorystorage "github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/storage/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := NewConfig(configFile)
	log := logger.New(config.Logger.Level, config.Logger.Path)

	var (
		storage common.Storage
		err     error
	)
	if config.Storage.Remote {
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			config.Storage.Host,
			config.Storage.Port,
			"calendar",
			"calendar",
			config.Storage.Database,
			config.Storage.Ssl)
		storage, err = sqlstorage.New(log, dsn)
		if err != nil {
			log.Fatalf("failed to connect to database: %s", err)
		}
	} else {
		storage = memorystorage.New(log)
	}

	calendar := app.New(log, storage)

	server := internalhttp.NewServer(calendar, storage, log, version, config.HTTP.Port)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGHUP)

		select {
		case <-ctx.Done():
			return
		case <-signals:
		}

		signal.Stop(signals)
		cancel()

		if err := server.Stop(); err != nil {
			log.Error("failed to stop http server: " + err.Error())
		}
	}()

	log.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		log.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}

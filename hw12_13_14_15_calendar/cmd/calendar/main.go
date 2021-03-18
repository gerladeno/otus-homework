package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/storage/common"
	memorystorage "github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/storage/sql"
	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/app"
	"github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/server/http"
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

	_, _ = storage.AddEvent(context.Background(), common.Event{
		ID:         0,
		Title:      "jopa",
		StartTime:  time.Now(),
		Duration:   0,
		InviteList: "",
		Comment:    "gavno",
	})
	_ = storage.EditEvent(context.Background(), 0, common.Event{
		ID:         0,
		Title:      "triJopi",
		StartTime:  time.Now(),
		Duration:   10,
		InviteList: "",
		Comment:    "PoloeGavno",
	})
	_ = storage.RemoveEvent(context.Background(), 0)

	calendar := app.New(log, storage)

	server := internalhttp.NewServer(calendar)
	ctx, cancel := context.WithCancel(context.Background())
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

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
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

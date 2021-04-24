package main

import (
	"context"
	"flag"
	"github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/cmd"
	"github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/app"
	"github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/server/http"
	_ "github.com/jackc/pgx/v4/stdlib"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/calendar_config.json", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := NewConfig(configFile)
	log := logger.New(config.Logger.Level, config.Logger.Path)

	ctx, cancel := context.WithCancel(context.Background())

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
	}()

	storage, err := cmd.GetStorage(ctx, log, config.Storage)
	if err != nil {
		log.Fatalf("failed to connect to database: %s", err)
	}

	calendar := app.New(log, storage)
	handler := internalhttp.NewEventHandler(calendar, log)
	router := internalhttp.NewRouter(handler, log, version)
	httpServer := internalhttp.NewServer(log, router, config.HTTP.Port)
	grpcServer := internalgrpc.NewRPCServer(calendar, log, config.GRPC.Network, config.GRPC.Port)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Infof("starting http server on %d", config.HTTP.Port)
		if err := httpServer.Start(ctx); err != nil {
			log.Error("failed to start http server: " + err.Error())
			cancel()
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Infof("starting grpc server on %d", config.GRPC.Port)
		if err := grpcServer.Start(ctx); err != nil {
			log.Error("failed to start grpc server: " + err.Error())
			cancel()
		}
	}()
	wg.Wait()
}

package main

import (
	"context"
	"github.com/spf13/pflag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	timeoutFlag := pflag.String("timeout", "30s", "specify duration, default - 30s")
	pflag.Parse()
	timeout, err := time.ParseDuration(*timeoutFlag)
	if err != nil {
		log.Fatal("invalid timeout")
	}
	args := pflag.CommandLine.Args()
	if len(args) < 2 {
		log.Fatal("too few arguments")
	}

	telnet := NewTelnetClient(net.JoinHostPort(args[0], args[1]), timeout, os.Stdin, os.Stdout)

	if err := startTelnet(telnet); err != nil {
		log.Fatal(err)
	}
}

func startTelnet(telnet TelnetClient) error {
	err := telnet.Connect()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())

	go worker(telnet.Receive, cancel)
	go worker(telnet.Send, cancel)

	signalChan := make(chan os.Signal)

	go func() {
		signal.Notify(signalChan, syscall.SIGINT)
	}()
	select {
	case <-signalChan:
		cancel()
		signal.Stop(signalChan)
	case <-ctx.Done():
		close(signalChan)
	}
	return nil
}

func worker(handler func() error, cancel context.CancelFunc) {
	if err := handler(); err != nil {
		cancel()
	}
}

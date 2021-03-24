package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/spf13/pflag"
)

func main() {
	timeout := pflag.Duration("timeout", 30*time.Second, "specify duration, default - 30s")
	pflag.Parse()
	args := pflag.CommandLine.Args()
	if len(args) < 2 {
		log.Fatal("too few arguments")
	}
	host := args[0]
	port := args[1]
	if i, err := strconv.Atoi(port); err != nil || i < 1 || i > 65535 {
		log.Fatal("invalid port")
	}

	telnet := NewTelnetClient(net.JoinHostPort(host, port), *timeout, os.Stdin, os.Stdout)

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

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	select {
	case <-signalChan:
		cancel()
		log.Print("terminated by SIGINT...")
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
	cancel()
}

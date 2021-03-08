package main

import (
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

	telnet := NewTelnetClient(net.JoinHostPort(args[0], args[1]), timeout, os.Stdout, os.Stdin)

	if err := startTelnet(telnet); err != nil {
		log.Fatal(err)
	}
}

func startTelnet(telnet TelnetClient) error {
	err := telnet.Connect()
	if err != nil {
		return err
	}

	errCh := make(chan error)
	go func() {
		signalChan := make(chan os.Signal)
		signal.Notify(signalChan, syscall.SIGINT)
		<-signalChan
		err := telnet.Close()
		errCh <- err
	}()
	go func() {
		err := telnet.Send()
		errCh <- err
	}()
	go func() {
		err := telnet.Receive()
		errCh <- err
	}()
	return <-errCh
}

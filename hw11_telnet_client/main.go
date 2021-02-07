package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	timeout time.Duration
	host    string
	port    string
)

func init() {
	flag.DurationVar(&timeout, "timeout", 10, "specify connection timeout")
	flag.StringVar(&host, "host", "localhost", "specify the hosting server")
	flag.StringVar(&port, "port", "3302", "specify server tcp port")
}

func main() {
	flag.Parse()
	if flag.NArg() == 2 {
		host = os.Args[2]
		port = os.Args[3]
	} else {
		flag.Usage()
	}
	ctx, cancel := context.WithCancel(context.Background())

	address := net.JoinHostPort(host, port)
	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		log.Fatalf("Cannot accept: %v", err)
	}
	defer client.Close()
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := client.Receive()
		if err != nil {
			log.Printf("Cannot start receiving goroutine: %v", err.Error())
		}
		cancel()
	}()

	go func() {
		err := client.Send()
		if err != nil {
			log.Printf("Cannot start sending goroutine: %v", err.Error())
		}
		cancel()
	}()

	select {
	case <-signals:
		signal.Stop(signals)
	case <-ctx.Done():
	}
}

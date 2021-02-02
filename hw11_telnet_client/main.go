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
	flag.Duration("timeout", 10, "connection timeout")
}

func main() {
	flag.Parse()
	if len(os.Args) == 3 {
		host = os.Args[1]
		port = os.Args[2]
	} else {
		host = os.Args[2]
		port = os.Args[3]
	}
	address := net.JoinHostPort(host, port)
	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	ctx, cancel := context.WithCancel(context.Background())
	err := client.Connect()
	if err != nil {
		log.Fatalf("Cannot accept: %v", err)
	}
	signals := make(chan os.Signal, 1)

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

	signal.Notify(signals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-signals:
		signal.Stop(signals)
		break
	case <-ctx.Done():
		break
	}
	err = client.Close()
	if err != nil {
		log.Printf("Close client error: %v", err.Error())
	}
	os.Exit(0)
}

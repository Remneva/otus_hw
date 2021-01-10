package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
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

	err := client.Connect()
	if err != nil {
		log.Fatalf("Cannot accept: %v", err)
	}
	wg := sync.WaitGroup{}
	wg.Add(2)
	signals := make(chan os.Signal, 1)

	go func() {
		defer client.Close()
		err := client.Receive()
		if err != nil {
			log.Fatalf("Cannot start receiving goroutine: %v", err.Error())
		}
	}()

	go func() {
		defer client.Close()
		err := client.Send()
		if err != nil {
			log.Fatalf("Cannot start sending goroutine: %v", err.Error())
		}
	}()

	go func() {
		signal.Notify(signals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		<-signals
		signal.Stop(signals)

		err := client.Close()
		if err != nil {
			log.Fatalf("Close client error: %v", err.Error())
		}
		os.Exit(0)
	}()
	wg.Wait()
}

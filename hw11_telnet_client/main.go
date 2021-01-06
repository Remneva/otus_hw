package main

import (
	"flag"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

var (
	host, port      string
	timeout         time.Duration
	timeoutDuration time.Duration
)

func init() {
	flag.Duration("timeout", 10, "connection timeout")
	flag.StringVar(&host, "host", "localhost", "host")
	flag.StringVar(&port, "port", "4242", "port")
}

func main() {
	flag.Parse()
	host := os.Args[2]
	port := os.Args[3]
	timeoutDuration = timeout * time.Second
	address := net.JoinHostPort(host, port)
	out := os.Stdout
	in := os.Stdin
	client := NewTelnetClient(address, timeoutDuration, in, out)

	err := client.Connect()
	if err != nil {
		log.Fatalf("Cannot accept: %v", err)
	}
	defer client.Close()
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		err := client.Receive()
		if err != nil {
			log.Fatalf("Cannot start receiving goroutine: %v", err.Error())
		}
	}()

	go func() {
		defer wg.Done()
		err := client.Send()
		if err != nil {
			log.Fatalf("Cannot start sending goroutine: %v", err.Error())
		}
	}()

	wg.Wait()
}

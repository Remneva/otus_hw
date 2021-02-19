package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/configs"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/pkg/logger"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/pkg/rabbit"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/pkg/storage/sql"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

var config string

func init() {
	flag.StringVar(&config, "config", "./calendar_sender/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config, err := configs.Read(config)
	if err != nil {
		log.Println("failed to read config")
	}
	logg, err := logger.NewLogger(config.Logger.Level, "prod", config.Logger.Path)
	if err != nil {
		log.Println("failed to create logger")
	}

	connection, err := amqp.Dial(config.AMQP.URI)
	if err != nil {
		logg.Fatal("dial: ", zap.Error(err))
	}
	defer connection.Close()

	logg.Info("AMQP", zap.String("got Connection, getting Channel", config.AMQP.URI))
	channel, err := connection.Channel()
	if err != nil {
		logg.Error("channel: ", zap.Error(err))
	}

	storage := sql.NewStorage(logg)
	if err := storage.Connect(ctx, config.PSQL.DSN); err != nil {
		logg.Fatal("fail connection")
	}
	defer storage.Close()

	c := rabbit.NewRabbit(ctx, channel, connection, logg, config, storage)

	err = c.Declare()
	if err != nil {
		logg.Error("declare:", zap.Error(err))
	}
	msgs, err := c.Consume()
	if err != nil {
		logg.Fatal("consume error", zap.Error(err))
	}
	go signalChan(*c)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		c.Handle(msgs, c.Done)
	}()
	wg.Wait()
}

func signalChan(q rabbit.Rabbit) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	fmt.Printf("Got %v...\n", <-signals)

	if err := q.Shutdown(); err != nil {
		q.Log.Error("consumer cancel failed", zap.Error(err))
	}
}

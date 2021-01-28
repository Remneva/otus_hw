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
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/rabbit"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

var config string

func init() {
	flag.StringVar(&config, "config", "./calendar_scheduler/config.toml", "Path to configuration file")
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

	defer cancel()
	storage := sql.NewStorage(logg)
	if err := storage.Connect(ctx, config.PSQL.DSN); err != nil {
		logg.Fatal("fail connection")
	}
	defer storage.Close()
	connection, err := amqp.Dial(config.AMQP.URI)
	if err != nil {
		logg.Fatal("dial: ", zap.Error(err))
	}
	defer connection.Close()

	logg.Info("AMQP", zap.String("got Connection, getting Channel", config.AMQP.URI))
	channel, err := connection.Channel()
	if err != nil {
		logg.Fatal("channel", zap.Error(err))
	}
	q := rabbit.NewRabbit(ctx, channel, connection, logg, config, storage)

	err = q.Declare()
	if err != nil {
		logg.Fatal("exchange Declare", zap.Error(err))
	}
	go signalChan(*q)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		cronStart(storage, *q)
	}()
	wg.Wait()
}

func signalChan(q rabbit.Rabbit) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	fmt.Printf("Got %v...\n", <-signals)

	close(q.Done)
	err := q.Shutdown()
	if err != nil {
		q.Log.Error("consumer cancel failed", zap.Error(err))
	}
}

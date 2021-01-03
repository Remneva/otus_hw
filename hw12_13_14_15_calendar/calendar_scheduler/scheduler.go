package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/configs"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/logger"
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
	logg, err := logger.NewLogger(config.Logger.Level, config.Logger.Path)
	if err != nil {
		log.Println("failed to create logger")
	}

	defer cancel()
	storage := new(sql.Storage)
	if err := storage.Connect(ctx, config.PSQL.DSN, logg); err != nil {
		logg.Fatal("fail connection")
	}
	connection, err := amqp.Dial(config.AMQP.URI)
	if err != nil {
		logg.Fatal("Dial: ", zap.Error(err))
	}
	defer connection.Close()

	logg.Info("AMQP", zap.String("got Connection, getting Channel", config.AMQP.URI))
	channel, err := connection.Channel()
	if err != nil {
		log.Fatalf("Channel: %s", err)
	}

	q := &Producer{
		Channel: channel,
		Conn:    connection,
		Log:     logg,
		C:       config,
		Ctx:     ctx,
		Done:    make(chan error),
	}

	err = q.Declare()
	if err != nil {
		log.Fatalf("Exchange Declare: %s", err)
	}

	// Reliable publisher confirms require confirm.select support from the
	// connection.
	if q.C.AMQP.Reliable {
		q.Log.Info("enabling publishing confirms.")
		if err := q.Channel.Confirm(false); err != nil {
			log.Fatalf("Channel could not be put into confirm mode: %s", err)
		}

		confirms := q.Channel.NotifyPublish(make(chan amqp.Confirmation, 1))

		defer ConfirmOne(confirms)
	}

	for true { //nolint
		select {
		case <-q.Ctx.Done():
			return
		default:
			q.CheckDB(storage)
			q.Log.Info("Waiting for the next checkup...")
			fmt.Println("TIMEOUT", q.C.AMQP.Timeout)
			time.Sleep(q.C.AMQP.Timeout * time.Second)
		}
	}

	err = q.Shutdown()
	if err != nil {
		q.Log.Error("Consumer cancel failed: %s", zap.Error(err))
	}
}

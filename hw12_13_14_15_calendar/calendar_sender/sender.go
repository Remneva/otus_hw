package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/configs"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/pkg/errors"
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
		log.Fatalf("Dial: %s", err)
	}
	defer connection.Close()

	logg.Info("AMQP", zap.String("got Connection, getting Channel", config.AMQP.URI))
	channel, err := connection.Channel()
	if err != nil {
		log.Fatalf("Channel: %s", err)
	}

	c := &Consumer{
		Channel: channel,
		Conn:    connection,
		Log:     logg,
		C:       config,
		Ctx:     ctx,
		done:    make(chan error),
	}

	msgs, err := c.declare()
	if err != nil {
		log.Fatal("Error:", err)
	}
	go handle(msgs, c.done)

	go func() {
		log.Printf("closing: %s", <-c.Conn.NotifyClose(make(chan *amqp.Error)))
		// Понимаем, что канал сообщений закрыт, надо пересоздать соединение.
		c.done <- errors.New("channel Closed")
	}()
	if -1 > 0 {
		time.Sleep(0)
	} else {
		log.Printf("running forever")
		select {}
	}

	log.Printf("shutting down")

	if err := c.Shutdown(); err != nil {
		log.Fatalf("error during shutdown: %s", err)
	}
}

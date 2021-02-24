package main

import (
	"context"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/configs"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/pkg/logger"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/pkg/rabbit"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/pkg/storage/sql"
	"log"
	"testing"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var dsn = "host=postgres port=5432 user=test password=test dbname=exampledb sslmode=disable"
var uri = "amqp://guest:guest@rabbit:5672/"

func TestStoreSuite(t *testing.T) {
	suite.Run(t, new(suiteTestIntegration))
}

type suiteTestIntegration struct {
	suite.Suite
	ctx context.Context
	r   *rabbit.Rabbit
	s   *sql.Storage
}

func (s *suiteTestIntegration) SetupTest() {
	ctx := context.Background()
	var z zapcore.Level

	logg, err := logger.NewLogger(z, "dev", "/dev/null")
	if err != nil {
		log.Fatal("failed to create logger")
	}

	storage := sql.NewStorage(logg)
	err = storage.Connect(ctx, dsn)
	if err != nil {
		log.Fatal("failed to connect db")
	}
	connection, err := amqp.Dial(uri)
	if err != nil {
		logg.Fatal("dial: ", zap.Error(err))
	}
	defer connection.Close()

	logg.Info("AMQP", zap.String("got Connection, getting Channel", uri))
	channel, err := connection.Channel()
	if err != nil {
		logg.Fatal("channel", zap.Error(err))
	}

	config, err := configs.Read("config.toml")
	if err != nil {
		log.Fatal("failed to read config")
	}

	r := rabbit.NewRabbit(ctx, channel, connection, logg, config, storage)
	err = r.Declare()
	r.Conn, _ = amqp.Dial(uri)
	r.Channel, _ = r.Conn.Channel()

	s.ctx = ctx
	s.s = storage
	s.r = r
}

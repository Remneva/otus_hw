package rabbit

import (
	"context"

	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/configs"
	store "github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

var _ rabbitmq = (*Rabbit)(nil)

type Rabbit struct {
	Channel *amqp.Channel
	Conn    *amqp.Connection
	Log     *zap.Logger
	C       configs.Config
	Done    chan error
	Ctx     context.Context
	Repo    store.EventsStorage
}

type rabbitmq interface {
	Declare() error
	Publish(ev store.Event) error
	Shutdown() error
	Consume() (<-chan amqp.Delivery, error)
	Handle(deliveries <-chan amqp.Delivery, done chan error)
}

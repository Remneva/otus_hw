package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/configs"
	store "github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

var _ EventQueue = (*Producer)(nil)

type (
	EvChan = chan store.Event
)

type Producer struct {
	Channel *amqp.Channel
	Conn    *amqp.Connection
	Log     *zap.Logger
	C       configs.Config
	Ctx     context.Context
	Done    chan error
}

type EventQueue interface {
	checkingDB(storage *sql.Storage)
	publish(ev store.Event) error
}

func (q *Producer) Shutdown() error {
	// will close() the deliveries channel
	if err := q.Channel.Cancel(q.C.AMQP.ConsumerTag, true); err != nil {
		q.Log.Error("Consumer cancel failed: %s", zap.Error(err))
		return errors.Wrap(err, "Consumer cancel failed")
	}

	if err := q.Conn.Close(); err != nil {
		q.Log.Error("AMQP connection close error:", zap.Error(err))
		return errors.Wrap(err, "AMQP connection close error")
	}

	defer log.Printf("AMQP shutdown OK")

	// wait for handle() to exit
	return <-q.Done
}

func (q *Producer) checkingDB(storage *sql.Storage) {
	if err := storage.Connect(q.Ctx, q.C.PSQL.DSN, q.Log); err != nil {
		q.Log.Fatal("fail connection")
	}
	eventChan := make(EvChan)

	go func() {
		for {
			select {
			case <-q.Ctx.Done():
				return
			default:
				events, err := storage.GetEvents(q.Ctx)
				if err != nil {
					q.Log.Fatal("fail request")
				}
				for _, ev := range events {
					q.Log.Info("ev: ", zap.String("Nahuatl name", ev.Title))

					if ev.StartTime.After(time.Now().Truncate(time.Duration(-60) * time.Minute)) {
						eventChan <- ev
					}
				}
				time.Sleep(5 * time.Second)
			}
		}
	}()

	go func() {
		for {
			select {
			case ev := <-eventChan:
				if err := q.publish(ev); err != nil {
					q.Log.Error("Publish message error", zap.Error(err))
				}
				q.Log.Info("Published OK", zap.Int("id", int(ev.ID)))
			case <-q.Ctx.Done():
				return
			default:
				time.Sleep(5 * time.Second)
			}
		}
	}()
}

func (q *Producer) CheckDB(storage *sql.Storage) {
	start := time.Now()

	oneYearLater := start.AddDate(-1, 0, 0)

	q.Conn, _ = amqp.Dial(q.C.AMQP.URI)
	q.Channel, _ = q.Conn.Channel()
	events, err := storage.GetEvents(q.Ctx)
	if err != nil {
		q.Log.Fatal("fail request", zap.Error(err))
	}

	for _, ev := range events {
		if ev.StartTime.After(time.Now().Add((-30) * time.Minute)) {
			go func(ev store.Event) {
				if err := q.publish(ev); err != nil {
					q.Log.Error("Publish message error", zap.Error(err))
				}
			}(ev)
		} else if ev.StartTime.Before(oneYearLater) {
			err = storage.DeleteEvent(q.Ctx, ev.ID)
			if err != nil {
				q.Log.Error("Delete failed", zap.Int("id", int(ev.ID)))
			}
			q.Log.Info("Outdated event deleted", zap.Int("id", int(ev.ID)))
		}
	}
}

func (q *Producer) publish(ev store.Event) error {
	reqBodyBytes := new(bytes.Buffer)
	err := json.NewEncoder(reqBodyBytes).Encode(ev)
	if err != nil {
		return errors.Wrap(err, "Encoding error")
	}
	body := reqBodyBytes.Bytes()

	if err := q.Channel.Publish(
		q.C.AMQP.ExchangeName, // publish to an exchange
		q.C.AMQP.RoutingKey,   // routing to 0 or more queues
		false,                 // mandatory
		false,                 // immediate
		amqp.Publishing{
			Body: body,
		},
	); err != nil {
		q.Log.Info("Error publish")
		return errors.Wrap(err, "Exchange Publish")
	}
	q.Log.Info("Published OK", zap.Int("id", int(ev.ID)))
	return nil
}

func ConfirmOne(confirms <-chan amqp.Confirmation) {
	log.Printf("waiting for confirmation of one publishing")

	if confirmed := <-confirms; confirmed.Ack {
		log.Printf("confirmed delivery with delivery tag: %d", confirmed.DeliveryTag)
	} else {
		log.Printf("failed delivery of delivery tag: %d", confirmed.DeliveryTag)
	}
}

func (q *Producer) Declare() error {
	if err := q.Channel.ExchangeDeclare(
		q.C.AMQP.ExchangeName, // name
		q.C.AMQP.ExchangeType, // type
		true,                  // durable
		false,                 // auto-deleted
		false,                 // internal
		false,                 // noWait
		nil,                   // arguments
	); err != nil {
		return errors.Errorf("Exchange Declare: %s", err)
	}
	q.Log.Info("Exchange declared")
	return nil
}

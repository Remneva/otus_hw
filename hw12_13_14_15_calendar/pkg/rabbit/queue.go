package rabbit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/configs"
	store "github.com/Remneva/otus_hw/hw12_13_14_15_calendar/pkg/storage"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type Message struct {
	Message string
}

func NewRabbit(ctx context.Context, channel *amqp.Channel, connection *amqp.Connection, logg *zap.Logger, config configs.Config, ev store.EventsStorage) *Rabbit {
	r := &Rabbit{
		Channel: channel,
		Conn:    connection,
		Log:     logg,
		C:       config,
		Ctx:     ctx,
		Repo:    ev,
		Done:    make(chan error),
	}
	return r
}

func (r *Rabbit) Declare() error {
	if err := r.Channel.ExchangeDeclare(
		r.C.AMQP.ExchangeName, // name
		r.C.AMQP.ExchangeType, // type
		true,                  // durable
		false,                 // auto-deleted
		false,                 // pkg
		false,                 // noWait
		nil,
	); err != nil {
		return fmt.Errorf("exchange Declare error: %w", err)
	}
	r.Log.Info("exchange declared")
	if _, err := r.Channel.QueueDeclare(
		r.C.AMQP.Queue, // name of the queue
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // noWait
		nil,            // arguments
	); err != nil {
		return fmt.Errorf("queue Declare error: %w", err)
	}
	// Число сообщений, которые можно подтвердить за раз.
	err := r.Channel.Qos(50, 0, false)
	if err != nil {
		return fmt.Errorf("etting qos error: %w", err)
	}

	// Создаём биндинг (правило маршрутизации), если оно ещё не создано
	if err := r.Channel.QueueBind(
		r.C.AMQP.Queue,
		r.C.AMQP.RoutingKey,
		r.C.AMQP.ExchangeName,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("queue Bind error: %w", err)
	}
	return nil
}

func (r *Rabbit) Publish(ev store.Event) error {
	reqBodyBytes := new(bytes.Buffer)
	err := json.NewEncoder(reqBodyBytes).Encode(ev)
	if err != nil {
		return fmt.Errorf("encoding error: %w", err)
	}
	body := reqBodyBytes.Bytes()

	if err := r.Channel.Publish(
		r.C.AMQP.ExchangeName, // publish to an exchange
		r.C.AMQP.RoutingKey,   // routing to 0 or more

		// s
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			Body: body,
		},
	); err != nil {
		r.Log.Info("error publish")
		return fmt.Errorf("exchange Publish error: %w", err)
	}
	r.Log.Info("published OK", zap.Int("id", int(ev.ID)))
	return nil
}

func (r *Rabbit) Shutdown() error {
	// will close() the deliveries channel
	if err := r.Channel.Cancel(r.C.AMQP.ConsumerTag, true); err != nil {
		r.Log.Error("consumer cancel failed: %s", zap.Error(err))
		return fmt.Errorf("consumer cancel failed: %w", err)
	}

	if err := r.Conn.Close(); err != nil {
		r.Log.Error("AMQP connection close error:", zap.Error(err))
		return fmt.Errorf("AMQP connection close error: %w", err)
	}
	r.Log.Info("AMQP connection is closed")
	defer r.Log.Info("AMQP shutdown OK")

	// wait for handle() to exit
	return <-r.Done
}

func (r *Rabbit) Consume() (<-chan amqp.Delivery, error) {
	queue, err := r.Channel.QueueDeclare(
		r.C.AMQP.Queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("queue Declare error: %w", err)
	}
	msgs, err := r.Channel.Consume(
		queue.Name,
		r.C.AMQP.ConsumerTag,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("queue Consume error: %w", err)
	}
	return msgs, nil
}

func (r *Rabbit) Handle(deliveries <-chan amqp.Delivery, done chan error) {
	var m Event

	for d := range deliveries {
		r.Log.Info("got msg:", zap.ByteString("body:", d.Body),
			zap.Int("byte:", len(d.Body)), zap.String("tag:", strconv.FormatUint(d.DeliveryTag, 10)))

		err := json.Unmarshal(d.Body, &m)
		if err != nil {
			r.Log.Error("unmarshal error", zap.Error(err))
		}
		err = r.Repo.ChangeStatusByID(r.Ctx, m.ID)
		if err != nil {
			r.Log.Error("status query error", zap.Error(err))
		}
		err = d.Ack(false)
		if err != nil {
			r.Log.Error("queue declare error", zap.Error(err))
		}
	}
	r.Log.Info("handle: deliveries channel closed")
	done <- nil
}

type Event struct {
	ID          int64  `json:"ID"`
	Owner       int64  `json:"Owner"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
	StartDate   string `json:"StartDate"`
	StartTime   string `json:"StartTime"`
	EndDate     string `json:"EndDate"`
	EndTime     string `json:"EndTime"`
}

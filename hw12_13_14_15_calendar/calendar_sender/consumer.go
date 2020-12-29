package main

import (
	"context"
	"fmt"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/configs"
	"github.com/cenkalti/backoff/v3"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"log"
	"strconv"
	"time"
)

type Message struct {
	Ctx  context.Context
	Data []byte
}
type Consumer struct {
	Channel     *amqp.Channel
	Conn        *amqp.Connection
	Log         *zap.Logger
	C           configs.Config
	done        chan error
	maxInterval time.Duration
	Ctx         context.Context
}

func (c *Consumer) Shutdown() error {
	// will close() the deliveries channel
	if err := c.Channel.Cancel(c.C.AMQP.ConsumerTag, true); err != nil {
		c.Log.Error("Consumer cancel failed: %s", zap.Error(err))
		return err
	}

	if err := c.Conn.Close(); err != nil {
		c.Log.Error("AMQP connection close error:", zap.Error(err))
		return err
	}

	defer log.Printf("AMQP shutdown OK")

	// wait for handle() to exit
	return <-c.done
}

func (c *Consumer) connect() error {
	var err error

	c.Conn, err = amqp.Dial(c.C.AMQP.URI)
	if err != nil {
		return fmt.Errorf("dial: %s", err)
	}

	c.Channel, err = c.Conn.Channel()
	if err != nil {
		return fmt.Errorf("channel: %s", err)
	}

	go func() {
		log.Printf("closing: %s", <-c.Conn.NotifyClose(make(chan *amqp.Error)))
		// Понимаем, что канал сообщений закрыт, надо пересоздать соединение.
		c.done <- errors.New("channel Closed")
	}()

	if err = c.Channel.ExchangeDeclare(
		c.C.AMQP.RoutingKey,
		c.C.AMQP.ExchangeName,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("exchange declare: %s", err)
	}

	return nil
}

func handle(deliveries <-chan amqp.Delivery, done chan error) {
	for d := range deliveries {
		log.Printf(
			"got %dB delivery: [%v] %q",
			len(d.Body),
			d.DeliveryTag,
			d.Body,
		)
		d.Ack(false)
	}
	log.Printf("handle: deliveries channel closed")
	done <- nil
}

func (c *Consumer) declare() (<-chan amqp.Delivery, error) {
	queue, err := c.Channel.QueueDeclare(
		c.C.AMQP.Queue,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("queue Declare: %s", err)
	}

	// Число сообщений, которые можно подтвердить за раз.
	err = c.Channel.Qos(50, 0, false)
	if err != nil {
		return nil, fmt.Errorf("error setting qos: %s", err)
	}

	// Создаём биндинг (правило маршрутизации).
	if err = c.Channel.QueueBind(
		queue.Name,
		c.C.AMQP.RoutingKey,
		c.C.AMQP.ExchangeName,
		false,
		nil,
	); err != nil {
		return nil, fmt.Errorf("queue Bind: %s", err)
	}

	msgs, err := c.Channel.Consume(
		queue.Name,
		c.C.AMQP.ConsumerTag,
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("queue Consume: %s", err)
	}

	err = c.Channel.Publish(
		c.C.AMQP.ExchangeName, // exchange
		c.C.AMQP.RoutingKey,   // routing key
		false,                 // mandatory
		false,                 // immediate
		amqp.Publishing{
			Body: []byte(strconv.Itoa(111222)),
		})
	if err != nil {
		fmt.Println("queue Publish: %s", err.Error())
		return nil, fmt.Errorf("queue Publish: %s", err.Error())
	}
	fmt.Println("declare queue.name: ", queue.Name)

	return msgs, nil
}

func (c *Consumer) reConnect(ctx context.Context) (<-chan amqp.Delivery, error) {
	be := backoff.NewExponentialBackOff()
	be.MaxElapsedTime = time.Minute
	be.InitialInterval = 1 * time.Second
	be.Multiplier = 2
	be.MaxInterval = 15 * time.Second

	b := backoff.WithContext(be, ctx)
	for {
		d := b.NextBackOff()
		if d == backoff.Stop {
			return nil, fmt.Errorf("stop reconnecting")
		}

		select {
		case <-time.After(d):
			if err := c.connect(); err != nil {
				log.Printf("could not connect in reconnect call: %+v", err)
				continue
			}
			msgs, err := c.declare()
			if err != nil {
				fmt.Printf("Couldn't connect: %+v", err)
				continue
			}

			return msgs, nil
		}
	}
}

func (c *Consumer) Handle(ctx context.Context, threads int) error {
	var err error
	if err = c.connect(); err != nil {
		return fmt.Errorf("error: %v", err)
	}
	out := make(chan Message)
	msgs, err := c.declare()
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}

	for {
		for i := 0; i < threads; i++ {
			del := <-msgs
			time.Sleep(3 * time.Second)
			fmt.Println("body", string(del.Body))
			if err := del.Ack(false); err != nil {
				log.Println(err)
			}
			fmt.Println("done")
			msg := Message{
				Ctx:  context.TODO(),
				Data: del.Body,
			}
			out <- msg
		}

		if <-c.done != nil {
			msgs, err = c.reConnect(ctx)
			if err != nil {
				return fmt.Errorf("reconnecting Error: %s", err)
			}
		}
		fmt.Println("Reconnected... possibly")
	}
}

package main

import (
	"context"
	"flag"
	"fmt"
	configs "github.com/Remneva/otus_hw/hw12_13_14_15_calendar/calendar/configs"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/calendar/internal/logger"
	"go.uber.org/zap"
	"log"


	"github.com/streadway/amqp"
)

var config string
var (
//  uri          = flag.String("uri", "amqp://guest:guest@localhost:5672/", "AMQP URI")
//	exchangeName = flag.String("exchange", "test-exchange", "Durable AMQP exchange name")
//	exchangeType = flag.String("exchange-type", "direct", "Exchange type - direct|fanout|topic|x-custom")
//	routingKey   = flag.String("key", "test-key", "AMQP routing key")
//	body         = flag.String("body", "foobar", "Body of message")
//	reliable = flag.Bool("reliable", true, "Wait for the publisher confirmation before exiting")
)

func init() {
	flag.StringVar(&config, "config", "config.toml", "Path to configuration file")
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
	if err := publish(ctx, config, logg); err != nil {
		log.Printf("%s", err)
	}

}

func publish(ctx context.Context, c configs.Config, logger *zap.Logger) error {

	// This function dials, connects, declares, publishes, and tears down,
	// all in one go. In a real service, you probably want to maintain a
	// long-lived connection as state, and publish against that.

	logger.Info("published OK", zap.Int("size", len(c.AMQP.Body)))
	connection, err := amqp.Dial(c.AMQP.Uri)
	if err != nil {
		return fmt.Errorf("Dial: %s", err)
	}
	defer connection.Close()

	logger.Info("AMQP", zap.String("got Connection, getting Channel", c.AMQP.Uri))
	channel, err := connection.Channel()
	if err != nil {
		return fmt.Errorf("Channel: %s", err)
	}

	logger.Info("AMQP", zap.String("got Channel, declaring", c.AMQP.ExchangeType), zap.String("Exchange", c.AMQP.ExchangeName))
	if err := channel.ExchangeDeclare(
		c.AMQP.ExchangeName, // name
		c.AMQP.ExchangeType, // type
		true,                // durable
		false,               // auto-deleted
		false,               // internal
		false,               // noWait
		nil,                 // arguments
	); err != nil {
		return fmt.Errorf("Exchange Declare: %s", err)
	}

	// Reliable publisher confirms require confirm.select support from the
	// connection.
	if c.AMQP.Reliable {
		logger.Info("AMQP", zap.String("enabling publishing confirms", c.AMQP.Uri))
		if err := channel.Confirm(false); err != nil {
			return fmt.Errorf("Channel could not be put into confirm mode: %s", err)
		}

		confirms := channel.NotifyPublish(make(chan amqp.Confirmation, 1))

		defer confirmOne(confirms)
	}

	body := c.AMQP.Body
	logger.Info("AMQP", zap.String("declared Exchange, publishing body", body), zap.Int("size", len(body)))
	if err = channel.Publish(
		c.AMQP.ExchangeName, // publish to an exchange
		c.AMQP.RoutingKey,   // routing to 0 or more queues
		false,               // mandatory
		false,               // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            []byte(body),
			DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,              // 0-9
			// a bunch of application/implementation-specific fields
		},
	); err != nil {
		return fmt.Errorf("Exchange Publish: %s", err)
	}

	return nil
}

// One would typically keep a channel of publishings, a sequence number, and a
// set of unacknowledged sequence numbers and loop until the publishing channel
// is closed.
func confirmOne(confirms <-chan amqp.Confirmation) {
	log.Printf("waiting for confirmation of one publishing")

	if confirmed := <-confirms; confirmed.Ack {
		log.Printf("confirmed delivery with delivery tag: %d", confirmed.DeliveryTag)
	} else {
		log.Printf("failed delivery of delivery tag: %d", confirmed.DeliveryTag)
	}
}

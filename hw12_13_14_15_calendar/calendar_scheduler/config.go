package main

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

type Config struct {
	AMQP AmqpConfig
}

type AmqpConfig struct {
	Uri          string
	ExchangeName string
	ExchangeType string
	RoutingKey   string
	Body         string
	Reliable     bool
}

func Read(path string) (c Config, err error) {
	_, err = toml.DecodeFile(path, &c)
	if err != nil {
		return Config{}, errors.Wrap(err, "DecodeFile failed")
	}
	fmt.Println("Uri: ", c.AMQP.Uri)
	fmt.Println("ExchangeName: ", c.AMQP.ExchangeName)
	fmt.Println("ExchangeType: ", c.AMQP.ExchangeType)
	fmt.Println("RoutingKey: ", c.AMQP.RoutingKey)
	fmt.Println("Body: ", c.AMQP.Body)
	fmt.Println("Reliable: ", c.AMQP.Reliable)
	return
}

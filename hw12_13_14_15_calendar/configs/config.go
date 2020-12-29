package configs

import (
	"fmt"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Logger LoggerConfig
	PSQL   PSQLConfig
	Port   PortConfig
	AMQP   AMQPConfig
	Mode   ModeConfig
}

type ModeConfig struct {
	MemMode bool
}

type LoggerConfig struct {
	Level zapcore.Level
	Path  string
}

type PSQLConfig struct {
	DSN string
}

type PortConfig struct {
	HTTP string
	Grpc string
}

type AMQPConfig struct {
	URI          string
	ExchangeName string
	ExchangeType string
	RoutingKey   string
	Body         string
	Reliable     bool
	Timeout      time.Duration
	ConsumerTag  string
	Queue        string
}

func Read(path string) (c Config, err error) {
	_, err = toml.DecodeFile(path, &c)
	if err != nil {
		return Config{}, errors.Wrap(err, "DecodeFile failed")
	}
	fmt.Println("DSN: ", c.PSQL.DSN)
	fmt.Println("logger: ", c.Logger.Level)
	fmt.Println("path: ", c.Logger.Path)
	fmt.Println("port http: ", c.Port.HTTP)
	fmt.Println("port grpc: ", c.Port.Grpc)
	fmt.Println("URI: ", c.AMQP.URI)
	fmt.Println("ExchangeName: ", c.AMQP.ExchangeName)
	fmt.Println("ExchangeType: ", c.AMQP.ExchangeType)
	fmt.Println("RoutingKey: ", c.AMQP.RoutingKey)
	fmt.Println("Body: ", c.AMQP.Body)
	fmt.Println("Reliable: ", c.AMQP.Reliable)
	fmt.Println("MemMode: ", c.Mode.MemMode)
	fmt.Println("Queue: ", c.AMQP.Queue)
	return
}

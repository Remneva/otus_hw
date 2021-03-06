package configs

import (
	"fmt"

	"github.com/BurntSushi/toml"
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
	Timeout      uint64
	ConsumerTag  string
	Queue        string
}

func Read(path string) (c Config, err error) {
	_, err = toml.DecodeFile(path, &c)
	if err != nil {
		return Config{}, fmt.Errorf("decodeFile failed: %w", err)
	}
	return
}

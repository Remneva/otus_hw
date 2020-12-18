package configs

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Logger LoggerConfig
	PSQL   PSQLConfig
	Port   PortConfig
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
	return
}

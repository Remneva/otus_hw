package configs

import (
	"errors"
	"fmt"

	"github.com/BurntSushi/toml"
	"go.uber.org/zap/zapcore"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.

var ErrConfig = errors.New("cant`t read conf.toml file")

type Config struct {
	Logger LoggerConf
	PSQL   PSQLConfig
}

type LoggerConf struct {
	Level zapcore.Level
	Path  string
}

type PSQLConfig struct {
	DSN string
}

func NewConfig(fpath string) (Config, error) {
	_, err := Read(fpath)
	if err != nil {
		return Config{}, err
	}
	return Config{}, nil
}

func Read(fpath string) (c Config, err error) {
	_, err = toml.DecodeFile(fpath, &c)
	if err != nil {
		return Config{}, err
	}
	fmt.Println("logger: ", &c.Logger.Level)
	fmt.Println("path: ", &c.Logger.Path)

	return
}

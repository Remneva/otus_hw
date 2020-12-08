package configs

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/apex/log"
	"go.uber.org/zap/zapcore"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
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

func NewConfig(fpath string) Config {
	_, err := Read(fpath)
	if err != nil {
		log.Fatal("cant`t read conf.toml file")
	}
	return Config{}
}

func Read(fpath string) (c Config, err error) {
	_, err = toml.DecodeFile(fpath, &c)
	if err != nil {
		fmt.Println("error: ", err)
	}
	fmt.Println("logger: ", &c.Logger.Level)
	fmt.Println("path: ", &c.Logger.Path)

	return
}

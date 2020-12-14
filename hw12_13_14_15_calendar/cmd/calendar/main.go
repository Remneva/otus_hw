package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/configs"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/apex/log"
	"os"
	"os/signal"
	"time"
)

var config string

func init() {
	flag.StringVar(&config, "config", "config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config, _ := configs.Read(config)
	logg, err := logger.NewLogger(config.Logger.Level, config.Logger.Path)
	if err != nil {
		log.Fatal("failed to create logger")
	}
	storage := new(memorystorage.Storage)

	if err := storage.Connect(ctx, config.PSQL.DSN); err != nil {
		logg.Fatal("fail connection")
	}
	application, err := app.New(logg, storage)
	if err != nil {
		fmt.Println(err.Error())
		logg.Fatal("failed to start application")
	}

	logg.Info("calendar is running...")
	err = application.Run(ctx)
	if err != nil {
		fmt.Println(err)
		logg.Fatal("failed to start application")
	}

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals)

		<-signals
		signal.Stop(signals)
		cancel()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		err := application.Stop(ctx)
		if err != nil {
			logg.Error(err.Error())
		}
	}()
}

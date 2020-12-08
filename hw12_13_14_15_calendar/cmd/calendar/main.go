package main //nolint:golint,stylecheck

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/configs"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/server/http"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/storage/memory"
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
	logg := logger.NewLogger(config.Logger.Level, config.Logger.Path)
	storage := new(memorystorage.Storage)

	if err := storage.Connect(ctx, config.PSQL.DSN); err != nil {
		fmt.Println(err)
	}
	app, err := app.New(logg, storage)
	if err != nil {
		fmt.Println(err.Error())
		logg.Fatal("failed to start app")
	}
	logg.Info("calendar is running...")
	err = app.Run(ctx)
	if err != nil {
		fmt.Println("run error: ", err)
	}
	server := internalhttp.NewServer(logg, app)
	defer cancel()

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals)

		<-signals
		signal.Stop(signals)
		cancel()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
	}
}

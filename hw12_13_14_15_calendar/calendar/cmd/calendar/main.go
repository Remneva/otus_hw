package main

import (
	"context"
	"flag"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/storage/sql"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/configs"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/server/http"
	"github.com/apex/log"
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

	config, err := configs.Read(config)
	if err != nil {
		log.Fatal("failed to read config")
	}
	logg, err := logger.NewLogger(config.Logger.Level, config.Logger.Path)
	if err != nil {
		log.Fatal("failed to create logger")
	}

	storage := new(sql.Storage)
	if err := storage.Connect(ctx, config.PSQL.DSN, logg); err != nil {
		log.Fatal("fail connection")
	}
	application, err := app.New(logg, storage)
	if err != nil {
		logg.Fatal("failed to start application")
	}

	logg.Info("calendar is running...")
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		application.Log.Info("http is running...")
		_, mux := internalhttp.NewHandler(ctx, application)
		server := internalhttp.New()
		http, err := server.NewServer(mux, config.Port.HTTP)
		if err != nil {
			application.Log.Fatal("failed to start http server: " + err.Error())
		}
		if err := http.Start(ctx); err != nil {
			application.Log.Error("failed to start http server: " + err.Error())
			os.Exit(1)
		}
		wg.Done()
	}()

	go func() {
		application.Log.Info("grpc is running...")
		service := internalgrpc.New(ctx, application)
		grpc, err := service.NewServer(config.Port.Grpc)
		if err != nil {
			application.Log.Fatal("failed to start grpc server: " + err.Error())
		}
		if err := grpc.Start(ctx); err != nil {
			application.Log.Error("failed to start grpc server: " + err.Error())
			os.Exit(1)
		}
		wg.Done()
	}()

	wg.Wait()
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

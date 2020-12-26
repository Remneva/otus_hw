package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/configs"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/server/http"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/storage/sql"
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
	config, err := configs.Read(config)
	if err != nil {
		log.Fatal("failed to read config")
	}
	logg, err := logger.NewLogger(config.Logger.Level, config.Logger.Path)
	if err != nil {
		log.Fatal("failed to create logger")
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	storage := new(sql.Storage)
	if err := storage.Connect(ctx, config.PSQL.DSN, logg); err != nil {
		logg.Fatal("fail connection")
	}
	application := app.New(logg, storage, config)
	if err != nil {
		logg.Fatal("failed to start application")
	}

	wg := sync.WaitGroup{}
	wg.Add(2)
	var http *internalhttp.Server
	var grpc *internalgrpc.Server
	go func() {
		_, mux := internalhttp.NewHandler(ctx, application)
		http, err := internalhttp.NewServer(mux, config.Port.HTTP, logg)
		if err != nil {
			application.Log.Fatal("failed to start http server: " + err.Error())
		}
		if err = http.Start(ctx); err != nil {
			application.Log.Error("failed to start http server: " + err.Error())
			os.Exit(1)
		}
		wg.Done()
	}()

	go func() {
		//service := internalgrpc.New(ctx, application)
		grpc, err := internalgrpc.NewServer(application, config.Port.Grpc)
		if err != nil {
			application.Log.Fatal("failed to start grpc server: " + err.Error())
		}
		if err = grpc.Start(ctx); err != nil {
			application.Log.Error("failed to start grpc server: " + err.Error())
			os.Exit(1)
		}
		wg.Done()
	}()

	wg.Wait()
	go signalChan(application, http, grpc)
}

func signalChan(app *app.App, http *internalhttp.Server, grpc *internalgrpc.Server) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals)

	<-signals
	signal.Stop(signals)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err := http.Stop(ctx)
	if err != nil {
		app.Log.Error(err.Error())
	}
	err = grpc.Stop(ctx)
	if err != nil {
		app.Log.Error(err.Error())
	}
}

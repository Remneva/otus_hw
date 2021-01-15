package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
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
	flag.StringVar(&config, "config", "./configs/config.toml", "Path to configuration file")
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
	var application *app.App
	if !config.Mode.MemMode {
		storage := new(sql.Storage)
		application = app.NewStoreApp(logg, storage, config)
		if err := storage.Connect(ctx, config.PSQL.DSN, logg); err != nil {
			logg.Fatal("fail connection")
		}
	} else {
		application = app.NewMemApp(logg, config)
	}

	var http *internalhttp.Server
	http, err = internalhttp.NewHTTP(ctx, application, config.Port.HTTP)
	if err != nil {
		logg.Fatal("failed to start http server: " + err.Error())
	}
	var grpc *internalgrpc.Server
	grpc, err = internalgrpc.NewServer(application, config.Port.Grpc)
	if err != nil {
		logg.Fatal("failed to start grpc server: " + err.Error())
	}

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		if err = http.Start(); err != nil {
			logg.Error("failed to start http server: " + err.Error())
		}
	}()

	go func() {
		defer wg.Done()
		if err = grpc.Start(); err != nil {
			logg.Error("failed to start grpc server: " + err.Error())
		}
	}()
	go signalChan(application, http, grpc)

	wg.Wait()
}

func signalChan(app *app.App, http *internalhttp.Server, grpc *internalgrpc.Server) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	fmt.Printf("Got %v...\n", <-signals)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err := http.Stop(ctx)
	if err != nil {
		app.Log.Error(err.Error())
	}
	app.Log.Info("http server shutdown")

	err = grpc.Stop()
	if err != nil {
		app.Log.Error(err.Error())
	}
	app.Log.Info("grpc server shutdown")
}

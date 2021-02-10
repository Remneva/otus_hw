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
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/pkg/app"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/pkg/logger"
	srv "github.com/Remneva/otus_hw/hw12_13_14_15_calendar/pkg/server"
	internalgrpc "github.com/Remneva/otus_hw/hw12_13_14_15_calendar/pkg/server/grpc"
	internalhttp "github.com/Remneva/otus_hw/hw12_13_14_15_calendar/pkg/server/http"
	memorystorage "github.com/Remneva/otus_hw/hw12_13_14_15_calendar/pkg/storage/memory"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/pkg/storage/sql"
	"github.com/apex/log"
	"go.uber.org/zap"
)

var config string
var env string

func init() {
	flag.StringVar(&config, "config", "./configs/config.toml", "Path to configuration file")
	flag.StringVar(&env, "env", "prod", "environmental")
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

	logg, err := logger.NewLogger(config.Logger.Level, env, config.Logger.Path)
	if err != nil {
		log.Fatal("failed to create logger")
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var application *app.App
	if !config.Mode.MemMode {
		storage := sql.NewStorage(logg)
		if err := storage.Connect(ctx, config.PSQL.DSN); err != nil {
			logg.Fatal("fail connection")
		}
		application = app.NewApp(logg, storage)
		defer storage.Close()
	} else {
		storage := memorystorage.NewMap(logg)
		application = app.NewApp(logg, storage)
	}
	var http *internalhttp.Server
	http, err = internalhttp.NewHTTP(ctx, application, logg, config.Port.HTTP)
	if err != nil {
		logg.Fatal("failed to start http server: " + err.Error())
	}
	var grpc *internalgrpc.Server
	grpc, err = internalgrpc.NewServer(application, logg, config.Port.Grpc)
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
	go signalChan(logg, http, grpc)

	wg.Wait()
}

func signalChan(log *zap.Logger, srv ...srv.Stopper) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	fmt.Printf("Got %v...\n", <-signals)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	for _, s := range srv {
		err := s.Stop(ctx)
		if err != nil {
			log.Error("failed to stop", zap.Error(err))
		}
	}
}

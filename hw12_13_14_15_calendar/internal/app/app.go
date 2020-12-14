package app

import (
	"context"
	internalgrpc "github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/server/http"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"os"
	"sync"

	sqlstorage "github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/storage/sql"
	"go.uber.org/zap"
)

type App struct {
	Repo   sqlstorage.BaseStorage
	log    *zap.Logger
	server *internalhttp.Server
	grpc   *grpc.Server
	Hand   *internalhttp.MyHandler
}

type Logger interface {
	// TODO
}

type Storage interface {
	// TODO
}

// вывод данных из базы оставила пока для проверки
func (a *App) Run(ctx context.Context) error {
	events, err := a.Repo.GetEvents(ctx)
	if err != nil {
		return errors.New("select query error")
	}
	for _, ev := range events {
		a.log.Info("ev: ", zap.String("Nahuatl name", ev.Title))
	}
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		a.log.Info("http is running...")
		a.server = internalhttp.NewServer(a.log)
		if err := a.server.Start(ctx); err != nil {
			a.log.Error("failed to start http server: " + err.Error())
			os.Exit(1)
		}
		wg.Done()
	}()

	go func() {
		a.log.Info("grpc is running...")
		a.grpc, err = internalgrpc.NewServer(a.log)
		if err != nil {
			a.log.Error("failed to start grpc server: " + err.Error())
		}
		wg.Done()
	}()

	wg.Wait()
	return nil
}

func (a *App) httpServerNew() *internalhttp.Server {
	return internalhttp.NewServer(a.log)
}

func New(logger *zap.Logger, r sqlstorage.BaseStorage) (*App, error) {
	return &App{Repo: r, log: logger}, nil
}

func (a *App) Stop(ctx context.Context) error {
	if err := a.Stop(ctx); err != nil {
		a.log.Error("failed to stop http server: " + err.Error())
		return errors.Wrap(err, "failed to stop http server")
	}
	return nil
}

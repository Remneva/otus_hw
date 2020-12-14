package app

import (
	"context"
	"errors"
	internalgrpc "github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/server/http"
	"google.golang.org/grpc"
	"os"
	"sync"

	sqlstorage "github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/storage/sql"
	"go.uber.org/zap"
)

type App struct {
	r      sqlstorage.BaseStorage
	l      *zap.Logger
	server *internalhttp.Server
	grpc   *grpc.Server
}

type Logger interface {
	// TODO
}

type Storage interface {
	// TODO
}

// вывод данных из базы оставила пока для проверки
func (a *App) Run(ctx context.Context) error {
	events, err := a.r.GetEvents(ctx)
	if err != nil {
		return errors.New("select query error")
	}
	for _, ev := range events {
		a.l.Info("ev: ", zap.String("Nahuatl name", ev.Title))
	}
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		a.server = internalhttp.NewServer(a.l)
		if err := a.server.Start(ctx); err != nil {
			a.l.Error("failed to start http server: " + err.Error())
			os.Exit(1)
		}
		wg.Done()
	}()

	go func() {
		a.grpc, err = internalgrpc.NewServer(a.l)
		if err != nil {
			a.l.Error("failed to start grpc server: " + err.Error())
		}
		a.l.Info("grpc is running...")
		wg.Done()
	}()

	wg.Wait()
	return nil
}

func (a *App) httpServerNew() *internalhttp.Server {
	return internalhttp.NewServer(a.l)
}

func New(logger *zap.Logger, r sqlstorage.BaseStorage) (*App, error) {
	return &App{r: r, l: logger}, nil
}

func (a *App) Stop(ctx context.Context) error {
	if err := a.Stop(ctx); err != nil {
		a.l.Error("failed to stop http server: " + err.Error())
	}
	return nil
}

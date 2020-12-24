package app

import (
	"context"
	sqlstorage "github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type App struct {
	Repo sqlstorage.EventsStorage
	Log  *zap.Logger
}

type Logger interface {
	// TODO
}

type Storage interface {
	// TODO
}

//func (a *App) httpServerNew() *internalhttp.Server {
//	return internalhttp.NewServer(a, a.log)
//}

func New(logger *zap.Logger, r sqlstorage.BaseStorage) (*App, error) {
	return &App{Repo: r, Log: logger}, nil
}

func (a *App) Stop(ctx context.Context) error {
	if err := a.Stop(ctx); err != nil {
		a.Log.Error("failed to stop http server: " + err.Error())
		return errors.Wrap(err, "failed to stop http server")
	}
	return nil
}

package app

import (
	"context"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type App struct {
	Repo storage.EventsStorage
	Log  *zap.Logger
}

type Logger interface {
	// TODO
}

type Storage interface {
	// TODO
}

func New(logger *zap.Logger, r storage.EventsStorage) (*App, error) {
	return &App{Repo: r, Log: logger}, nil
}

func (a *App) Stop(ctx context.Context) error {
	if err := a.Stop(ctx); err != nil {
		a.Log.Error("failed to stop http server: " + err.Error())
		return errors.Wrap(err, "failed to stop http server")
	}
	return nil
}

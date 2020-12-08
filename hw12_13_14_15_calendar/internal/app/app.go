package app

import (
	"context"
	"errors"

	sqlstorage "github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/storage/sql"
	"go.uber.org/zap"
)

type App struct {
	r sqlstorage.BaseStorage
	l *zap.Logger
}

type Logger interface {
	// TODO
}

type Storage interface {
	// TODO
}

func (a *App) Run(ctx context.Context) error {
	events, err := a.r.GetEvents(ctx)
	if err != nil {
		return errors.New("select query error")
	}
	for _, ev := range events {
		a.l.Info("ev: ", zap.String("Nahuatl name", ev.Title))
	}

	return nil
}

func New(logger *zap.Logger, r sqlstorage.BaseStorage) (*App, error) {
	return &App{r: r, l: logger}, nil
}

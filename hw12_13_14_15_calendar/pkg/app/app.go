package app

import (
	"context"
	"fmt"

	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/pkg/storage"
	"go.uber.org/zap"
)

type App struct {
	repo storage.EventsStorage
	log  *zap.Logger
}

var _ Application = (*App)(nil)

type Application interface {
	Create(ctx context.Context, eve storage.Event) (int64, error)
	Update(ctx context.Context, eve storage.Event) error
	Delete(ctx context.Context, id int64) error
	Get(ctx context.Context, id int64) (storage.Event, error)
}

func NewApp(logger *zap.Logger, r storage.EventsStorage) *App {
	return &App{
		repo: r,
		log:  logger,
	}
}

func (a *App) Create(ctx context.Context, eve storage.Event) (int64, error) {
	id, err := a.repo.AddEvent(ctx, eve)
	if err != nil {
		a.log.Info("Create Event method", zap.String("error", err.Error()))
		return 0, fmt.Errorf("create error: %w", err)
	}
	return id, nil
}

func (a *App) Update(ctx context.Context, eve storage.Event) error {
	err := a.repo.UpdateEvent(ctx, eve)
	if err != nil {
		a.log.Info("Update Event psql method", zap.String("error", err.Error()))
		return fmt.Errorf("update error: %w", err)
	}
	return nil
}

func (a *App) Delete(ctx context.Context, id int64) error {
	err := a.repo.DeleteEvent(ctx, id)
	if err != nil {
		a.log.Error("Delete Event psql method", zap.Error(err))
		return fmt.Errorf("delete error: %w", err)
	}
	return nil
}

func (a *App) Get(ctx context.Context, id int64) (storage.Event, error) {
	eve, err := a.repo.GetEvent(ctx, id)
	if err != nil {
		a.log.Error("Get Event psql method", zap.Error(err))
		return eve, fmt.Errorf("get error: %w", err)
	}
	return eve, nil
}

package app

import (
	"context"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/configs"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/storage/memory"
	"go.uber.org/zap"
)

type App struct {
	Mode bool
	Repo storage.EventsStorage
	Mem  *memorystorage.EventMap
	Log  *zap.Logger
}

var _ Application = (*App)(nil)

type Application interface {
	Create(ctx context.Context, eve storage.Event) (int64, error)
	CreateInMemory(ctx context.Context, eve storage.Event) (int64, error)
	Update(ctx context.Context, eve storage.Event) error
	UpdateInMemory(ctx context.Context, eve storage.Event) error
	Delete(ctx context.Context, id int64) error
	DeleteInMemory(ctx context.Context, id int64) error
	Get(ctx context.Context, id int64) (storage.Event, error)
	GetInMemory(ctx context.Context, id int64) (storage.Event, error)
}

func New(logger *zap.Logger, r storage.EventsStorage, c configs.Config) *App {
	if !c.Mode.MemMode {
		return &App{Mem: nil, Repo: r, Log: logger, Mode: c.Mode.MemMode}
	}
	return &App{Mem: memorystorage.NewMap(), Repo: nil, Log: logger, Mode: c.Mode.MemMode}
}

func NewMemApp(logger *zap.Logger, c configs.Config) *App {
	return &App{Mem: memorystorage.NewMap(), Repo: nil, Log: logger, Mode: c.Mode.MemMode}
}

func NewStoreApp(logger *zap.Logger, r storage.EventsStorage, c configs.Config) *App {
	return &App{Mem: nil, Repo: r, Log: logger, Mode: c.Mode.MemMode}
}

func (a *App) Create(ctx context.Context, eve storage.Event) (int64, error) {
	id, err := a.Repo.AddEvent(ctx, eve)
	if err != nil {
		a.Log.Info("Create Event method", zap.String("error", err.Error()))
		return 0, err
	}
	return id, nil
}

func (a *App) CreateInMemory(ctx context.Context, eve storage.Event) (int64, error) {
	id, err := a.Mem.AddEvent(ctx, eve)
	if err != nil {
		a.Log.Info("Create Event memory method", zap.String("error", err.Error()))
		return 0, err
	}
	return id, nil
}

func (a *App) Update(ctx context.Context, eve storage.Event) error {
	err := a.Repo.UpdateEvent(ctx, eve)
	if err != nil {
		a.Log.Info("Update Event psql method", zap.String("error", err.Error()))
		return err
	}
	return nil
}

func (a *App) UpdateInMemory(ctx context.Context, eve storage.Event) error {
	err := a.Mem.UpdateEvent(ctx, eve)
	if err != nil {
		a.Log.Info("Update Event memory method", zap.String("error", err.Error()))
		return err
	}
	return nil
}

func (a *App) Delete(ctx context.Context, id int64) error {
	err := a.Repo.DeleteEvent(ctx, id)
	if err != nil {
		a.Log.Error("Delete Event psql method", zap.Error(err))
		return err
	}
	return nil
}

func (a *App) DeleteInMemory(ctx context.Context, id int64) error {
	err := a.Mem.DeleteEvent(ctx, id)
	if err != nil {
		a.Log.Error("Delete Event memory method", zap.Error(err))
		return err
	}
	return nil
}

func (a *App) Get(ctx context.Context, id int64) (storage.Event, error) {
	eve, err := a.Repo.GetEvent(ctx, id)
	if err != nil {
		a.Log.Error("Delete Event psql method", zap.Error(err))
		return eve, err
	}
	return eve, nil
}

func (a *App) GetInMemory(ctx context.Context, id int64) (storage.Event, error) {
	eve, err := a.Mem.GetEvent(ctx, id)
	if err != nil {
		a.Log.Error("Delete Event memory method", zap.Error(err))
		return eve, err
	}
	return eve, nil
}

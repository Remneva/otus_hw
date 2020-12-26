package app

import (
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

func New(logger *zap.Logger, r storage.EventsStorage, c configs.Config) *App {
	if !c.Mode.MemMode {
		return &App{Mem: nil, Repo: r, Log: logger, Mode: c.Mode.MemMode}
	}
	return &App{Mem: memorystorage.NewMap(), Repo: nil, Log: logger, Mode: c.Mode.MemMode}
}

package storage

import (
	"context"
	"time"

	// Postgres driver.
	_ "github.com/jackc/pgx/v4/stdlib"
	"go.uber.org/zap"
)

// go get github.com/golang/mock/mockgen
//go:generate mockgen -destination=mock_db_test.go -package=storage . EventsStorage

type BaseStorage interface {
	Connect(ctx context.Context, dsn string, l *zap.Logger) error
	Close(ctx context.Context) error
	EventsStorage
}

type EventsStorage interface {
	GetEvents(ctx context.Context) ([]Event, error)
	GetEvent(ctx context.Context, ID int64) (Event, error)
	AddEvent(ctx context.Context, ev Event) (int64, error)
	DeleteEvent(ctx context.Context, ID int64) error
	UpdateEvent(ctx context.Context, ev Event) error
}

type Event struct {
	ID          int64
	Owner       int64
	Title       string
	Description string
	StartDate   string
	StartTime   time.Time
	EndDate     string
	EndTime     time.Time
}

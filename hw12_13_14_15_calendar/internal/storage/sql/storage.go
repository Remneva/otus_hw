package sqlstorage //nolint:golint,stylecheck

import (
	"context"
	"time"

	// Postgres driver.
	_ "github.com/jackc/pgx/v4/stdlib"
)

type BaseStorage interface {
	Connect(ctx context.Context, dsn string) error
	Close(ctx context.Context) error
	EventsStorage
}

type EventsStorage interface {
	GetEvents(ctx context.Context) ([]Event, error)
	GetEvent(ctx context.Context, ID int64) (Event, error)
	SetEvent(ctx context.Context, title string, descr string, startDate time.Time, startTime time.Time, endDate time.Time, endTime time.Time) error
	DeleteEvent(ctx context.Context, ID int64) error
	UpdateEvent(ctx context.Context, FieldToChange string, NewValue interface{}, ID int64) (Event, error)
	CreateEvent(ctx context.Context, ev Event) error
}

type Event struct {
	ID        int64
	Owner     int64
	Title     string
	Descr     string
	StartDate time.Time
	StartTime string
	EndDate   time.Time
	EndTime   string
}

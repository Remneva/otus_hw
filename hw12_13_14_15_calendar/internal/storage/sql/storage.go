package sqlstorage

import (
	"context"
	_ "github.com/jackc/pgx/v4/stdlib"
	"time"
)

type BaseStorage interface {
	Connect(ctx context.Context, dsn string) error
	Close(ctx context.Context) error
	EventsStorage
}

type EventsStorage interface {
	GetEvents(ctx context.Context) ([]Event, error)
	GetEvent(ctx context.Context, Id int64) (Event, error)
	SetEvent(ctx context.Context, title string, descr string, start_date time.Time, start_time time.Time, end_date time.Time, end_time time.Time) error
	DeleteEvent(ctx context.Context, Id int64) error
	UpdateEvent(ctx context.Context, FieldToChange string, NewValue interface{}, Id int64) (Event, error)
	CreateEvent(ctx context.Context, ev Event) error
}

type Event struct {
	Id        int64
	Owner     int64
	Title     string
	Descr     string
	StartDate time.Time
	StartTime string
	EndDate   time.Time
	EndTime   string
}

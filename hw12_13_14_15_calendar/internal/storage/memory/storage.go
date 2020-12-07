package memorystorage

import (
	"context"
	"database/sql"
	"fmt"
	sqlstorage "github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/apex/log"
	"go.uber.org/zap"
	"strconv"
	"sync"
	"time"
)

var _ sqlstorage.BaseStorage = (*Storage)(nil)

type Storage struct {
	db *sql.DB
	mu sync.RWMutex
	l  *zap.Logger
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Connect(ctx context.Context, dsn string) (err error) {
	fmt.Println("dsn: ", dsn)
	s.db, err = sql.Open("pgx", dsn)
	if err != nil {
		return
	}
	s.db.Stats()
	return s.db.PingContext(ctx)
}

func (s *Storage) Close(ctx context.Context) error {
	return s.db.Close()
}

func (s *Storage) DeleteEvent(ctx context.Context, Id int64) error {
	row, err := s.db.ExecContext(ctx, `
		DELETE from events where Id = $1
		`, Id)
	if err != nil {
		return err
	}
	rowAffected, _ := row.RowsAffected()
	log.Debug(strconv.FormatInt(rowAffected, 10))
	return nil
}

func (s *Storage) UpdateEvent(ctx context.Context, FieldToChange string, NewValue interface{}, Id int64) (sqlstorage.Event, error) {
	var ev sqlstorage.Event
	row, err := s.db.QueryContext(ctx, `
		Update events 
		set $1 = $2
		where Id = $3
		`, FieldToChange, NewValue, Id)
	if err != nil {
		return ev, err
	}
	defer row.Close()
	err = row.Scan(
		&ev.Id,
		&ev.Title,
		&ev.Owner,
		&ev.StartDate,
		&ev.StartTime,
		&ev.EndDate,
		&ev.EndTime)

	if err != nil {
		return ev, err
		log.Debug("cant't update")
	}
	return ev, row.Err()
}

func (s *Storage) GetEvent(ctx context.Context, Id int64) (sqlstorage.Event, error) {
	var ev sqlstorage.Event
	row := s.db.QueryRowContext(ctx, `
		SELECT title, descr, start_date, start_time, end_date, end_time FROM events where Id = $1`, Id)
	err := row.Scan(
		&ev.Id,
		&ev.Title,
		&ev.Owner,
		&ev.StartDate,
		&ev.StartTime,
		&ev.EndDate,
		&ev.EndTime)

	if err != nil {
		return ev, err
	}
	if err == sql.ErrNoRows {
		return ev, nil
	} else if err != nil {
		return ev, err
	}

	return ev, row.Err()
}

func (s *Storage) GetEvents(ctx context.Context) ([]sqlstorage.Event, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, owner, title, descr, start_date, start_time, end_date, end_time FROM events
	`)
	if err != nil {
		fmt.Println("query error:", err)
		return nil, err
	}
	defer rows.Close()

	var events []sqlstorage.Event

	for rows.Next() {
		var ev sqlstorage.Event

		if err := rows.Scan(
			&ev.Id,
			&ev.Owner,
			&ev.Title,
			&ev.Descr,
			&ev.StartDate,
			&ev.StartTime,
			&ev.EndDate,
			&ev.EndTime,
		); err != nil {
			return nil, err
		}

		events = append(events, ev)
	}
	return events, rows.Err()
}

func (s *Storage) SetEvent(ctx context.Context, title string, descr string, start_date time.Time, start_time time.Time, end_date time.Time, end_time time.Time) error {
	query := `INSERT INTO events (title, descr, start_date, start_time, end_date, end_time)
VALUES($1, $2, $3, $4, $5, $6)`
	row, err := s.db.ExecContext(ctx, query, title, descr, start_date, start_time, end_date, end_time)
	if err != nil {
		return err
	}
	rowAffected, _ := row.RowsAffected()
	log.Debug(strconv.FormatInt(rowAffected, 10))
	return nil
}

func (s *Storage) CreateEvent(ctx context.Context, ev sqlstorage.Event) error {
	owner := ev.Owner
	title := ev.Title
	descr := ev.Descr
	start_date := ev.StartDate
	start_time := ev.StartTime
	end_date := ev.EndDate
	end_time := ev.EndTime
	query := `INSERT INTO events (title, descr, start_date, start_time, end_date, end_time)
VALUES($1, $2, $3, $4, $5, $6, $7)`
	row, err := s.db.ExecContext(ctx, query, owner, title, descr, start_date, start_time, end_date, end_time)
	if err != nil {
		return err
	}
	rowAffected, _ := row.RowsAffected()
	log.Debug(strconv.FormatInt(rowAffected, 10))
	return nil
}

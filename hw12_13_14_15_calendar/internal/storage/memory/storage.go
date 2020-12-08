package memorystorage

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"time"

	sqlstorage "github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/apex/log"
	"go.uber.org/zap"
)

var ErrConnection = errors.New("connection error")
var _ sqlstorage.BaseStorage = (*Storage)(nil)

type Storage struct {
	db *sql.DB
	l  *zap.Logger
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Connect(ctx context.Context, dsn string) (err error) {
	s.db, err = sql.Open("pgx", dsn)
	if err != nil {
		s.l.Error("Error", zap.String("Driver", err.Error()))
		return ErrConnection
	}
	s.db.Stats()
	return s.db.PingContext(ctx)
}

func (s *Storage) Close(ctx context.Context) error {
	return s.db.Close()
}

func (s *Storage) DeleteEvent(ctx context.Context, id int64) error {
	row, err := s.db.ExecContext(ctx, `
		DELETE from events where ID = $1
		`, id)
	if err != nil {
		s.l.Error("Error", zap.String("Connection", err.Error()))
		return ErrConnection
	}
	rowAffected, _ := row.RowsAffected()
	log.Debug(strconv.FormatInt(rowAffected, 10))
	return nil
}

func (s *Storage) UpdateEvent(ctx context.Context, fieldToChange string, newValue interface{}, id int64) (sqlstorage.Event, error) {
	var ev sqlstorage.Event
	row, err := s.db.QueryContext(ctx, `
		Update events 
		set $1 = $2
		where ID = $3
		`, fieldToChange, newValue, id)
	if err != nil {
		return ev, ErrConnection
	}
	defer row.Close()
	err = row.Scan(
		&ev.ID,
		&ev.Title,
		&ev.Owner,
		&ev.StartDate,
		&ev.StartTime,
		&ev.EndDate,
		&ev.EndTime)

	if err != nil {
		s.l.Error("Update error", zap.String("query", row.Err().Error()))
		return ev, ErrConnection
	}
	return ev, row.Err()
}

func (s *Storage) GetEvent(ctx context.Context, id int64) (sqlstorage.Event, error) {
	var ev sqlstorage.Event
	row := s.db.QueryRowContext(ctx, `
		SELECT title, descr, start_date, start_time, end_date, end_time FROM events where ID = $1`, id)
	err := row.Scan(
		&ev.ID,
		&ev.Title,
		&ev.Owner,
		&ev.StartDate,
		&ev.StartTime,
		&ev.EndDate,
		&ev.EndTime)

	if err.Error() == sql.ErrNoRows.Error() {
		return ev, nil
	} else if err != nil {
		return ev, ErrConnection
	}

	return ev, ErrConnection
}

func (s *Storage) GetEvents(ctx context.Context) ([]sqlstorage.Event, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, owner, title, descr, start_date, start_time, end_date, end_time FROM events
	`)
	if err != nil {
		return nil, ErrConnection
	}
	defer rows.Close()

	var events []sqlstorage.Event

	for rows.Next() {
		var ev sqlstorage.Event

		if err := rows.Scan(
			&ev.ID,
			&ev.Owner,
			&ev.Title,
			&ev.Descr,
			&ev.StartDate,
			&ev.StartTime,
			&ev.EndDate,
			&ev.EndTime,
		); err != nil {
			s.l.Error("Get error", zap.String("query", rows.Err().Error()))
			return nil, ErrConnection
		}

		events = append(events, ev)
	}
	return events, rows.Err()
}

func (s *Storage) SetEvent(ctx context.Context, title string, descr string, startDate time.Time, startTime time.Time, endDate time.Time, endTime time.Time) error {
	query := `INSERT INTO events (title, descr, start_date, start_time, end_date, end_time)
VALUES($1, $2, $3, $4, $5, $6)`
	row, err := s.db.ExecContext(ctx, query, title, descr, startDate, startTime, endDate, endTime)
	if err != nil {
		return ErrConnection
	}
	rowAffected, _ := row.RowsAffected()
	log.Debug(strconv.FormatInt(rowAffected, 10))
	return nil
}

func (s *Storage) CreateEvent(ctx context.Context, ev sqlstorage.Event) error {
	owner := ev.Owner
	title := ev.Title
	descr := ev.Descr
	startDate := ev.StartDate
	startTime := ev.StartTime
	endDate := ev.EndDate
	endTime := ev.EndTime
	query := `INSERT INTO events (title, descr, start_date, start_time, end_date, end_time)
VALUES($1, $2, $3, $4, $5, $6, $7)`
	row, err := s.db.ExecContext(ctx, query, owner, title, descr, startDate, startTime, endDate, endTime)
	if err != nil {
		return ErrConnection
	}
	rowAffected, _ := row.RowsAffected()
	log.Debug(strconv.FormatInt(rowAffected, 10))
	return nil
}

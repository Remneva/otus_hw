package sql

import (
	"context"
	"database/sql"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/storage"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var _ storage.BaseStorage = (*Storage)(nil)
var OpenConnectionErr = errors.New("Open connection error")
var SqlQueryErr = errors.New("Database query failed")

type Storage struct {
	db *sql.DB
	l  *zap.Logger
	s  storage.EventsStorage
}

func New(s storage.EventsStorage) *Storage {
	return &Storage{s: s}
}

func (s *Storage) Connect(ctx context.Context, dsn string, l *zap.Logger) (err error) {
	s.l = l
	s.db, err = sql.Open("pgx", dsn)
	if err != nil {
		s.l.Error("Error", zap.String("Open connection", err.Error()))
		return errors.Wrapf(OpenConnectionErr, err.Error())
	}
	err = s.db.PingContext(ctx)
	if err != nil {
		s.l.Error("Error", zap.String("Ping", err.Error()))
		return errors.Wrap(err, "Ping error")
	}
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	return s.db.Close()
}

func (s *Storage) DeleteEvent(ctx context.Context, id int64) error {
	_, err := s.db.ExecContext(ctx, `
		DELETE from events where ID = $1
		`, id)
	if err != nil {
		s.l.Error("Error", zap.String("Connection", err.Error()))
		return errors.Wrapf(SqlQueryErr, err.Error())
	}
	s.l.Info("Event Deleted", zap.Int64("id", id))
	return nil
}

func (s *Storage) UpdateEvent(ctx context.Context, ev storage.Event) error {
	query := "Update events SET owner = $1, title = $2, description = $3, start_date = $4, start_time = $5, end_date = $6, end_time = $7 WHERE id = $8 "
	_, err := s.db.ExecContext(ctx, query, ev.Owner, ev.Title, ev.Description, ev.StartDate, ev.StartTime, ev.EndDate, ev.EndTime, ev.ID)
	if err != nil {
		s.l.Error("Exec query error", zap.Error(err))
		return errors.Wrapf(SqlQueryErr, err.Error())
	}
	s.l.Info("Event Updated", zap.Int64("id", ev.ID))
	return nil
}

func (s *Storage) GetEvent(ctx context.Context, id int64) (storage.Event, error) {
	var ev storage.Event
	row := s.db.QueryRowContext(ctx, `
		SELECT id, owner, title, description, start_date, start_time, end_date, end_time FROM events where id = $1`, id)
	err := row.Scan(
		&ev.ID,
		&ev.Owner,
		&ev.Title,
		&ev.Description,
		&ev.StartDate,
		&ev.StartTime,
		&ev.EndDate,
		&ev.EndTime)
	if err != nil {
		return ev, errors.Wrapf(SqlQueryErr, err.Error())
	}
	return ev, nil
}

func (s *Storage) GetEvents(ctx context.Context) ([]storage.Event, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, owner, title, description, start_date, start_time, end_date, end_time FROM events
	`)
	if err != nil {
		return nil, errors.Wrapf(SqlQueryErr, err.Error())
	}
	defer rows.Close()
	var events []storage.Event
	for rows.Next() {
		var ev storage.Event
		if err := rows.Scan(
			&ev.ID,
			&ev.Owner,
			&ev.Title,
			&ev.Description,
			&ev.StartDate,
			&ev.StartTime,
			&ev.EndDate,
			&ev.EndTime,
		); err != nil {
			s.l.Error("Get event error", zap.String("query", rows.Err().Error()))
			return nil, errors.Wrapf(SqlQueryErr, err.Error())
		}
		events = append(events, ev)
	}
	return events, rows.Err()
}

func (s *Storage) AddEvent(ctx context.Context, ev storage.Event) (int64, error) {
	query := `INSERT INTO events (owner, title, description, start_date, start_time, end_date, end_time)
VALUES($1, $2, $3, $4, $5, $6, $7)`
	_, err := s.db.ExecContext(ctx, query, ev.Owner, ev.Title, ev.Description, ev.StartDate, ev.StartTime, ev.EndDate, ev.EndTime)
	if err != nil {
		return 0, errors.Wrapf(SqlQueryErr, err.Error())
	}
	var id int64
	err = s.db.QueryRowContext(ctx, `
		SELECT id FROM events ORDER BY id DESC LIMIT 1`).Scan(&id)
	if err != nil {
		return 0, errors.Wrapf(SqlQueryErr, err.Error())
	}
	s.l.Info("Event Created", zap.Int64("id", id))
	return id, nil
}

package memorystorage

import (
	"context"
	"database/sql"

	sqlstorage "github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var _ sqlstorage.BaseStorage = (*Storage)(nil)

type Storage struct {
	db *sql.DB
	l  *zap.Logger
	s  sqlstorage.EventsStorage
}

func New(s sqlstorage.EventsStorage) *Storage {
	return &Storage{s: s}
}

func (s *Storage) Connect(ctx context.Context, dsn string, l *zap.Logger) (err error) {
	s.l = l
	s.db, err = sql.Open("pgx", dsn)
	if err != nil {
		s.l.Error("Error", zap.String("Open connection", err.Error()))
		return errors.Wrap(err, "Open connection error")
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
		return errors.Wrap(err, "Database query failed")
	}
	s.l.Info("Event Deleted", zap.Int64("id", id))
	return nil
}

func (s *Storage) UpdateEvent(ctx context.Context, ev sqlstorage.Event) error {
	query := "Update events SET owner = $1, title = $2, description = $3, start_date = $4, start_time = $5, end_date = $6, end_time = $7 WHERE id = $8 "
	_, err := s.db.ExecContext(ctx, query, ev.Owner, ev.Title, ev.Description, ev.StartDate, ev.StartTime, ev.EndDate, ev.EndTime, ev.ID)
	if err != nil {
		s.l.Error("Exec query error", zap.Error(err))
		return errors.Wrap(err, "Database query failed")
	}
	s.l.Info("Event Updated", zap.Int64("id", ev.ID))
	return nil
}

func (s *Storage) GetEvent(ctx context.Context, id int64) (sqlstorage.Event, error) {
	var ev sqlstorage.Event
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
		return ev, errors.Wrap(err, "Database query failed")
	}
	return ev, nil
}

func (s *Storage) GetEvents(ctx context.Context) ([]sqlstorage.Event, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, owner, title, description, start_date, start_time, end_date, end_time FROM events
	`)
	if err != nil {
		return nil, errors.Wrap(err, "Database query failed")
	}
	defer rows.Close()
	var events []sqlstorage.Event
	for rows.Next() {
		var ev sqlstorage.Event
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
			return nil, errors.Wrap(err, "Database query failed")
		}
		events = append(events, ev)
	}
	return events, rows.Err()
}

func (s *Storage) AddEvent(ctx context.Context, ev sqlstorage.Event) (int64, error) {
	err := s.Insert(ctx, ev)
	if err != nil {
		return 0, errors.Wrap(err, "Database query failed")
	}
	id, err := s.GetLastId(ctx)
	if err != nil {
		return 0, errors.Wrap(err, "Getting the last id error")
	}
	s.l.Info("Event Created", zap.Int64("id", id))
	return id, nil
}

func (s *Storage) Insert(ctx context.Context, ev sqlstorage.Event) error {
	query := `INSERT INTO events (owner, title, description, start_date, start_time, end_date, end_time)
VALUES($1, $2, $3, $4, $5, $6, $7)`
	_, err := s.db.ExecContext(ctx, query, ev.Owner, ev.Title, ev.Description, ev.StartDate, ev.StartTime, ev.EndDate, ev.EndTime)
	if err != nil {
		return errors.Wrap(err, "Database query failed")
	}
	return nil
}

func (s *Storage) GetLastId(ctx context.Context) (int64, error) {
	var id int64
	err := s.db.QueryRowContext(ctx, `
		SELECT id FROM events ORDER BY id DESC LIMIT 1`).Scan(&id)
	if err != nil {
		return 0, errors.Wrap(err, "Getting the last id error")
	}
	return id, nil
}

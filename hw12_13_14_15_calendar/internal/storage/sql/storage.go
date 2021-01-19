package sql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/storage"

	// Postgres driver.
	_ "github.com/jackc/pgx/v4/stdlib"
	"go.uber.org/zap"
)

var _ storage.BaseStorage = (*Storage)(nil)

type Storage struct {
	db *sql.DB
	l  *zap.Logger
	storage.EventsStorage
}

func NewStorage(l *zap.Logger) *Storage {
	s := &Storage{
		l: l}
	return s
}
func (s *Storage) Connect(ctx context.Context, dsn string) (err error) {
	s.db, err = sql.Open("pgx", dsn)
	if err != nil {
		s.l.Error("Error", zap.String("Open connection", err.Error()))
		return fmt.Errorf("open connection error %w", err)
	}
	err = s.db.PingContext(ctx)
	if err != nil {
		s.l.Error("Error", zap.String("Ping", err.Error()))
		return fmt.Errorf("ping error: %w", err)
	}
	return nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) DeleteEvent(ctx context.Context, id int64) error {
	exist, err := s.eventExistValidationByID(id)
	if err != nil {
		return fmt.Errorf("SELECT query error %w", err)
	}
	if !exist {
		return fmt.Errorf("event does not exist %w", err)
	}
	_, err = s.db.ExecContext(ctx, `
		DELETE from events where ID = $1
		`, id)
	if err != nil {
		s.l.Error("Error", zap.String("Connection", err.Error()))
		return fmt.Errorf("open connection error %w", err)
	}
	s.l.Info("Event Deleted", zap.Int64("id", id))
	return nil
}

func (s *Storage) UpdateEvent(ctx context.Context, ev storage.Event) error {
	exist, err := s.eventExistValidationByID(ev.ID)
	if err != nil {
		return fmt.Errorf("SELECT query error %w", err)
	}
	if !exist {
		return fmt.Errorf("event does not exist %w", err)
	}
	query := "Update events SET owner = $1, title = $2, description = $3, start_date = $4, start_time = $5, end_date = $6, end_time = $7 WHERE id = $8 "
	result, err := s.db.ExecContext(ctx, query, ev.Owner, ev.Title, ev.Description, ev.StartDate, ev.StartTime, ev.EndDate, ev.EndTime, ev.ID)
	if err != nil {
		s.l.Error("Exec query error", zap.Error(err))
		return fmt.Errorf("open connection error %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected error %w", err)
	}
	if rowsAffected > 0 {
		s.l.Info("Event updated", zap.Int64("id", ev.ID))
	} else {
		s.l.Info("Event does not exist", zap.Int64("id", ev.ID))
		return fmt.Errorf("event does not exist")
	}

	return nil
}

func (s *Storage) GetEvent(ctx context.Context, id int64) (storage.Event, error) {
	var ev storage.Event
	exist, err := s.eventExistValidationByID(id)
	if err != nil {
		return ev, fmt.Errorf("SELECT query error %w", err)
	}
	if !exist {
		return ev, fmt.Errorf("event does not exist %w", err)
	}
	row := s.db.QueryRowContext(ctx, `
		SELECT id, owner, title, description, start_date, start_time, end_date, end_time FROM events where id = $1`, id)
	err = row.Scan(
		&ev.ID,
		&ev.Owner,
		&ev.Title,
		&ev.Description,
		&ev.StartDate,
		&ev.StartTime,
		&ev.EndDate,
		&ev.EndTime)
	if err != nil {
		return ev, fmt.Errorf("query error %w", err)
	}
	return ev, nil
}

func (s *Storage) GetEvents(ctx context.Context) ([]storage.Event, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, owner, title, description, start_date, start_time, end_date, end_time FROM events
	`)
	if err != nil {
		return nil, fmt.Errorf("open connection error %w", err)
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
			return nil, fmt.Errorf("open connection error %w", err)
		}
		events = append(events, ev)
	}
	return events, rows.Err()
}

func (s *Storage) AddEvent(ctx context.Context, ev storage.Event) (int64, error) {
	exist, err := s.eventExistValidation(ev.Owner, ev.StartTime)
	if err != nil {
		return 0, fmt.Errorf("SELECT query error %w", err)
	}
	if exist {
		return 0, fmt.Errorf("event already exist at this time")
	}
	query := `INSERT INTO events (owner, title, description, start_date, start_time, end_date, end_time)
VALUES($1, $2, $3, $4, $5, $6, $7)`
	_, err = s.db.ExecContext(ctx, query, ev.Owner, ev.Title, ev.Description, ev.StartDate, ev.StartTime, ev.EndDate, ev.EndTime)
	if err != nil {
		return 0, fmt.Errorf("open connection error %w", err)
	}
	var id int64
	err = s.db.QueryRowContext(ctx, `
		SELECT id FROM events ORDER BY id DESC LIMIT 1`).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("open connection error %w", err)
	}
	s.l.Info("Event Created", zap.Int64("id", id))
	return id, nil
}

func (s *Storage) eventExistValidationByID(id int64) (bool, error) {
	var exists bool
	row := s.db.QueryRow("SELECT EXISTS(SELECT * FROM events WHERE id = $1)", id)
	if err := row.Scan(&exists); err != nil {
		s.l.Error("Select query error", zap.Error(err))
		return exists, fmt.Errorf("query error %w", err)
	}
	return exists, nil
}

func (s *Storage) eventExistValidation(owner int64, time time.Time) (bool, error) {
	var exists bool
	row := s.db.QueryRow("SELECT EXISTS(SELECT * FROM events WHERE owner = $1 AND start_time = $2)", owner, time)
	if err := row.Scan(&exists); err != nil {
		s.l.Error("Select query error", zap.Error(err))
		return exists, fmt.Errorf("query error %w", err)
	}
	return exists, nil
}

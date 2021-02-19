package memorystorage

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/pkg/storage"
	"go.uber.org/zap"
)

var ErrNoSuchEvent = errors.New("no such event")

func NewMap(l *zap.Logger) *EventMap {
	eve := &EventMap{log: l}
	eve.ev = make(map[int]storage.Event)
	return eve
}

var id int
var _ storage.EventsStorage = (*EventMap)(nil)

type EventMap struct {
	ev  map[int]storage.Event
	mu  sync.Mutex
	log *zap.Logger
}

func (e *EventMap) GetEventsByPeriod(ctx context.Context, starttime time.Time, endtime time.Time) ([]storage.Event, error) {
	panic("implement me")
}

func (e *EventMap) ChangeStateByID(ctx context.Context, id int64) error {
	panic("implement me")
}

func (e *EventMap) GetEvents(ctx context.Context) ([]storage.Event, error) {
	count := len(e.ev)
	slice := make([]storage.Event, count)
	e.mu.Lock()
	defer e.mu.Unlock()
	for _, event := range e.ev {
		slice = append(slice, event) //nolint:makezero
	}
	return slice, nil
}

func (e *EventMap) GetEvent(ctx context.Context, id int64) (storage.Event, error) {
	var ev storage.Event
	e.mu.Lock()
	defer e.mu.Unlock()
	if _, ok := e.ev[int(id)]; ok {
		ev = e.ev[int(id)]
		e.log.Info("getting event from memory", zap.Int64("id", id))
		return ev, nil
	}
	e.log.Error("No such event in memory", zap.Int64("id", id))
	return ev, ErrNoSuchEvent
}

func (e *EventMap) AddEvent(ctx context.Context, ev storage.Event) (int64, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	id = len(e.ev)
	id++
	ev.ID = int64(id)
	e.ev[id] = ev
	e.log.Info("create event in memory", zap.Int("id", id))
	return int64(id), nil
}

func (e *EventMap) DeleteEvent(ctx context.Context, id int64) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if _, ok := e.ev[int(id)]; ok {
		delete(e.ev, int(id))
		e.log.Info("delete event from memory", zap.Int64("id", id))
		return nil
	}
	e.log.Error("No such event in memory", zap.Int64("id", id))
	return ErrNoSuchEvent
}

func (e *EventMap) UpdateEvent(ctx context.Context, ev storage.Event) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if _, ok := e.ev[int(ev.ID)]; ok {
		ev = e.ev[id]
		return nil
	}
	return ErrNoSuchEvent
}

func (e *EventMap) ChangeStatusByID(ctx context.Context, id int64) error {
	panic("implement me")
}

func (e *EventMap) GetStatusByID(ctx context.Context, id int64) (int64, error) {
	panic("implement me")
}

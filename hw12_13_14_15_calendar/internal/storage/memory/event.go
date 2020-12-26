package memorystorage

import (
	"context"
	"errors"

	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/storage"
)

var ErrNoSuchEvent = errors.New("no such event")

func NewMap() *EventMap {
	eve := &EventMap{}
	eve.ev = make(map[int]storage.Event)
	return eve
}

var id int
var _ storage.EventsStorage = (*EventMap)(nil)

type EventMap struct {
	ev map[int]storage.Event
}

func (e *EventMap) GetEvents(ctx context.Context) ([]storage.Event, error) {
	count := len(e.ev)
	slice := make([]storage.Event, count)
	for _, event := range e.ev {
		slice = append(slice, event)
	}
	return slice, nil
}

func (e *EventMap) GetEvent(ctx context.Context, id int64) (storage.Event, error) {
	var ev storage.Event

	if _, ok := e.ev[int(id)]; ok {
		ev = e.ev[int(id)]
		return ev, nil
	}
	return ev, ErrNoSuchEvent
}

func (e *EventMap) AddEvent(ctx context.Context, ev storage.Event) (int64, error) {
	id++
	e.ev[id] = ev
	ev.ID = int64(id)
	return int64(id), nil
}

func (e *EventMap) DeleteEvent(ctx context.Context, id int64) error {
	if _, ok := e.ev[int(id)]; ok {
		delete(e.ev, int(id))
		return nil
	}
	return ErrNoSuchEvent
}

func (e *EventMap) UpdateEvent(ctx context.Context, ev storage.Event) error {
	if _, ok := e.ev[int(ev.ID)]; ok {
		ev = e.ev[id]
		return nil
	}
	return ErrNoSuchEvent
}

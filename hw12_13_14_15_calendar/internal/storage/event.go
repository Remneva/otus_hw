//nolint
package storage

import (
	"errors"
	"sync"
	"time"
)

var ErrNoSuchEvent = errors.New("no such event")

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

var id int

type eventMap struct {
	mu sync.Mutex
	ev map[int]Event
}

func (e *eventMap) Add(ev Event) error {
	id++
	e.ev[id] = ev
	return nil
}

func (e *eventMap) Get(id int) (Event, error) {
	var ev Event
	if _, ok := e.ev[id]; ok {
		ev = e.ev[id]
		return ev, nil
	}
	return ev, ErrNoSuchEvent
}

type EventMap interface {
	Add(ev Event) error
	Get(id int) (Event, error)
}

func NewMap() EventMap {
	eve := &eventMap{}
	eve.ev = make(map[int]Event)
	return eve
}

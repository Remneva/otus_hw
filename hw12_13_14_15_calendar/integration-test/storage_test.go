package main

import (
	store "/hw12_13_14_15_calendar/pkg/storage"
	"fmt"
	"time"
)

func (s *suiteTestIntegration) TestStorage1() {

	var ev store.Event

	id, err := s.s.AddEvent(s.ctx, ev)
	s.Require().Errorf(err, "Database query failed")
	s.Require().Equal(id, int64(0))

	err = s.s.UpdateEvent(s.ctx, ev)
	s.Require().Errorf(err, "Database query failed")

	err = s.s.DeleteEvent(s.ctx, 1000)
	s.Require().Errorf(err, "event does not exist")

}

func (s *suiteTestIntegration) TestStorage2() {

	ev := store.Event{
		Owner:       1,
		Title:       "test",
		Description: "test",
		StartDate:   "2020-03-01",
		StartTime:   time.Now().Add((15) * time.Minute),
		EndDate:     "2020-03-02",
		EndTime:     time.Now().Add((5) * time.Minute),
	}

	id, err := s.s.AddEvent(s.ctx, ev)
	s.Require().NoError(err)

	ev = store.Event{
		ID:          id,
		Owner:       2,
		Title:       "test test",
		Description: "test test",
		StartDate:   "2020-03-01",
		StartTime:   time.Now().Add((15) * time.Minute),
		EndDate:     "2020-03-02",
		EndTime:     time.Now().Add((5) * time.Minute),
	}

	err = s.s.UpdateEvent(s.ctx, ev)
	s.Require().NoError(err)

	event, err := s.s.GetEvent(s.ctx, id)
	s.Require().NoError(err)
	s.Require().Equal(int64(2), event.Owner)
	s.Require().Equal("test test", event.Title)

	err = s.s.DeleteEvent(s.ctx, id)
	s.Require().NoError(err)

	_, err = s.s.GetEvent(s.ctx, id)
	s.Require().Equal("event does not exist", err.Error())
}

func (s *suiteTestIntegration) TestStorage3() {
	startTime := time.Now().Add((-1000) * time.Minute)
	endTime := time.Now().Add((5) * time.Minute)

	events, err := s.s.GetEventsByPeriod(s.ctx, startTime, endTime)
	for _, event := range events {
		fmt.Println(event)
	}
	result := len(events)

	s.Require().Equal(5, result)
	s.Require().NoError(err)
}

func (s *suiteTestIntegration) TestStorage4() {
	startTime := time.Now().Add((168) * time.Hour)
	endTime := time.Now().Add((336) * time.Hour)

	events, err := s.s.GetEventsByPeriod(s.ctx, startTime, endTime)
	result := len(events)

	s.Require().Equal(0, result)
	s.Require().NoError(err)
}

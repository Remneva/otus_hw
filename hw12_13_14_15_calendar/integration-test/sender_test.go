package main

import (
	store "github.com/Remneva/otus_hw/hw12_13_14_15_calendar/pkg/storage"
	"time"
)

func (s *suiteTestIntegration) TestSender() {

	ev := store.Event{
		Owner:       0,
		Title:       "test",
		Description: "test",
		StartDate:   "2020-03-01",
		StartTime:   time.Now().Add((-15) * time.Minute),
		EndDate:     "2020-03-02",
		EndTime:     time.Now().Add((-5) * time.Minute),
	}

	id, err := s.s.AddEvent(s.ctx, ev)
	s.Require().NoError(err)

	status, err := s.s.GetStatusByID(s.ctx, id)
	s.Require().NoError(err)
	s.Require().Equal(int64(0), status)

	event, err := s.s.GetEvent(s.ctx, id)
	s.Require().NoError(err)

	err = s.r.Publish(event)
	s.Require().NoError(err)

	time.Sleep(4 * time.Second)

	status, err = s.s.GetStatusByID(s.ctx, id)
	s.Require().NoError(err)
	s.Require().Equal(int64(1), status)

}

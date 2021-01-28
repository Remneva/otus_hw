package main

import (
	"time"

	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/rabbit"
	store "github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/go-co-op/gocron"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

func cronStart(storage *sql.Storage, q rabbit.Rabbit) {
	s := gocron.NewScheduler(time.UTC)
	_, err := s.Every(q.C.AMQP.Timeout).Second().Do(checkEventsTask, storage, q)
	if err != nil {
		q.Log.Error("cron execution error", zap.Error(err))
		s.Stop()
	}
	s.StartAsync()

	<-q.Done
	s.Clear()
	s.Stop()
	q.Log.Info("go cron stopping...")
}

func checkEventsTask(storage *sql.Storage, q rabbit.Rabbit) {
	start := time.Now()

	oneYearLater := start.AddDate(-1, 0, 0)

	q.Conn, _ = amqp.Dial(q.C.AMQP.URI)
	q.Channel, _ = q.Conn.Channel()
	events, err := storage.GetEvents(q.Ctx)
	if err != nil {
		q.Log.Error("fail request", zap.Error(err))
	}

	for _, ev := range events {
		if ev.StartTime.After(time.Now().Add((-30) * time.Minute)) {
			go func(ev store.Event) {
				if err := q.Publish(ev); err != nil {
					q.Log.Error("publish message error", zap.Error(err))
				}
			}(ev)
		} else if ev.StartTime.Before(oneYearLater) {
			err = storage.DeleteEvent(q.Ctx, ev.ID)
			if err != nil {
				q.Log.Error("delete failed", zap.Int("id", int(ev.ID)))
			}
			q.Log.Info("outdated event deleted", zap.Int("id", int(ev.ID)))
		}
	}
	q.Log.Info("waiting for the next checkup...")
}

// +build integration

package test

import (
	"context"
	"flag"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/pkg/logger"
	s "github.com/Remneva/otus_hw/hw12_13_14_15_calendar/pkg/storage"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/pkg/storage/sql"
	"github.com/apex/log"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"testing"
)

var dsn = "host=postgres port=5432 user=test password=test dbname=exampledb sslmode=disable"

func TestStorage(t *testing.T) {

	var z zapcore.Level
	flag.Parse()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var ev s.Event
	logg, err := logger.NewLogger(z, "dev", "/dev/null")
	if err != nil {
		log.Fatal("failed to create logger")
	}
	storage := sql.NewStorage(logg)
	err = storage.Connect(ctx, dsn)
	if err != nil {
		log.Fatal("failed to connect db")
	}

	t.Run("CRUD query", func(t *testing.T) {
		id, err := storage.AddEvent(ctx, ev)
		require.Errorf(t, err, "Database query failed")
		require.Equal(t, id, int64(0))

		err = storage.UpdateEvent(ctx, ev)
		require.Errorf(t, err, "Database query failed")

		err = storage.DeleteEvent(ctx, 0)
		require.NoError(t, err)
	})

}

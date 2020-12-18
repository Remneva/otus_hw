package memorystorage

import (
	"context"
	"flag"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/logger"
	_ "github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/storage/sql"
	sqlstorage "github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/apex/log"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"testing"
)

var dsn = "host=localhost port=5432 user=mary password=mary dbname=exampledb sslmode=disable"

func TestStorage(t *testing.T) {

	var z zapcore.Level
	flag.Parse()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	storage := new(Storage)
	var ev sqlstorage.Event
	logg, err := logger.NewLogger(z, "/dev/null")
	if err != nil {
		log.Fatal("failed to create logger")
	}
	err = storage.Connect(ctx, dsn, logg)
	if err != nil {
		log.Fatal("failed to connect db")
	}

	t.Run("CRUD query", func(t *testing.T) {
		id, err := storage.SetEvent(ctx, ev)
		require.Errorf(t, err, "Database query failed")
		require.Equal(t, id, int64(0))

		err = storage.UpdateEvent(ctx, ev)
		require.Errorf(t, err, "Database query failed")

		err = storage.DeleteEvent(ctx, 0)
		require.NoError(t, err)
	})

}

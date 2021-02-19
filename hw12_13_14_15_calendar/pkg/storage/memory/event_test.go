package memorystorage

import (
	"/hw12_13_14_15_calendar/pkg/logger"
	"/hw12_13_14_15_calendar/pkg/storage"
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestMap(t *testing.T) {

	var ev1, ev2, ev3 storage.Event
	ev1.Owner = 1
	ev1.Title = "Title 1"
	ev2.Title = "Title 2"
	ev3.Title = "Title 3"

	t.Run("Events add to map", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		var z zapcore.Level
		logg, _ := logger.NewLogger(z, "dev", "/dev/null")
		e := NewMap(logg)
		_, err := e.AddEvent(ctx, ev1)
		_, err = e.AddEvent(ctx, ev2)
		_, err = e.AddEvent(ctx, ev3)
		actual1, err1 := e.GetEvent(ctx, 1)
		actual2, err2 := e.GetEvent(ctx, 2)
		fmt.Println("err1: ", err1)
		fmt.Println("err2: ", err2)

		require.NoError(t, err1)
		require.NoError(t, err2)
		require.Contains(t, actual1.Title, "Title")
		require.Contains(t, actual2.Title, "Title")
		fmt.Printf("%+v\n", actual1)
		fmt.Printf("%+v\n", actual2)
		require.Equal(t, actual1.ID, int64(1))
		require.Equal(t, actual2.ID, int64(2))

		err = e.DeleteEvent(ctx, 1)
		if err != nil {
			require.NoError(t, err)
		}
	})

	t.Run("No such event in map", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		var z zapcore.Level
		logg, _ := logger.NewLogger(z, "dev", "/dev/null")
		e := NewMap(logg)

		_, error := e.GetEvent(ctx, 10)
		require.EqualError(t, ErrNoSuchEvent, error.Error())

	})
}

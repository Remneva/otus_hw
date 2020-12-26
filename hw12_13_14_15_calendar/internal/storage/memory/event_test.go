package memorystorage

import (
	"context"
	"fmt"
	"testing"

	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestMap(t *testing.T) {

	var ev1, ev2 storage.Event
	ev1.Owner = 1
	ev1.Title = "Title 1"
	ev2.Title = "Title 2"

	t.Run("Events add to map", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		e := NewMap()
		_, err := e.AddEvent(ctx, ev1)
		_, err = e.AddEvent(ctx, ev2)
		actual1, err1 := e.GetEvent(ctx, 1)
		actual2, err2 := e.GetEvent(ctx, 2)
		fmt.Println("err1: ", err1)
		fmt.Println("err2: ", err2)

		require.NoError(t, err1)
		require.NoError(t, err2)
		require.Contains(t, actual1.Title, "Title")
		require.Contains(t, actual2.Title, "Title")

		err = e.DeleteEvent(ctx, 1)
		if err != nil {
			require.NoError(t, err)
		}
	})

	t.Run("No such event in map", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		e := NewMap()

		_, error := e.GetEvent(ctx, 10)
		require.EqualError(t, ErrNoSuchEvent, error.Error())

	})
}

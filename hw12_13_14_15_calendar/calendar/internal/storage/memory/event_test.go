package memorystorage

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestMap(t *testing.T) {

	var ev1, ev2 storage.Event
	ev1.Owner = 1
	ev1.Title = "Title 1"
	ev2.Title = "Title 2"
	var events = []storage.Event{ev1, ev2}

	t.Run("Events add to map", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		wg := sync.WaitGroup{}
		e := NewMap()
		wg.Add(2)

		for _, ev := range events {
			go func(ev storage.Event) {
				e.AddEvent(ctx, ev)
				fmt.Println("done")
				wg.Done()
			}(ev)
		}

		wg.Wait()

		actual1, err1 := e.GetEvent(ctx, 1)
		actual2, err2 := e.GetEvent(ctx, 2)
		fmt.Println("err1: ", err1)
		fmt.Println("err2: ", err2)

		require.NoError(t, err1)
		require.NoError(t, err2)
		require.Contains(t, actual1.Title, "Title")
		require.Contains(t, actual2.Title, "Title")

		err := e.DeleteEvent(ctx, int64(int(ev1.ID)))
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

package storage

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

func TestMap(t *testing.T) {

	var ev1, ev2 Event
	ev1.Owner = 1
	ev1.Title = "Title 1"
	ev2.Title = "Title 2"
	var events = []Event{ev1, ev2}

	t.Run("Events add to map", func(t *testing.T) {
		wg := sync.WaitGroup{}
		e := NewMap()
		wg.Add(2)

		for _, ev := range events {
			go func(ev Event) {
				e.Add(ev)
				fmt.Println("done")
				wg.Done()
			}(ev)
		}

		wg.Wait()

		actual1, err1 := e.Get(1)
		actual2, err2 := e.Get(2)
		fmt.Println("err1: ", err1)
		fmt.Println("err2: ", err2)

		require.NoError(t, err1)
		require.NoError(t, err2)
		require.Contains(t, actual1.Title, "Title")
		require.Contains(t, actual2.Title, "Title")

		err := e.Delete(int(ev1.ID))
		if err != nil {
			require.NoError(t, err)
		}
	})

	t.Run("No such event in map", func(t *testing.T) {

		e := NewMap()

		_, error := e.Get(10)
		require.EqualError(t, ErrNoSuchEvent, error.Error())

	})
}

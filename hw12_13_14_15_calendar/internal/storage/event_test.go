package storage

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMap(t *testing.T) {

	var ev1, ev2 Event
	ev1.Owner = 1
	ev1.Title = "Title 1"
	ev2.Title = "Title 2"

	t.Run("Events add to map", func(t *testing.T) {

		e := NewMap()

		e.Add(ev1)
		e.Add(ev2)
		actual1, bool1, _ := e.Get(1)
		actual2, bool2, _ := e.Get(2)

		require.Equal(t, true, bool1)
		require.Equal(t, true, bool2)
		require.Equal(t, "Title 1", actual1.Title)
		require.Equal(t, "Title 2", actual2.Title)
	})

	t.Run("No such event in map", func(t *testing.T) {

		e := NewMap()

		_, actual, error := e.Get(10)
		require.Equal(t, false, actual)
		require.EqualError(t, ErrNoSuchEvent, error.Error())

	})
}

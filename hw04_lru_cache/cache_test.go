package hw04_lru_cache //nolint:golint,stylecheck

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("Capacity limits logic", func(t *testing.T) {
		c := NewCache(3)
		c.Set("1", 1)
		c.Set("2", 2)
		c.Set("3", 3)
		c.Set("4", 4)

		_, actual := c.Get("1")
		require.False(t, actual)
		_, actual = c.Get("4")
		require.True(t, actual)
	})

	t.Run("Seldom used items first out", func(t *testing.T) {
		c := NewCache(3)
		c.Set("aaa", 1)
		c.Set("bbb", 2)
		c.Set("ccc", 3)

		c.Get("aaa")
		c.Get("ccc")
		c.Get("bbb")
		c.Get("ccc")
		c.Get("bbb")

		c.Set("ddd", 4)

		_, actual := c.Get("aaa")
		require.False(t, actual)

	})

	t.Run("Purge queue logic", func(t *testing.T) {
		c := NewCache(3)
		c.Set("aaa", 1)
		c.Set("bbb", 2)
		c.Set("ccc", 3)

		c.Clear()

		_, actual := c.Get("aaa")
		require.False(t, actual)
		_, actual = c.Get("bbb")
		require.False(t, actual)
		_, actual = c.Get("ccc")
		require.False(t, actual)
	})
}

func TestCacheMultithreading(t *testing.T) {
	t.Skip() // Remove if task with asterisk completed

	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}

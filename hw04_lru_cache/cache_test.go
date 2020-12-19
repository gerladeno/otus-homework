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

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(5)
		_ = c.Set("a", 100)
		_ = c.Set("b", 200)
		_ = c.Set("c", 300)
		c.Clear()
		val, ok := c.Get("a")
		require.False(t, ok)
		require.Nil(t, val)
		val, ok = c.Get("b")
		require.False(t, ok)
		require.Nil(t, val)
		val, ok = c.Get("c")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("cache overflow", func(t *testing.T) {

		c := NewCache(3)
		_ = c.Set("a", 100)
		_ = c.Set("b", 200)
		_ = c.Set("c", 300)
		_ = c.Set("d", 400)

		val, ok := c.Get("a")
		require.False(t, ok)
		require.Nil(t, val)

		val, ok = c.Get("b")
		require.True(t, ok)
		require.Equal(t, 200, val)
		val, ok = c.Get("c")
		require.True(t, ok)
		require.Equal(t, 300, val)
		val, ok = c.Get("d")
		require.True(t, ok)
		require.Equal(t, 400, val)

		_ = c.Set("e", 500)
		_ = c.Set("f", 600)
		_ = c.Set("g", 700)
		_ = c.Set("h", 800)

		val, ok = c.Get("b")
		require.False(t, ok)
		require.Nil(t, val)
		val, ok = c.Get("c")
		require.False(t, ok)
		require.Nil(t, val)
		val, ok = c.Get("d")
		require.False(t, ok)
		require.Nil(t, val)
		val, ok = c.Get("e")
		require.False(t, ok)
		require.Nil(t, val)

		val, ok = c.Get("f")
		require.True(t, ok)
		require.Equal(t, 600, val)
		val, ok = c.Get("g")
		require.True(t, ok)
		require.Equal(t, 700, val)
		val, ok = c.Get("h")
		require.True(t, ok)
		require.Equal(t, 800, val)
	})
}

func TestCacheMultithreading(t *testing.T) {
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

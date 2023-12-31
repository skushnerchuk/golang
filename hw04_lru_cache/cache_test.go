package hw04lrucache

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
		c := NewCache(2)
		_ = c.Set("1", 1)
		_ = c.Set("2", 2)
		c.Clear()

		_, ok := c.Get("1")
		require.False(t, ok)

		_, ok = c.Get("2")
		require.False(t, ok)
	})

	t.Run("push logic", func(t *testing.T) {
		c := NewCache(3)
		_ = c.Set("1", 1)
		_ = c.Set("2", 2)
		_ = c.Set("3", 3)
		_ = c.Set("4", 4)

		_, ok := c.Get("1")
		require.False(t, ok)
	})

	t.Run("remove oldest element", func(t *testing.T) {
		c := NewCache(3)
		_ = c.Set("1", 1)
		_ = c.Set("2", 2)
		_ = c.Set("3", 3)

		_, _ = c.Get("1")
		_, _ = c.Get("2")
		_, _ = c.Get("3")
		_ = c.Set("3", 3)
		_ = c.Set("2", 2)
		_ = c.Set("1", 1)
		_ = c.Set("4", 4)

		_, ok := c.Get("3")
		require.False(t, ok)
	})

	t.Run("check capacity", func(t *testing.T) {
		c := NewCache(3)
		_ = c.Set("1", 1)
		_ = c.Set("2", 2)
		_ = c.Set("3", 3)

		_, ok1 := c.Get("1")
		_, ok2 := c.Get("2")
		_, ok3 := c.Get("3")
		_, ok4 := c.Get("4")

		require.True(t, ok1 && ok2 && ok3)
		require.False(t, ok4)

		_ = c.Set("4", 4)

		_, ok1 = c.Get("1")
		_, ok2 = c.Get("2")
		_, ok3 = c.Get("3")
		_, ok4 = c.Get("4")

		require.False(t, ok1)
		require.True(t, ok2 && ok3 && ok4)
	})
}

func TestCacheMultithreading(_ *testing.T) {
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

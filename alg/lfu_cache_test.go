package alg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLFUCachePut(t *testing.T) {
	cache := NewLFUCache[int, int](4)
	for i := 0; i < 100; i++ {
		cache.Put(i, i)
	}

	a, err := cache.Get(99)
	assert.Equal(t, nil, err)
	assert.Equal(t, 99, a)
	a, err = cache.Get(10)
	assert.Equal(t, ErrLFUCacheKeyNotFound, err)
	assert.NotEqual(t, 10, a)

	cache.Put(99, 99)
	cache.Put(99, 99)
	cache.Put(99, 99)
	cache.Put(99, 99)
	cache.Put(99, 99)
	cache.Put(98, 98)
	cache.Put(98, 98)
	cache.Put(98, 98)
	cache.Put(98, 98)
	cache.Put(97, 97)
	cache.Put(97, 97)
	cache.Put(97, 97)
	cache.Put(96, 96)
	cache.Put(96, 96)

	assert.Equal(t, []int{99, 98, 97, 96}, cache.GetTopHotKeys(4))
}

func TestLFUCacheDelete(t *testing.T) {
	cache := NewLFUCache[int, int](4)
	for i := 0; i < 100; i++ {
		cache.Put(i, i)
	}

	a, err := cache.Get(99)
	assert.Equal(t, nil, err)
	assert.Equal(t, 99, a)
	a, err = cache.Get(10)
	assert.Equal(t, ErrLFUCacheKeyNotFound, err)
	assert.NotEqual(t, 10, a)

	cache.Put(99, 99)
	cache.Put(99, 99)
	cache.Put(99, 99)
	cache.Put(99, 99)
	cache.Put(99, 99)
	cache.Put(98, 98)
	cache.Put(98, 98)
	cache.Put(98, 98)
	cache.Put(98, 98)
	cache.Put(97, 97)
	cache.Put(97, 97)
	cache.Put(97, 97)
	cache.Put(96, 96)
	cache.Put(96, 96)

	assert.Equal(t, []int{99, 98, 97, 96}, cache.GetTopHotKeys(4))

	cache.Delete(97)
	cache.Put(101, 101)

	assert.Equal(t, []int{99, 98, 96, 101}, cache.GetTopHotKeys(4))
}

func BenchmarkLFUCachePut(b *testing.B) {
	cache := NewLFUCache[int, int](400000)
	for i := 0; i < b.N; i++ {
		cache.Put(i, i)
	}
}

func BenchmarkLFUCacheGet(b *testing.B) {
	cache := NewLFUCache[int, int](400000)
	for i := 0; i < 400000; i++ {
		cache.Put(i, i)
	}

	for i := 0; i < b.N; i++ {
		cache.Get(399999)
	}
}

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

	assert.Equal(t, []int{99, 98, 97, 96}, cache.GetTopHotKeys(4))

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

func TestLFUCacheTopHotKeys(t *testing.T) {
	cache := NewLFUCache[int, int](4)

	for i := 0; i < 400000; i++ {
		cache.Put(i, i)
	}
	assert.Equal(t, []int{399999, 399998, 399997, 399996}, cache.GetTopHotKeys(4))

	for i := 400000; i >= 400000-4; i-- {
		// input 400000 | 1 => 400000, 399999, 399998, 399997
		// input 399999 | 1 => 400000, 399998, 399997 | 2 => 399999
		// input 399998 | 1 => 400000, 399997 | 2 => 399998,399999
		// input 399997 | 1 => 400000 | 2 => 399997, 399998, 399999
		// input 399996 | 1 => 399996 | 2 => 399997, 399998, 399999
		// linked list order -> 399997, 399998, 399999, 399996
		// println("----input ", i)
		cache.Put(i, i)
	}
	// cache.Debug()

	assert.Equal(t, []int{399997, 399998, 399999, 399996}, cache.GetTopHotKeys(4))
}

func BenchmarkLFUCachePut(b *testing.B) {
	cache := NewLFUCache[int, int](400000)
	for i := 0; i < b.N; i++ {
		cache.Put(i, i)
	}
}

func BenchmarkLFUCacheGet(b *testing.B) {
	cache := NewLFUCache[int, int](b.N)
	for i := 0; i <= b.N; i++ {
		cache.Put(i, i)
	}

	b.ResetTimer()

	for i := 0; i <= b.N; i++ {
		cache.Get(i)
	}
}

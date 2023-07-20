package alg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLFUCache(t *testing.T) {
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

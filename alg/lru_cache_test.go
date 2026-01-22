package alg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLRUCacheBasic(t *testing.T) {
	cache := NewLRUCache[string, int](3)

	// 测试 Put 和 Get
	cache.Put("a", 1)
	cache.Put("b", 2)
	cache.Put("c", 3)

	assert.Equal(t, 3, cache.Len())

	val, err := cache.Get("a")
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, val)

	val, err = cache.Get("b")
	assert.Equal(t, nil, err)
	assert.Equal(t, 2, val)
}

func TestLRUCacheEviction(t *testing.T) {
	cache := NewLRUCache[string, int](3)

	cache.Put("a", 1)
	cache.Put("b", 2)
	cache.Put("c", 3)

	// 容量已满，插入新节点应该淘汰最久未使用的 "a"
	cache.Put("d", 4)

	assert.Equal(t, 3, cache.Len())

	// "a" 应该被淘汰
	_, err := cache.Get("a")
	assert.NotEqual(t, nil, err)

	// 其他节点应该存在
	val, err := cache.Get("b")
	assert.Equal(t, nil, err)
	assert.Equal(t, 2, val)

	val, err = cache.Get("c")
	assert.Equal(t, nil, err)
	assert.Equal(t, 3, val)

	val, err = cache.Get("d")
	assert.Equal(t, nil, err)
	assert.Equal(t, 4, val)
}

func TestLRUCacheAccessOrder(t *testing.T) {
	cache := NewLRUCache[string, int](3)

	cache.Put("a", 1)
	cache.Put("b", 2)
	cache.Put("c", 3)

	// 访问 "a"，使其成为最近使用
	cache.Get("a")

	// 插入新节点应该淘汰最久未使用的 "b"
	cache.Put("d", 4)

	assert.Equal(t, 3, cache.Len())

	// "a" 应该还在（最近访问过）
	val, err := cache.Get("a")
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, val)

	// "b" 应该被淘汰
	_, err = cache.Get("b")
	assert.NotEqual(t, nil, err)
}

func TestLRUCacheUpdate(t *testing.T) {
	cache := NewLRUCache[string, int](2)

	cache.Put("a", 1)
	cache.Put("b", 2)

	// 更新 "a" 的值
	cache.Put("a", 10)

	val, err := cache.Get("a")
	assert.Equal(t, nil, err)
	assert.Equal(t, 10, val)

	assert.Equal(t, 2, cache.Len())
}

func TestLRUCacheDelete(t *testing.T) {
	cache := NewLRUCache[string, int](3)

	cache.Put("a", 1)
	cache.Put("b", 2)
	cache.Put("c", 3)

	// 删除 "b"
	cache.Delete("b")

	assert.Equal(t, 2, cache.Len())

	// "b" 应该不存在
	_, err := cache.Get("b")
	assert.NotEqual(t, nil, err)

	// 其他节点应该存在
	val, err := cache.Get("a")
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, val)

	val, err = cache.Get("c")
	assert.Equal(t, nil, err)
	assert.Equal(t, 3, val)
}

func TestLRUCacheDeleteNonExistent(t *testing.T) {
	cache := NewLRUCache[string, int](3)

	cache.Put("a", 1)

	// 删除不存在的 key 不应该 panic
	cache.Delete("b")

	assert.Equal(t, 1, cache.Len())
}

func TestLRUCacheGetSize(t *testing.T) {
	cache := NewLRUCache[string, int](5)
	assert.Equal(t, 5, cache.GetSize())

	cache.Put("a", 1)
	assert.Equal(t, 5, cache.GetSize())
}

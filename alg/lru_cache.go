package alg

import (
	"errors"
)

func NewLRUCache[KT comparable, VT any](size int) (t *LRUCache[KT, VT]) {
	t = new(LRUCache[KT, VT])
	t.orderList = new(LinkedList[*LRUCacheNode[KT, VT]])
	t.data = make(map[KT]*LinkedListNode[*LRUCacheNode[KT, VT]])
	t.size = size
	return t
}

type LRUCache[KT comparable, VT any] struct {
	// orderList 访问顺序链表，头部是最近访问，尾部是最久未使用
	orderList *LinkedList[*LRUCacheNode[KT, VT]]
	// data 实际数据存储，key 到链表节点的映射
	data map[KT]*LinkedListNode[*LRUCacheNode[KT, VT]]
	// size 缓存的最大长度
	size int
}

// LRUCacheNode LRU 缓存节点
type LRUCacheNode[KT comparable, VT any] struct {
	// Key 缓存键
	Key KT
	// Value 缓存值
	Value VT
}

// Len return length of current cache item count
// Len 返回当前长度大小
func (t *LRUCache[KT, VT]) Len() int {
	return len(t.data)
}

// Put update or insert new key and value
// Put 更新或插入新的键值对
func (t *LRUCache[KT, VT]) Put(key KT, value VT) {
	if node, exists := t.data[key]; exists {
		// key 已存在，更新值并移到头部
		node.Data.Value = value
		t.orderList.MovePrePend(node, t.orderList.Head)
		return
	}

	// 容量已满，移除最久未使用的节点（尾部）
	if t.Len() >= t.size && t.orderList.Tail != nil {
		tailKey := t.orderList.Tail.Data.Key
		t.orderList.Remove(t.orderList.Tail)
		delete(t.data, tailKey)
	}

	// 创建新节点并插入头部
	cacheNode := &LRUCacheNode[KT, VT]{Key: key, Value: value}
	listNode := NewLinkedListNode(cacheNode)
	if t.orderList.Head == nil {
		t.orderList.Append(listNode)
	} else {
		t.orderList.PrePend(t.orderList.Head, listNode)
	}
	t.data[key] = listNode
}

// GetSize return max length of storage slots can use
// GetSize 返回最大可用的存储槽个数
func (t *LRUCache[KT, VT]) GetSize() int {
	return t.size
}

// Get get value by key name
// Get 通过key取得value，返回值可能为零值或空值
func (t *LRUCache[KT, VT]) Get(key KT) (value VT, err error) {
	node, exists := t.data[key]
	if !exists {
		return value, errors.New("key not found")
	}
	// 将访问的节点移到头部
	t.orderList.MovePrePend(node, t.orderList.Head)
	return node.Data.Value, nil
}

// Delete delete key and update cache storage, do not forget call ptr variable's Close() method
// Delete 删除key并更新缓存存储， 别忘了调用指针变量的 Close() 方法
func (t *LRUCache[KT, VT]) Delete(key KT) {
	node, exists := t.data[key]
	if !exists {
		return
	}
	t.orderList.Remove(node)
	delete(t.data, key)
}

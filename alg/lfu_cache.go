package alg

import (
	"fmt"
	"sync"
)

var (
	ErrLFUCacheKeyNotFound = fmt.Errorf("key not found")
)

func NewLFUCache[KT string | int, VT any](size int) (t *LFUCache[KT, VT]) {
	t = new(LFUCache[KT, VT])
	t.size = size
	t.frequency = make([]KT, size)
	t.data = make(map[KT]*LFUCacheNode[VT])
	return t
}

type LFUCacheNode[VT any] struct {
	Value     VT
	Count     uint64
	SlotIndex int
}

// LFUCache thread safe latest frequency use cache
// LFUCache 线程安全的LFU缓存
type LFUCache[KT string | int, VT any] struct {
	// frequency 链表，当插入新的值时，访问尾端节点的值，即得到应该应该操作那个key
	frequency []KT
	// data 实际数据存储
	data map[KT]*LFUCacheNode[VT]
	// size 缓存的最大长度
	size int
	// currentSize 当前存储的容量
	currentSize int
	lock        sync.RWMutex
}

func (t *LFUCache[KT, VT]) GetSize() int {
	return t.size
}

func (t *LFUCache[KT, VT]) Put(key KT, value VT) {
	t.lock.Lock()
	defer t.lock.Unlock()

	if _, found := t.data[key]; found {
		t.frequencyInc(key)
		return
	}

	slotIndex := t.currentSize
	if t.currentSize >= t.size { // 满了就要找最后一位进行剔除,这会导致老key永远清不掉，而新key进来后下一次插入其他key则刚插入的key就立马被删除了
		delete(t.data, t.frequency[t.currentSize-1])
	} else {
		t.currentSize++
	}

	if slotIndex >= t.size {
		slotIndex = t.size - 1
	}

	t.frequency[slotIndex] = key

	t.data[key] = &LFUCacheNode[VT]{
		Value:     value,
		Count:     1,
		SlotIndex: slotIndex,
	}
}

// Get 通过key取得value，返回值可能为零值或空值
func (t *LFUCache[KT, VT]) Get(key KT) (value VT, err error) {
	t.lock.Lock()
	defer t.lock.Unlock()
	if node, found := t.data[key]; found {
		t.frequencyInc(key)
		value = node.Value
		return value, nil
	} else {
		return value, ErrLFUCacheKeyNotFound
	}
}

// frequencyInc 增加频率
func (t *LFUCache[KT, VT]) frequencyInc(key KT) {
	node := t.data[key]
	if node != nil {
		node.Count++
		for range t.frequency { // 无限向前排序
			prevSlotIndex := node.SlotIndex - 1
			nodeSlotIndex := node.SlotIndex
			if prevSlotIndex > -1 {
				prevNode := t.data[t.frequency[prevSlotIndex]]
				if prevNode != nil && prevNode.Count < node.Count {
					prevNodeKey := t.frequency[prevSlotIndex]
					t.frequency[prevSlotIndex] = key
					t.frequency[nodeSlotIndex] = prevNodeKey
					node.SlotIndex = prevSlotIndex
					prevNode.SlotIndex = nodeSlotIndex
				} else {
					break
				}
			} else {
				break
			}
		}
	}
}

// GetTopHotKeys 取得最热的几个key，热度按数组下标从左（小）到右（大）排列,如果要取得全部key的排列信息，则 topNum = GetSize() 即可
func (t *LFUCache[KT, VT]) GetTopHotKeys(topNum int) (keys []KT) {
	t.lock.RLock()
	defer t.lock.RUnlock()
	for i, key := range t.frequency {
		if i < topNum {
			keys = append(keys, key)
		}
	}
	return
}

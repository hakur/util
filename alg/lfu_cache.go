package alg

import (
	"fmt"
	"sort"
	"time"
)

var (
	ErrLFUCacheKeyNotFound = fmt.Errorf("lfu cache key not found")
)

// NewLFUCache create new lfu cache with fixed size
// NewLFUCache 新建固定大小的LFU缓存
func NewLFUCache[KT comparable, VT any](size int) (t *LFUCache[KT, VT]) {
	t = new(LFUCache[KT, VT])
	t.size = size
	t.frequency = NewLinkedList[KT]()

	t.data = make(map[KT]*LFUCacheNode[KT, VT])
	return t
}

// LFUCacheNode lfu cache node
// LFUCacheNode lfu缓存节点
type LFUCacheNode[KT comparable, VT any] struct {
	// Value 缓存存储的实际内容
	Value VT
	// Count 访问次数统计，溢出后将成为新的值
	Count uint64
	// LinkedListNode 指向频率链表节点，用于交换位置
	LinkedListNode *LinkedListNode[KT]
	// AccessTime 访问时间，使用简单的int64而不是复杂的time.Time
	AccessTime int64
	SlotIndex  int
}

// LFUCache thread safe latest frequency use cache, when add lock to NewLFUCache instance, must use sync.Mutex
// LFUCache 非线程安全的LFU缓存,对 NewLFUCache 产生的实例加锁时，务必使用互斥锁。
type LFUCache[KT comparable, VT any] struct {
	// frequency 链表，当插入新的值时，访问尾端节点的值，即得到应该应该操作那个key
	frequency *LinkedList[KT]
	// slots 排序用的数组
	slots []*LFUCacheNode[KT, VT]
	// data 实际数据存储
	data map[KT]*LFUCacheNode[KT, VT]
	// size 缓存的最大长度
	size int
	// currentSize 当前存储的容量
	currentSize int
	// lock        sync.RWMutex // 并不是每个人都想要锁住的，看个人的编码情况自行决定是否加锁
}

// Len return length of current cache item count
// Len 返回当前长度大小
func (t *LFUCache[KT, VT]) Len() int {
	return t.currentSize
}

// GetSize return max length of storage slots can use
// GetSize 返回最大可用的存储槽个数
func (t *LFUCache[KT, VT]) GetSize() int {
	return t.size
}

// Put update or insert new key and value
// Put 更新或插入新的键值对
func (t *LFUCache[KT, VT]) Put(key KT, value VT) {
	// t.lock.Lock()
	// defer t.lock.Unlock()

	if _, found := t.data[key]; found {
		t.frequencyInc(key)
		return
	}

	if t.currentSize >= t.size { // 满了就要找最后一位进行剔除,这会导致老key永远清不掉，而新key进来后下一次插入其他key则刚插入的key就立马被删除了
		delete(t.data, t.frequency.Tail.Data)
		t.frequency.Tail.Data = key
	} else {
		t.currentSize++
		t.frequency.Append(NewLinkedListNode(key))
	}

	node := &LFUCacheNode[KT, VT]{
		Value:          value,
		Count:          1,
		LinkedListNode: t.frequency.Tail,
		AccessTime:     time.Now().UnixNano(),
		SlotIndex:      t.currentSize - 1,
	}
	t.data[key] = node

	t.slots = append(t.slots, node)
}

// Get get value by key name
// Get 通过key取得value，返回值可能为零值或空值
func (t *LFUCache[KT, VT]) Get(key KT) (value VT, err error) {
	// t.lock.Lock()
	// defer t.lock.Unlock()
	if node, found := t.data[key]; found {
		t.frequencyInc(key)
		value = node.Value
		return value, nil
	} else {
		return value, ErrLFUCacheKeyNotFound
	}
}

// frequencyInc increase key use count and update access time
// frequencyInc 增加key的频率和更新访问时间
func (t *LFUCache[KT, VT]) frequencyInc(key KT) {
	node := t.data[key]
	if node != nil {
		node.Count++
		node.AccessTime = time.Now().UnixNano()
		if node.SlotIndex < 1 {
			return
		}

		prevNodeIndex := sort.Search(node.SlotIndex, func(i int) bool {
			return t.slots[i].SlotIndex < node.SlotIndex || t.slots[i].AccessTime < node.AccessTime
		})

		if prevNodeIndex < node.SlotIndex {
			prevNode := t.slots[prevNodeIndex]
			t.slots[node.SlotIndex], t.slots[prevNodeIndex] = t.slots[prevNodeIndex], t.slots[node.SlotIndex] // 交换排序指针
			node.SlotIndex, prevNode.SlotIndex = prevNode.SlotIndex, node.SlotIndex                           // 交换排序的槽位
			node.LinkedListNode, prevNode.LinkedListNode = prevNode.LinkedListNode, node.LinkedListNode       // 交换链表的值
			t.frequency.SwapData(prevNode.LinkedListNode, node.LinkedListNode)                                // 交换链表的位置
		}
	}
}

// GetTopHotKeys get most biggest access count of key
// GetTopHotKeys 取得最热的几个key，热度按数组下标从左（大）到右（小）排列,如果要取得全部key的排列信息，则 topNum = GetSize() 即可
func (t *LFUCache[KT, VT]) GetTopHotKeys(topNum int) (keys []KT) {
	// t.lock.RLock()
	// defer t.lock.RUnlock()
	var count int
	node := t.frequency.Head
	for {
		if node == nil {
			break
		}
		if count < topNum {
			keys = append(keys, node.Data)
		}
		count++
		node = node.Next
	}
	return
}

// Delete delete key and update cache storage, do not forget call ptr variable's Close() method
// Delete 删除key并更新缓存存储， 别忘了调用指针变量的 Close() 方法
func (t *LFUCache[KT, VT]) Delete(key KT) {
	var found bool
	var node *LFUCacheNode[KT, VT]
	if node, found = t.data[key]; !found {
		return
	}
	t.frequency.Remove(node.LinkedListNode)
	delete(t.data, key)
	t.slots[node.SlotIndex] = nil
	t.slots = append(t.slots[:node.SlotIndex], t.slots[node.SlotIndex+1:]...)

	t.currentSize--
}

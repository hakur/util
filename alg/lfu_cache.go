package alg

import (
	"fmt"
	"sort"
	"time"
)

var (
	ErrLFUCacheKeyNotFound = fmt.Errorf("lfu cache key not found")
)

// NewLFUCache create new lfu cache with pre allocated size, but with some garbage collection, there is extra memory "in use"
// NewLFUCache 新建固定大小的lfu缓存, 因为GC的问题导致看起来有一些额外的内存“占用”
func NewLFUCache[KT comparable, VT any](size int) (t *LFUCache[KT, VT]) {
	t = new(LFUCache[KT, VT])
	t.size = size
	t.frequency = make([]KT, 0, size) // 保证slice的底层数组是同一个

	t.data = make(map[KT]*LFUCacheNode[VT])
	return t
}

// LFUCacheNode lfu cache node, use array instead of linked list, fast get key's value, but takes big cost when delete key
// LFUCacheNode lfu缓存节点, 使用数组替代链表，快速取得key的value，但删除key的接口消耗巨大
type LFUCacheNode[VT any] struct {
	// Value 缓存存储的实际内容
	Value VT
	// Count 访问次数统计，溢出后将成为新的值
	Count uint64
	// SlotIndex 存在于 LFUCache.data 的那个数组下标
	SlotIndex int
	// AccessTime 访问时间，使用简单的int64而不是复杂的time.Time
	AccessTime int64
}

// LFUCache thread safe latest frequency use cache, when add lock to NewLFUCache instance, must use sync.Mutex
// LFUCache 非线程安全的LFU缓存,对 NewLFUCache 产生的实例加锁时，务必使用互斥锁。
type LFUCache[KT comparable, VT any] struct {
	// frequency 链表，当插入新的值时，访问尾端节点的值，即得到应该应该操作那个key
	frequency []KT
	// data 实际数据存储
	data map[KT]*LFUCacheNode[VT]
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

	slotIndex := t.currentSize
	if t.currentSize >= t.size { // 满了就要找最后一位进行剔除,这会导致老key永远清不掉，而新key进来后下一次插入其他key则刚插入的key就立马被删除了
		delete(t.data, t.frequency[t.currentSize-1])
		t.frequency = append(t.frequency[:t.currentSize-1], t.frequency[t.currentSize:]...)
	} else {
		t.currentSize++
	}

	// 填充逻辑，此时data存储槽还没有满
	if slotIndex >= t.size {
		slotIndex = t.size - 1
	}

	// t.frequency[slotIndex] = key
	t.frequency = append(t.frequency, key)

	t.data[key] = &LFUCacheNode[VT]{
		Value:      value,
		Count:      1,
		SlotIndex:  slotIndex,
		AccessTime: time.Now().UnixNano(),
	}
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

		if t.currentSize > 2048 {
			// 数据量更大的时候，sort.Search内置的算法会更快一些
			nodeSlotIndex := node.SlotIndex
			if nodeSlotIndex > 0 {
				prevSlotIndex := sort.Search(nodeSlotIndex, func(i int) bool {
					prevNode, ok := t.data[t.frequency[i]]
					if !ok { // key的值被删除了
						return true
					}
					if prevNode.Count < node.Count {
						return true
					} else if prevNode.Count == node.Count {
						if prevNode.AccessTime < node.AccessTime {
							return true
						}
					}
					return false
				})

				// 时间小于当前node，则当前node可以被认为是更活跃的，将当前node移动到data数组的左侧一位
				if prevSlotIndex >= t.currentSize {
					prevSlotIndex = t.currentSize - 1
				} else if prevSlotIndex < 0 {
					prevSlotIndex = 0
				}

				if prevSlotIndex != nodeSlotIndex {
					prevNode, ok := t.data[t.frequency[prevSlotIndex]]
					if ok {
						prevNodeKey := t.frequency[prevSlotIndex]
						t.frequency[prevSlotIndex] = key
						t.frequency[nodeSlotIndex] = prevNodeKey
						node.SlotIndex = prevSlotIndex
						prevNode.SlotIndex = nodeSlotIndex
					} else {
						t.frequency[prevSlotIndex] = key
						node.SlotIndex = prevSlotIndex
					}
				}
			}
		} else {
			// 数据量更低的时候，线性会更快一些
			for range t.frequency { // 无限向前搜索并排序
				prevSlotIndex := node.SlotIndex - 1
				nodeSlotIndex := node.SlotIndex
				if prevSlotIndex > -1 {
					prevNode, found := t.data[t.frequency[prevSlotIndex]]
					// 时间小于当前node，则当前node可以被认为是更活跃的，将当前node移动到data数组的左侧一位
					var needSwitch bool
					if !found { // key的值被删除了
						needSwitch = true
					} else if prevNode.Count < node.Count {
						needSwitch = true
					} else if prevNode.Count == node.Count {
						if prevNode.AccessTime < node.AccessTime {
							needSwitch = true
						}
					}
					if prevSlotIndex != nodeSlotIndex && needSwitch {
						if prevNode != nil {
							prevNodeKey := t.frequency[prevSlotIndex]
							t.frequency[prevSlotIndex] = key
							t.frequency[nodeSlotIndex] = prevNodeKey
							node.SlotIndex = prevSlotIndex
							prevNode.SlotIndex = nodeSlotIndex
						} else {
							t.frequency[prevSlotIndex] = key
							node.SlotIndex = prevSlotIndex
						}
					} else {
						break
					}
				} else {
					break
				}
			}
		}
	}
}

// GetTopHotKeys get most biggest access count of key
// GetTopHotKeys 取得最热的几个key，热度按数组下标从左（小）到右（大）排列,如果要取得全部key的排列信息，则 topNum = GetSize() 即可
func (t *LFUCache[KT, VT]) GetTopHotKeys(topNum int) (keys []KT) {
	// t.lock.RLock()
	// defer t.lock.RUnlock()
	var count int
	for _, key := range t.frequency {
		if _, found := t.data[key]; found {
			if count < topNum {
				keys = append(keys, key)
				count++
			}
		}
	}
	return
}

// Delete delete key and update cache storage, this methos takes big cost
// Delete 删除key并更新缓存存储， 这个接口消耗是巨大的
func (t *LFUCache[KT, VT]) Delete(key KT) {
	var found bool
	var node *LFUCacheNode[VT]
	if node, found = t.data[key]; !found {
		return
	}
	delete(t.data, key)
	t.currentSize--

	// very big cost
	t.frequency = append(t.frequency[:node.SlotIndex], t.frequency[node.SlotIndex+1:]...)
	for _, n := range t.data {
		n.SlotIndex++
	}
}

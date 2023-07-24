package alg

import (
	"fmt"
	"sort"
	"time"
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
type LFUCache[KT string | int, VT any] struct {
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

func (t *LFUCache[KT, VT]) GetSize() int {
	return t.size
}

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
	} else {
		t.currentSize++
	}

	// 填充逻辑，此时data存储槽还没有满
	if slotIndex >= t.size {
		slotIndex = t.size - 1
	}

	t.frequency[slotIndex] = key

	t.data[key] = &LFUCacheNode[VT]{
		Value:      value,
		Count:      1,
		SlotIndex:  slotIndex,
		AccessTime: time.Now().UnixNano(),
	}
}

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

// frequencyInc 增加频率
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
					prevNode := t.data[t.frequency[i]]
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
				}

				if prevSlotIndex != nodeSlotIndex {
					prevNode := t.data[t.frequency[prevSlotIndex]]
					prevNodeKey := t.frequency[prevSlotIndex]
					t.frequency[prevSlotIndex] = key
					t.frequency[nodeSlotIndex] = prevNodeKey
					node.SlotIndex = prevSlotIndex
					prevNode.SlotIndex = nodeSlotIndex
				}
			}
		} else {
			// 数据量更低的时候，线性会更快一些
			for range t.frequency { // 无限向前搜索并排序
				prevSlotIndex := node.SlotIndex - 1
				nodeSlotIndex := node.SlotIndex
				if prevSlotIndex > -1 {
					prevNode := t.data[t.frequency[prevSlotIndex]]
					// 时间小于当前node，则当前node可以被认为是更活跃的，将当前node移动到data数组的左侧一位
					var needSwitch bool
					if prevNode.Count < node.Count {
						needSwitch = true
					} else if prevNode.Count == node.Count {
						if prevNode.AccessTime < node.AccessTime {
							needSwitch = true
						}
					}
					if prevNode != nil && needSwitch {
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
}

// GetTopHotKeys 取得最热的几个key，热度按数组下标从左（小）到右（大）排列,如果要取得全部key的排列信息，则 topNum = GetSize() 即可
func (t *LFUCache[KT, VT]) GetTopHotKeys(topNum int) (keys []KT) {
	// t.lock.RLock()
	// defer t.lock.RUnlock()
	for i, key := range t.frequency {
		if i < topNum {
			keys = append(keys, key)
		}
	}
	return
}

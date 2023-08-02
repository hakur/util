package alg

import (
	"encoding/json"
	"fmt"
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
	t.frequencyGroup = make(map[uint64]*LFUFrequencyGroupEntry[KT])
	t.data = make(map[KT]*LFUCacheNode[KT, VT])
	return t
}

// LFUFrequencyGroupEntry 使用频率组
type LFUFrequencyGroupEntry[KT comparable] struct {
	// // List 引用频率链表，用于修改Head和Tail
	// List *LinkedList[KT]
	// Head 这个频率组内访问时间最新的一条
	Head *LinkedListNode[KT]
	// Tail 这个频率组内访问时间最旧的一条
	Tail *LinkedListNode[KT]
	// Length 这个频率组内一共有多少个链表节点
	Length int
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
}

// LFUCache thread safe latest frequency use cache, when add lock to NewLFUCache instance, must use sync.Mutex
// LFUCache 非线程安全的LFU缓存,对 NewLFUCache 产生的实例加锁时，务必使用互斥锁。
type LFUCache[KT comparable, VT any] struct {
	// frequency 链表，当插入新的值时，访问尾端节点的值，即得到应该应该操作那个key,使用频率的总体链表
	frequency *LinkedList[KT]
	// frequencyGroup 使用频率的链表头，将按照使用频次对同一个链表进行链表数据段分组,指针指向这段链表的第一个元素
	frequencyGroup map[uint64]*LFUFrequencyGroupEntry[KT]
	// data 实际数据存储
	data map[KT]*LFUCacheNode[KT, VT]
	// size 缓存的最大长度
	size int
}

// Len return length of current cache item count
// Len 返回当前长度大小
func (t *LFUCache[KT, VT]) Len() int {
	return len(t.data)
}

// GetSize return max length of storage slots can use
// GetSize 返回最大可用的存储槽个数
func (t *LFUCache[KT, VT]) GetSize() int {
	return t.size
}

// Put update or insert new key and value
// Put 更新或插入新的键值对
func (t *LFUCache[KT, VT]) Put(key KT, value VT) {
	if node, found := t.data[key]; found {
		t.frequencyInc(node)
		return
	}

	if len(t.data) >= t.size { // 满了就要找最后一位进行剔除,这会导致老key永远清不掉，而新key进来后下一次插入其他key则刚插入的key就立马被删除了
		t.Delete(t.frequency.Tail.Data)
		t.frequency.Append(NewLinkedListNode(key))
	} else {
		t.frequency.Append(NewLinkedListNode(key))
	}

	node := &LFUCacheNode[KT, VT]{
		Value:          value,
		Count:          0,
		LinkedListNode: t.frequency.Tail,
	}
	t.data[key] = node
	t.frequencyInc(node)
}

// Get get value by key name
// Get 通过key取得value，返回值可能为零值或空值
func (t *LFUCache[KT, VT]) Get(key KT) (value VT, err error) {
	if node, found := t.data[key]; found {
		t.frequencyInc(node)
		value = node.Value
		return value, nil
	} else {
		return value, ErrLFUCacheKeyNotFound
	}
}

// frequencyInc increase key use count and update access time
// frequencyInc 增加key的频率和更新访问时间
func (t *LFUCache[KT, VT]) frequencyInc(node *LFUCacheNode[KT, VT]) {
	node.Count++
	if _, ok := t.frequencyGroup[node.Count]; !ok {
		// 这个频率组不存在就新建
		t.frequencyGroup[node.Count] = &LFUFrequencyGroupEntry[KT]{
			Head:   node.LinkedListNode,
			Tail:   node.LinkedListNode,
			Length: 1,
		}
		if node.Count > 1 { // 说明这不是初始化节点，是确确实实的节点频率增加了
			t.frequency.MovePrePend(node.LinkedListNode, t.frequencyGroup[node.Count-1].Head) // 操作链表,插入到频率组头部节点的前面，因为当前节点的时间是最新的
			t.removeNodeFromFrequencyGroup(node, node.Count-1)                                // 逃离旧的频率组
		}
	} else {
		// 存在就表示需要追加到新的位置上, 按访问频次和访问时间排序
		t.frequency.MovePrePend(node.LinkedListNode, t.frequencyGroup[node.Count].Head) // 操作链表,插入到频率组头部节点的前面，因为当前节点的时间是最新的
		t.frequencyGroup[node.Count].Head = node.LinkedListNode                         // 修改对应组的头部
		t.frequencyGroup[node.Count].Length++                                           // 对应组的数量增加
		if node.Count > 1 {                                                             // 说明这不是初始化节点，是确确实实的节点频率增加了
			t.removeNodeFromFrequencyGroup(node, node.Count-1) // 逃离旧的频率组
		}
	}
}

func (t *LFUCache[KT, VT]) removeNodeFromFrequencyGroup(node *LFUCacheNode[KT, VT], nodeGroupKey uint64) {
	prevNodeGroupKey := nodeGroupKey // node.Count - 1
	group := t.frequencyGroup[prevNodeGroupKey]
	if group == nil {
		return
	}

	group.Length--

	if group.Length < 1 {
		delete(t.frequencyGroup, prevNodeGroupKey)
		return
	}

	if group.Head == node.LinkedListNode { // 如果移除的是组的头部,那么需要寻找当前节点的下一个节点来替代
		nextNode := node.LinkedListNode.Next
		if nextNode == nil { // 如果没有下一个节点则说明当前这个频率组已经没有后继节点，需要清理掉
			delete(t.frequencyGroup, prevNodeGroupKey)
			return
		}
		// 如果还有后续节点，则检查一下引用频次是否等于当前组的频率
		if t.data[nextNode.Data] != nil && (t.data[nextNode.Data].Count == prevNodeGroupKey) { // 如果相同，则可以后继
			group.Head = nextNode
		} else { // 否则不能后继，需要删除这个频率组
			delete(t.frequencyGroup, prevNodeGroupKey)
			return
		}
	} else if group.Tail == node.LinkedListNode { // 如果移除的是组的尾部，那么需要找到上一个节点来作为后继
		prevNode := node.LinkedListNode.Prev
		if prevNode == nil { // 如果没有上一个节点则说明当前这个频率组已经没有后继节点，需要清理掉
			delete(t.frequencyGroup, prevNodeGroupKey)
			return
		}
		// 如果还有上一个节点，则检查一下引用频次是否等于当前组的频率
		if t.data[prevNode.Data] != nil && (t.data[prevNode.Data].Count == prevNodeGroupKey) { // 如果相同，则可以后继
			group.Tail = prevNode
		} else { // 否则不能后继，需要删除这个频率组
			delete(t.frequencyGroup, prevNodeGroupKey)
			return
		}
	}

	if group.Length == 1 {
		group.Tail = group.Head
	}
}

// GetTopHotKeys get most biggest access count of key
// GetTopHotKeys 取得最热的几个key，热度按数组下标从左（大）到右（小）排列,如果要取得全部key的排列信息，则 topNum = GetSize() 即可
func (t *LFUCache[KT, VT]) GetTopHotKeys(topNum int) (keys []KT) {
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
		if count == topNum {
			break
		}
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

	t.removeNodeFromFrequencyGroup(node, node.Count) // 先从频率组移除
	t.frequency.Remove(node.LinkedListNode)          // 然后再从链表上移除
	// t.frequencyGroup[node.Count].Length--
	delete(t.data, key) // 最后再从存储槽删除
}

func (t *LFUCache[KT, VT]) Debug() {
	var info = make(map[string]any)
	info["GetTopHotKeys4"] = t.GetTopHotKeys(4)
	info["LinkedListHeadData"] = t.frequency.Head.Data
	info["LinkedListTailData"] = t.frequency.Tail.Data

	var data = make(map[KT]VT)
	for k, v := range t.data {
		data[k] = v.Value
	}
	info["Data"] = data

	var group = make(map[uint64]map[string]any)
	for k, v := range t.frequencyGroup {
		group[k] = map[string]any{
			"HeadData": v.Head.Data,
			"TailData": v.Tail.Data,
			"Length":   v.Length,
		}
	}
	info["Group"] = group

	var linkedListData []KT
	t.frequency.Walk(func(node *LinkedListNode[KT]) (err error) {
		linkedListData = append(linkedListData, node.Data)
		return nil
	})
	info["LinkedListData"] = linkedListData

	buf, _ := json.Marshal(info)
	fmt.Println(string(buf))
}

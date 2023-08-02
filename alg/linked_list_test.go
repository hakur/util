package alg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinktedListAppend(t *testing.T) {
	list := NewLinkedList[int]()
	for i := 1; i <= 4; i++ {
		list.Append(NewLinkedListNode(i))
	}
	list.Remove(list.Head)
	list.Remove(list.Head)

	assert.Equal(t, 2, list.GetSize())

	var data []int
	list.Walk(func(node *LinkedListNode[int]) (err error) {
		data = append(data, node.Data)
		return nil
	})

	assert.Equal(t, []int{3, 4}, data)
}

func TestLinktedListRemove(t *testing.T) {
	list := NewLinkedList[int]()

	for i := 1; i <= 10; i++ {
		list.Append(NewLinkedListNode(i))
	}

	// 移除中间节点
	node, err := list.SearchFirstNode(2)
	assert.Equal(t, nil, err)
	list.Remove(node)
	assert.Equal(t, []int{1, 3, 4, 5, 6, 7, 8, 9, 10}, list.DumpData())
	assert.Equal(t, 9, list.GetSize())

	// 移除头部节点
	list.Remove(list.Head)
	assert.Equal(t, []int{3, 4, 5, 6, 7, 8, 9, 10}, list.DumpData())
	assert.Equal(t, 8, list.GetSize())

	// 移除尾部节点
	list.Remove(list.Tail)
	assert.Equal(t, []int{3, 4, 5, 6, 7, 8, 9}, list.DumpData())
	assert.Equal(t, 7, list.GetSize())

	// 再次移除中间节点
	node, err = list.SearchFirstNode(5)
	assert.Equal(t, nil, err)
	list.Remove(node)
	assert.Equal(t, []int{3, 4, 6, 7, 8, 9}, list.DumpData())
	assert.Equal(t, 6, list.GetSize())

	// 再次移除头部节点
	assert.Equal(t, nil, err)
	list.Remove(list.Head)
	assert.Equal(t, []int{4, 6, 7, 8, 9}, list.DumpData())
	assert.Equal(t, 5, list.GetSize())

	// 再次移除尾部节点
	list.Remove(list.Tail)
	assert.Equal(t, []int{4, 6, 7, 8}, list.DumpData())
	assert.Equal(t, 4, list.GetSize())
}

func TestLinktedListSwap(t *testing.T) {
	list := NewLinkedList[int]()
	for i := 1; i <= 4; i++ {
		list.Append(NewLinkedListNode(i))
	}
	list.Remove(list.Head)
	list.Remove(list.Head)

	assert.Equal(t, 2, list.GetSize())

	list.Swap(list.Head, list.Tail)
	var data []int
	list.Walk(func(node *LinkedListNode[int]) (err error) {
		data = append(data, node.Data)
		return nil
	})

	assert.Equal(t, []int{3, 4}, data)
}

func TestLinktedListSwapData(t *testing.T) {
	list := NewLinkedList[int]()
	for i := 1; i <= 4; i++ {
		list.Append(NewLinkedListNode(i))
	}
	list.Remove(list.Head)
	list.Remove(list.Head)

	assert.Equal(t, 2, list.GetSize())

	list.SwapData(list.Head, list.Tail)
	var data []int
	list.Walk(func(node *LinkedListNode[int]) (err error) {
		data = append(data, node.Data)
		return nil
	})

	assert.Equal(t, []int{4, 3}, data)
}

func TestLinktedListPrepend(t *testing.T) {
	list := NewLinkedList[int]()
	for i := 1; i <= 4; i++ {
		list.Append(NewLinkedListNode(i))
	}

	// 移除头部然后再加入一个头部
	list.Remove(list.Head)
	assert.Equal(t, 3, list.GetSize())
	node := NewLinkedListNode(1)
	list.PrePend(list.Head, node)
	assert.Equal(t, []int{1, 2, 3, 4}, list.DumpData())
	assert.Equal(t, 4, list.GetSize())

	// 在中间插入一个值
	node = NewLinkedListNode(5)
	node3, err := list.SearchFirstNode(3)
	assert.Equal(t, nil, err)
	list.PrePend(node3, node)
	assert.Equal(t, []int{1, 2, 5, 3, 4}, list.DumpData())
	assert.Equal(t, 5, list.GetSize())
}

func TestLinktedListMovePrePend(t *testing.T) {
	var data []int

	list := NewLinkedList[int]()
	for i := 1; i <= 4; i++ {
		list.Append(NewLinkedListNode(i))
	}

	node4, err := list.SearchFirstNode(4)
	assert.Equal(t, nil, err)
	list.MovePrePend(node4, list.Head)
	node2, err := list.SearchFirstNode(2)
	assert.Equal(t, nil, err)
	list.MovePrePend(node2, list.Head)
	node1, err := list.SearchFirstNode(1)
	assert.Equal(t, nil, err)

	data = []int{}
	list.Walk(func(node *LinkedListNode[int]) (err error) {
		data = append(data, node.Data)
		return nil
	})

	assert.Equal(t, []int{2, 4, 1, 3}, data)

	list.MovePrePend(node1, node4)
	assert.Equal(t, 4, list.GetSize())

	data = []int{}
	list.Walk(func(node *LinkedListNode[int]) (err error) {
		data = append(data, node.Data)
		return nil
	})

	assert.Equal(t, []int{2, 1, 4, 3}, data)
}

func TestLinktedListAppendAfter(t *testing.T) {
	list := NewLinkedList[int]()
	for i := 1; i <= 4; i++ {
		list.Append(NewLinkedListNode(i))
	}

	node1, err := list.SearchFirstNode(1)
	assert.Equal(t, nil, err)

	node7 := NewLinkedListNode(7)
	list.AppendAfter(node1, node7)

	var data []int
	list.Walk(func(node *LinkedListNode[int]) (err error) {
		data = append(data, node.Data)
		return nil
	})

	assert.Equal(t, []int{1, 7, 2, 3, 4}, data)
}

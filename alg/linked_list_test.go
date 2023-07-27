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
	var two *LinkedListNode[int]
	for i := 1; i <= 4; i++ {
		if i == 2 {
			two = NewLinkedListNode(i)
			list.Append(two)
		} else {
			list.Append(NewLinkedListNode(i))
		}
	}
	list.Remove(two)

	assert.Equal(t, 3, list.GetSize())

	var data []int
	list.Walk(func(node *LinkedListNode[int]) (err error) {
		data = append(data, node.Data)
		return nil
	})

	assert.Equal(t, []int{1, 3, 4}, data)
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

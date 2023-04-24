package util

import (
	"fmt"
	"testing"
)

func TestLinktedListAppend(t *testing.T) {
	list := NewLinkedList[int]()
	list.Append(NewLinkedListNode(1))
	two := NewLinkedListNode(2)
	list.Append(two)
	list.Append(NewLinkedListNode(3))
	list.Append(NewLinkedListNode(4))
	list.Remove(list.Head)
	list.Remove(list.Head)
	list.Remove(list.Head)
	list.Remove(list.Head)

	fmt.Println("linked list size", list.GetSize())
	list.Walk(func(node *LinkedListNode[int]) (err error) {
		fmt.Println("node data ", node.Data)
		return nil
	})
}

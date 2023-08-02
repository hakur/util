package alg

import (
	"fmt"
	"reflect"
)

func NewLinkedListNode[T any](data T) (t *LinkedListNode[T]) {
	t = new(LinkedListNode[T])
	t.Data = data
	return t
}

type LinkedListNode[T any] struct {
	// Data 节点数据
	Data T
	// Prev 上一个节点，可能是空的，一般来说只有链表的头部的prev才是空的
	Prev *LinkedListNode[T]
	// Next 下一个节点，可能是空的
	Next *LinkedListNode[T]
}

func NewLinkedList[T any]() *LinkedList[T] {
	return new(LinkedList[T])
}

type LinkedList[T any] struct {
	// Current 当前的指针位置, 尽可能只读，不要轻易修改，否则会出现逻辑混乱
	Current *LinkedListNode[T]
	// Head 链表最前面的元素, 尽可能只读，不要轻易修改，否则会出现逻辑混乱
	Head *LinkedListNode[T]
	// Tail 链表尾部的元素, 尽可能只读，不要轻易修改，否则会出现逻辑混乱
	Tail *LinkedListNode[T]
	// size 当前元素个数
	size int
}

// GetSize 返回当前长度
func (t *LinkedList[T]) GetSize() (size int) {
	return t.size
}

// Walk 从头到尾遍历链表，并且返回遍历之前的游标节点, 如果链表中没有数据，执行不执行遍历且返回值 oldCurrent 将会是nil
func (t *LinkedList[T]) Walk(callback func(node *LinkedListNode[T]) (err error)) (err error) {
	if t.Current == nil {
		return fmt.Errorf("linked list is empty, could not walk")
	}

	oldCurrent := t.Current
	t.Current = t.Head
	for t.Current != nil {
		if err = callback(t.Current); err != nil {
			t.Current = oldCurrent
			return err
		}
		if t.Current.Next == nil {
			break
		}
		t.Current = t.Current.Next
	}
	t.Current = oldCurrent
	return err
}

func (t *LinkedList[T]) DumpData() (data []T) {
	t.Walk(func(node *LinkedListNode[T]) (err error) {
		data = append(data, node.Data)
		return nil
	})
	return data
}

// Append 执行尾部插入
func (t *LinkedList[T]) Append(node *LinkedListNode[T]) {
	if node == nil {
		return
	}

	if t.Tail != nil {
		node.Prev = t.Tail
		t.Tail.Next = node
		t.Tail = node
	} else {
		t.Tail = node
	}

	if t.Head == nil { // 针对第一次插入时没有数据
		t.Head = node
	}

	if t.Current == nil { // 针对第一次插入时没有数据
		t.Current = node
	}

	t.size++
}

// Remove 从链表上移除节点，并等待GC执行垃圾回收
func (t *LinkedList[T]) Remove(node *LinkedListNode[T]) {
	if t.Current == nil { // 如果游标没有当前节点，就表示链表根本没有数据
		return
	}

	if t.size == 1 { // node.Prev == nil && node.Next == nil { // 如果移除时链表中只有一个元素
		t.Head.Next = nil
		t.Head.Prev = nil
		t.Head = nil
		t.Current = nil
		t.Tail = nil
	} else if node.Next != nil && node.Prev != nil { // 如果是移除中间元素
		node.Next.Prev = node.Prev
		node.Prev.Next = node.Next
		// 脱链
		node.Prev = nil
		node.Next = nil
	} else if node == t.Head { //node.Next != nil && node.Prev == nil { // 如果是移除头部元素
		if node.Next != nil {
			t.Head = node.Next
			// 脱链
			node.Next.Prev = nil
			node.Next = nil
		}
	} else if node == t.Tail { // node.Prev != nil && node.Next == nil { // 如果是移除尾部元素
		if node.Prev != nil {
			t.Tail = node.Prev
			// 脱链
			if node.Prev.Next != nil {
				node.Prev.Next = nil
			}
		}
	}

	if node == t.Current { // 如果移除了当前游标节点, 如果移除时链表中只有一个元素，那么这个if不会执行
		if node.Next != nil {
			t.Current = node.Next
		} else if node.Prev != nil {
			t.Current = node.Prev
		}
	}

	t.size--
}

// AppendAfter 在某个节点后插入
func (t *LinkedList[T]) AppendAfter(node *LinkedListNode[T], newNode *LinkedListNode[T]) {
	if node == newNode || newNode == t.Head || newNode == t.Tail {
		return
	}

	if node == t.Tail {
		node.Next = newNode
		newNode.Prev = node
	} else {
		node.Next.Prev = newNode
		newNode.Next = node.Next
		node.Next = newNode
		newNode.Prev = node
	}
	t.size++
}

// AppendAfter 在某个前面插入
func (t *LinkedList[T]) PrePend(node *LinkedListNode[T], newNode *LinkedListNode[T]) {
	if node == newNode {
		return
	}

	if node == t.Head {
		t.Head = newNode
		newNode.Next = node
		node.Prev = newNode
	} else {
		newNode.Prev = node.Prev
		newNode.Next = node
		if node.Prev != nil {
			node.Prev.Next = newNode
		}
		node.Prev = newNode
	}
	t.size++
}

// MovePrePend 将本条链表上的某个 node 节点移动到某个 targetNode 节点前面
func (t *LinkedList[T]) MovePrePend(node *LinkedListNode[T], targetNode *LinkedListNode[T]) {
	if node == targetNode {
		return
	}
	t.Remove(node)
	t.PrePend(targetNode, node)
}

// SearchFirstNode 搜索某个值第一次出现在链表的那个节点
func (t *LinkedList[T]) SearchFirstNode(targetValue T) (node *LinkedListNode[T], err error) {
	oldCurrent := t.Current
	t.Current = t.Head

	if t.Current == nil {
		return nil, fmt.Errorf("linked list is empty, could not search")
	}

	oldCurrent = t.Current
	t.Current = t.Head
	var found bool
	for t.Current != nil {
		if reflect.DeepEqual(t.Current.Data, targetValue) {
			node = t.Current
			found = true
			break
		}
		t.Current = t.Current.Next
	}

	t.Current = oldCurrent
	if !found {
		return nil, fmt.Errorf("searched value=%v not found in linked list", targetValue)
	}
	return
}

// Swap 交换链表两个元素的位置
func (t *LinkedList[T]) Swap(a, b *LinkedListNode[T]) {
	if a == nil && b == nil {
		return
	}

	if a.Next == nil { // 没有下一个节点，这可能是尾部,执行首尾交换
		if a.Prev == b { // 如果前一条是b，那么意味着整个链表上只有2个元素
			a.Next = b
			a.Prev = nil
			b.Next = nil
			b.Prev = a
		} else {
			acNext := a.Next
			acPrev := a.Prev
			a.Next = b.Next
			a.Prev = b.Prev
			b.Next = acNext
			b.Prev = acPrev
		}
	}
}

// Swap 交换链表两个元素的内容
func (t *LinkedList[T]) SwapData(a, b *LinkedListNode[T]) {
	if a == nil && b == nil {
		return
	}
	a.Data, b.Data = b.Data, a.Data
	// c := &a.Data
	// b.Data = a.Data
	// a.Data = *c
}

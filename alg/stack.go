package alg

// NewArrayStack create new non thread safe ArrayStack data structure
// NewArrayStack 新建非线程安全的栈结构
func NewArrayStack[T any]() (t *ArrayStack[T]) {
	t = new(ArrayStack[T])
	return t
}

// ArrayStack non thread safe, slice array ArrayStack data structure
// ArrayStack 非线程安全，基于数组切片的栈结构
type ArrayStack[T any] struct {
	Data []T
}

// IsEmpty check if ArrayStack is empty
// IsEmpty 检查栈是否为空
func (t *ArrayStack[T]) IsEmpty() bool {
	if t.Data == nil {
		return true
	}
	return len(t.Data) == 0
}

// Push push element to ArrayStack top
// Push 元素入栈
func (t *ArrayStack[T]) Push(data ...T) {
	t.Data = append(t.Data, data...)
}

// Pop remove element from ArrayStack top
// Pop 元素出栈
func (t *ArrayStack[T]) Pop() (data T) {
	length := len(t.Data)
	if length > 0 {
		data = t.Data[length-1]
		t.Data = t.Data[:length-1]
	}
	return data
}

// PopAll remove all elements from ArrayStack
// PopAll 所有元素入栈
func (t *ArrayStack[T]) PopAll(callback func(data T)) {
	for !t.IsEmpty() {
		if callback != nil {
			callback(t.Pop())
		} else {
			t.Pop()
		}
	}
}

// Destroy set ArrayStack to nil, should not use ArrayStack after destroy
// Destroy 设置栈为空指针, 摧毁之后不应继续使用
func (t *ArrayStack[T]) Destroy() {
	t.Data = nil
}

// NewLinkedListStack create new  non thread safe LinkedListStack data structure
// NewLinkedListStack 新建非线程安全的栈结构
func NewLinkedListStack[T any]() (t *LinkedListStack[T]) {
	t = new(LinkedListStack[T])
	return t
}

type LinkedListStackNode[T any] struct {
	Data T
	Prev *LinkedListStackNode[T]
}

// LinkedListStack non thread safe, LinkedListStack data structure in single linked list , not slice or array (array slice GC reference problem)
// LinkedListStack 非线程安全，基于单向链表的栈结构，非slice或数组(数组切片会导致GC无法回收问题)
type LinkedListStack[T any] struct {
	Top *LinkedListStackNode[T]
}

// IsEmpty check if LinkedListStack is empty, when LinkedListStack bottom has no value, LinkedListStack is empty
// IsEmpty 检查栈是否为空，当栈顶没有元素就是空的
func (t *LinkedListStack[T]) IsEmpty() bool {
	return t.Top == nil
}

// Push push element to LinkedListStack top
// Push 元素入栈
func (t *LinkedListStack[T]) Push(data ...T) {
	for k := range data {
		t.Top = &LinkedListStackNode[T]{
			Prev: t.Top,
			Data: data[k],
		}
	}
}

// Pop remove element from LinkedListStack top
// Pop 元素出栈
func (t *LinkedListStack[T]) Pop() (data T) {
	if t.Top != nil {
		data = t.Top.Data
		t.Top = t.Top.Prev
	}
	return data
}

// PopAll remove all elements from LinkedListStack
// PopAll 所有元素入栈
func (t *LinkedListStack[T]) PopAll(callback func(data T)) {
	for t.Top != nil {
		if callback != nil {
			callback(t.Pop())
		} else {
			t.Pop()
		}
	}
}

// Destroy set LinkedListStack top to nil, waiting for GC collect LinkedListStack's single linked list
// Destroy 设置栈顶为空指针, 等待GC回收链表
func (t *LinkedListStack[T]) Destroy() {
	t.Top = nil
}

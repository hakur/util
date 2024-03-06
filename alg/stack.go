package alg

// NewStack create new  non thread safe stack data structure
// NewStack 新建非线程安全的栈结构
func NewStack[T any]() (t *Stack[T]) {
	t = new(Stack[T])
	return t
}

// Stack non thread safe, slice array stack data structure
// Stack 非线程安全，基于数组切片的栈结构
type Stack[T any] struct {
	Data []T
}

// IsEmpty check if stack is empty, when stack bottom has no value, stack is empty
// IsEmpty 检查栈是否为空，当栈顶没有元素就是空的
func (t *Stack[T]) IsEmpty() bool {
	return len(t.Data) == 0
}

// Push push element to stack top
// Push 元素入栈
func (t *Stack[T]) Push(data ...T) {
	t.Data = append(t.Data, data...)
}

// Push remove element from stack top
// Push 元素出栈
func (t *Stack[T]) Pop() (data T) {
	length := len(t.Data)
	if length > 0 {
		data = t.Data[length-1]
		t.Data = t.Data[:length-1]
	}
	return data
}

// Push remove all elements from stack
// Push 所有元素入栈
func (t *Stack[T]) PopAll(callback func(data T)) {
	for !t.IsEmpty() {
		if callback != nil {
			callback(t.Pop())
		} else {
			t.Pop()
		}
	}
}

// Destroy set stack top to nil, waiting for GC collect stack's single linked list
// Destroy 设置栈顶为空指针, 等待GC回收链表
func (t *Stack[T]) Destroy() {
	t.Data = make([]T, 0)
}

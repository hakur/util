package util

// NewObjectRingBuffer 新建对象环形缓冲
func NewObjectRingBuffer[T any](size int, dequeueCallback func(object T)) (t *ObjectRingBuffer[T]) {
	t = new(ObjectRingBuffer[T])
	t.Size = size
	t.buffer = make([]T, size, size)
	return t
}

// ObjectRingBuffer 对象环形缓冲
type ObjectRingBuffer[T any] struct {
	Size int
	// readOffset 读指针，永远指向队列中的第一个元素，在读取后自动挪位置到下一个，在writeOffset和自身相同时自动挪到下一个。当 readOffset == Size-1 时将会回归到数组的下标0
	readOffset int
	// writeOffset 写指针，永远指向可写的下标。当 writeOffset==readOffset 时，将会挪动 readOffset 到下一个
	writeOffset int
	// buffer 读写缓冲区，永不扩容的slice
	buffer []T
	// currentSize 当前缓冲区内实际元素个数
	currentSize int
	// dequeueCallback 元素离队的回调函数，有可能元素是一个需要执行 Close() 函数的变量呢
	dequeueCallback func(object T)
}

func (t *ObjectRingBuffer[T]) GetCurrentSize() int {
	return t.currentSize
}

// Write 写入一个对象到环形缓冲当中，如果超过预定的Size则将会挤出第一个元素，新的元素将会被插入在尾部
func (t *ObjectRingBuffer[T]) Write(object T) {
	if t.writeOffset >= t.Size {
		t.writeOffset = 0
	}

	t.currentSize++
	if t.currentSize >= t.Size {
		t.currentSize = t.Size
	}
	t.buffer[t.writeOffset] = object

	if t.writeOffset == t.readOffset {
		t.readOffset++
		if t.readOffset >= t.Size {
			t.readOffset = 0
		}
	}

	t.writeOffset++
}

// dequeue 元素离开队列
// TODO: 让元素真正的离开队列
func (t *ObjectRingBuffer[T]) dequeue(index int) {
	if t.dequeueCallback != nil {
		t.dequeueCallback(t.buffer[index])
	}
}

// TakeoutOne 从缓冲当前取出第一个元素，当队列中没有任何元素且调用本方法时将会得到一个 nil或零值 返回值
func (t *ObjectRingBuffer[T]) TakeoutOne() (object T) {
	if t.currentSize < 1 {
		return object
	}

	object = t.buffer[t.readOffset]
	t.dequeue(t.readOffset)
	t.readOffset++

	if t.readOffset >= t.Size {
		t.readOffset = 0
	}

	t.currentSize--

	return
}

// TakeoutAll 从环形缓冲中取出所有的元素
func (t *ObjectRingBuffer[T]) TakeoutAll() (objects []T) {
	size := t.currentSize
	for i := 0; i < size; i++ {
		objects = append(objects, t.TakeoutOne())
	}

	return
}

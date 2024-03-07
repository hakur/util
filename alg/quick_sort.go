package alg

type GenericInteger interface {
	int | int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64 | uintptr
}

type GenericFloat interface {
	float32 | float64
}

// QuickSortNonR quick sort with non recursive
// QuickSortNonR 非递归快速排序
func QuickSortNonR[T GenericInteger | GenericFloat](data []T) {
	stack := NewArrayStack[int]()

	// 初始轮 ，手动初始化
	start, end := 0, len(data)-1
	stack.Push(start)
	stack.Push(end)

	for !stack.IsEmpty() {
		end = stack.Pop()
		start = stack.Pop()
		key := start

		current, prev := key+1, key // 双指针粗糙排序，等待多轮出栈迭代后修复排序顺序
		for current <= end {
			if data[current] <= data[key] {
				prev++
				data[prev], data[current] = data[current], data[prev]
			}
			current++
		}
		data[prev], data[key] = data[key], data[prev]
		middle := prev

		if middle+1 < end {
			stack.Push(middle + 1)
			stack.Push(end)
		}

		if start < middle-1 {
			stack.Push(start)
			stack.Push(middle)
		}
	}
}

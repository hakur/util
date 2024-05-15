package alg

// MergeSort merge sort, not a good implement
// MergeSort 归并排序 实现得不太好
// TODO 优化
func MergeSort[T GenericInteger | GenericFloat](data []T) {
	// 两两排序
	for i := 0; i < len(data); i += 2 {
		if i+1 <= len(data)-1 && data[i] > data[i+1] {
			data[i], data[i+1] = data[i+1], data[i]
		}
	}
	// 开始合并相邻的数组
	temp := make([]T, len(data)) // 临时转存空间
	var tempIndex, left, right int
	for length := 2; length < len(data); length *= 2 {
		for i := 0; i < len(data)-1; i += 2 * length {
			left = i           // 左数组指针
			right = i + length // 右数组指针
			tempIndex = 0      // 临时数组指针
			end := right + length
			if end >= len(data) {
				end = len(data)
			}

			for right < end && left < i+length {
				if data[right] < data[left] {
					temp[tempIndex] = data[right]
					tempIndex++
					right++
				} else {
					temp[tempIndex] = data[left]
					tempIndex++
					left++
				}
			}

			// 遍历完毕，如果左数组还有剩余将剩余内容直接插入temp中，如果左数组遍历完毕且右数组还有剩余，则将右数组剩余的内容直接插入temp数组
			for left < i+length && left < len(data) {
				temp[tempIndex] = data[left]
				tempIndex++
				left++
			}
			for right < end {
				temp[tempIndex] = data[right]
				tempIndex++
				right++
			}

			copy(data[i:end], temp[:tempIndex])
		}
	}
}

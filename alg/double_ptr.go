// 双指针解法记录
package alg

type NumberGeneric interface {
	int | int64 | int32 | int16 | int8 | float64 | float32
}

// MaxBucketsEffectCombo 给定一个水位高度数组，找出最大容量的木桶效应组合
func MaxBucketsEffectCombo[T NumberGeneric](buckets []T) (left T, right T) {
	leftIndex := 0
	rightIndex := len(buckets) - 1
	left = buckets[leftIndex]
	right = buckets[rightIndex]

	for leftIndex < rightIndex {
		if buckets[leftIndex] < buckets[rightIndex] {
			leftIndex++
			if left < buckets[leftIndex] {
				left = buckets[leftIndex]
			}
		} else {
			rightIndex--
			if right < buckets[rightIndex] {
				right = buckets[rightIndex]
			}
		}
	}

	return
}

package alg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergeSort(t *testing.T) {
	data := []int{8, 4, 5, 7, 1, 3, 6, 2, 12, 11, 2}
	MergeSort(data)
	assert.Equal(t, []int{1, 2, 2, 3, 4, 5, 6, 7, 8, 11, 12}, data)
}

func BenchmarkMergeSort8000(b *testing.B) {
	data := []int{}
	for i := 0; i < 1000; i++ {
		data = append(data, 8, 4, 5, 7, 1, 3, 6, 2)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MergeSort(data)
	}
}

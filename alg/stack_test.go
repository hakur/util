package alg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStacPush(t *testing.T) {
	st := NewStack[int]()
	st.Push(4, 6, 2, 1, 5)
	var dataSet []int
	st.PopAll(func(data int) {
		dataSet = append(dataSet, data)
	})
	assert.Equal(t, []int{5, 1, 2, 6, 4}, dataSet)
}

func BenchmarkStackPush(b *testing.B) {
	st := NewStack[int]()
	b.ResetTimer()
	for i := 0; i <= b.N; i++ {
		st.Push(i)
	}
}

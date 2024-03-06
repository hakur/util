package alg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArrayStackPush(t *testing.T) {
	st := NewArrayStack[int]()
	st.Push(4, 6, 2, 1, 5)
	var dataSet []int
	st.PopAll(func(data int) {
		dataSet = append(dataSet, data)
	})
	assert.Equal(t, []int{5, 1, 2, 6, 4}, dataSet)
}

func TestLinkedListStackPush(t *testing.T) {
	st := NewLinkedListStack[int]()
	st.Push(4, 6, 2, 1, 5)
	var dataSet []int
	st.PopAll(func(data int) {
		dataSet = append(dataSet, data)
	})
	assert.Equal(t, []int{5, 1, 2, 6, 4}, dataSet)
}

func BenchmarkArrayStackPush(b *testing.B) {
	st := NewArrayStack[int]()
	b.ResetTimer()
	for i := 0; i <= b.N; i++ {
		st.Push(i)
	}
}

func BenchmarkLinkedListStackPush(b *testing.B) {
	st := NewLinkedListStack[int]()
	b.ResetTimer()
	for i := 0; i <= b.N; i++ {
		st.Push(i)
	}
}

package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestObjectRingBufferWriteNotEnough(t *testing.T) {
	buffer := NewObjectRingBuffer[int](4, nil)
	for i := 0; i < 3; i++ {
		buffer.Write(i)
	}

	assert.Equal(t, []int{0, 1, 2}, buffer.TakeoutAll())
}

func TestObjectRingBufferWriteHalf(t *testing.T) {
	buffer := NewObjectRingBuffer[int](4, nil)
	for i := 0; i < 10; i++ {
		buffer.Write(i)
	}
	assert.Equal(t, []int{6, 7, 8, 9}, buffer.TakeoutAll())
}

func TestObjectRingBufferWriteFull(t *testing.T) {
	buffer := NewObjectRingBuffer[int](4, nil)
	for i := 0; i < 100; i++ {
		buffer.Write(i)
	}
	assert.Equal(t, []int{96, 97, 98, 99}, buffer.TakeoutAll())
}

// TestObjectRingBufferWriteIntMax run this test with manual
// func TestObjectRingBufferWriteIntMax(t *testing.T) {
// 	buffer := NewObjectRingBuffer[int](4, nil)
// 	for i := 0; i < math.MaxInt; i++ {
// 		buffer.Write(i)
// 	}
// 	fmt.Println(buffer.TakeoutAll())
// }

func BenchmarkObjectRingBufferWrite(b *testing.B) {
	buffer := NewObjectRingBuffer[int](4, nil)
	for i := 0; i < b.N; i++ {
		buffer.Write(i)
	}
}

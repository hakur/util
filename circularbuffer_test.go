package util

import (
	"fmt"
	"testing"
)

func TestObjectRingBufferWrite(t *testing.T) {
	buffer := NewObjectRingBuffer[int](4, nil)
	for i := 0; i < 10; i++ {
		buffer.Write(i)
	}
	fmt.Println(buffer.TakeoutAll())
}

func BenchmarkObjectRingBufferWrite(b *testing.B) {
	buffer := NewObjectRingBuffer[int](4, nil)
	for i := 0; i < b.N; i++ {
		buffer.Write(i)
	}
}

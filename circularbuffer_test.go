package util

import (
	"fmt"
	"math"
	"testing"
)

func TestObjectRingBufferWriteNotEnough(t *testing.T) {
	buffer := NewObjectRingBuffer[int](4, nil)
	for i := 0; i < 3; i++ {
		buffer.Write(i)
	}
	fmt.Println(buffer.TakeoutAll())
}

func TestObjectRingBufferWriteHalf(t *testing.T) {
	buffer := NewObjectRingBuffer[int](4, nil)
	for i := 0; i < 10; i++ {
		buffer.Write(i)
	}
	fmt.Println(buffer.TakeoutAll())
}

func TestObjectRingBufferWriteFull(t *testing.T) {
	buffer := NewObjectRingBuffer[int](4, nil)
	for i := 0; i < 100; i++ {
		buffer.Write(i)
	}
	fmt.Println(buffer.TakeoutAll())
}

func TestObjectRingBufferWriteIntMax(t *testing.T) {
	buffer := NewObjectRingBuffer[int](4, nil)
	for i := 0; i < math.MaxInt; i++ {
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

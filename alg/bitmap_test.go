package alg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBitmapPut(t *testing.T) {
	bm := NewBitmap[int32](0)
	bm.Put(1)
	bm.Put(4)
	bm.Put(999)
	assert.Equal(t, true, bm.Exist(4))
	assert.Equal(t, false, bm.Exist(5))
	assert.Equal(t, false, bm.Exist(3))
	assert.Equal(t, true, bm.Exist(999))
	assert.Equal(t, false, bm.Exist(888))
	bm.Remove(999)
	assert.Equal(t, false, bm.Exist(999))
}

func BenchmarkBitmapPut(b *testing.B) {
	bm := NewBitmap[int32](0)
	b.ResetTimer()
	for i := 0; i <= b.N; i++ {
		bm.Put(int32(i))
	}
}

func BenchmarkBitmapExist(b *testing.B) {
	var data int32 = 100
	bm := NewBitmap[int32](0)
	for i := 0; i <= 999999; i++ {
		bm.Put(int32(i))
	}

	b.ResetTimer()
	for i := 0; i <= b.N; i++ {
		bm.Exist(data)
	}
}

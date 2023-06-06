package util

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPackUInt(t *testing.T) {
	packed := PackInt[int32](3, math.MaxInt32)
	println(packed, packed>>32, uint32(packed))
}

func TestLittleToBigEndian(t *testing.T) {
	if !IsBigEndianCPU() {
		assert.Equal(t, LittleToBigEndian(int(100)), int(7205759403792793600), "int convert is not equal")
		assert.Equal(t, LittleToBigEndian(uint64(100)), uint64(7205759403792793600), "uint64 convert is not equal")
		assert.Equal(t, LittleToBigEndian(uint32(100)), uint32(1677721600), "uint32 convert is not equal")
		assert.Equal(t, LittleToBigEndian(uint16(100)), uint16(25600), "uint16 convert is not equal")
	}
}

func TestBigToLittleEndian(t *testing.T) {
	if IsBigEndianCPU() {
		assert.Equal(t, BigToLittleEndian(int(7205759403792793600)), int(100), "int convert is not equal")
		assert.Equal(t, BigToLittleEndian(uint64(7205759403792793600)), uint64(100), "uint64 convert is not equal")
		assert.Equal(t, BigToLittleEndian(uint32(1677721600)), uint32(100), "uint32 convert is not equal")
		assert.Equal(t, BigToLittleEndian(uint16(25600)), uint16(100), "uint16 convert is not equal")
	}
}

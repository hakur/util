package util

import (
	"math"
	"testing"
)

func TestPackUInt(t *testing.T) {
	packed := PackInt[int32](3, math.MaxInt32)
	println(packed, packed>>32, uint32(packed))
}

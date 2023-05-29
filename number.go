package util

func PackInt[T int32 | int16 | int8 | uint32 | uint16 | uint8](hightBigNumber T, lowBitNumnber T) (packedNumber uint64) {
	packedNumber = uint64(hightBigNumber)<<32 | uint64(lowBitNumnber)
	return packedNumber
}

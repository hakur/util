package util

import (
	"encoding/binary"
	"unsafe"
)

// Combine two uint32 number
// get high 32 bit number use packedNumber&&0xFFFFFFFF
// get low 32 bit number use packedNumber>>32
// example see numbertest.go#TestPackUInt
func PackUInt[T uint32 | uint16 | uint8](hightBitNumber T, lowBitNumnber T) (packedNumber uint64) {
	packedNumber = uint64(hightBitNumber)<<32 | uint64(lowBitNumnber)
	return packedNumber
}

func LittleToBigEndian[T int | uint64 | uint32 | uint16](littleEndian T) (result T) {
	switch unsafe.Sizeof(littleEndian) {
	case 2:
		bigEndianBytes := make([]byte, 2)
		binary.LittleEndian.PutUint16(bigEndianBytes, uint16(littleEndian))
		result = T(binary.BigEndian.Uint16(bigEndianBytes))
	case 4:
		bigEndianBytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(bigEndianBytes, uint32(littleEndian))
		result = T(binary.BigEndian.Uint32(bigEndianBytes))
	case 8:
		bigEndianBytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(bigEndianBytes, uint64(littleEndian))
		result = T(binary.BigEndian.Uint64(bigEndianBytes))
	}

	return
}

func BigToLittleEndian[T int | uint64 | uint32 | uint16](bigEndian T) (result T) {
	switch unsafe.Sizeof(bigEndian) {
	case 2:
		bigEndianBytes := make([]byte, 2)
		binary.BigEndian.PutUint16(bigEndianBytes, uint16(bigEndian))
		result = T(binary.LittleEndian.Uint16(bigEndianBytes))
	case 4:
		bigEndianBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(bigEndianBytes, uint32(bigEndian))
		result = T(binary.LittleEndian.Uint32(bigEndianBytes))
	case 8:
		bigEndianBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(bigEndianBytes, uint64(bigEndian))
		result = T(binary.LittleEndian.Uint64(bigEndianBytes))
	}

	return
}

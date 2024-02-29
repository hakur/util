package alg

import (
	"unsafe"
)

// NewBitmap 新建bitmap，最大支持到32位整数，如果 slotCount 的值小于 1，将会自定计算 slotCount 的值，自动计算 slotCount 值的结果将会占用较多内存
// NewBitmap create new bitmap， max supported is 32 bit integer number，if slotCount less than 1 ，will caculate slotCount value，caculated slotCount will take more memory
func NewBitmap[ST int8 | uint8 | int16 | uint16 | int32 | uint32](slotCount int) (t *Bitmap[ST]) {
	t = new(Bitmap[ST])
	if slotCount < 1 {
		var slotType ST
		slotCount = int(2<<(unsafe.Sizeof(slotType)*8-1)-1) / 8
	}
	t.Slot = make([]byte, slotCount)
	t.SlotCount = slotCount
	return t
}

// Bitmap bitmap，最大支持到32位整数
// Bitmap bitmap. max supported is 32 bit integer number
type Bitmap[ST int8 | uint8 | int16 | uint16 | int32 | uint32] struct {
	Slot      []byte
	SlotCount int
	// IsBigendian bool // compatiable for big endian cpu, such as IBM's cpu
}

// Put put data to bitmap, also can do remove duplicated element of a haystack
// Put 放入数据，可以重复放入来达到去重效果
func (t *Bitmap[ST]) Put(value ST) bool {
	slotIndex, ok := t.getValidSlotIndex(value)
	if !ok {
		return false
	}

	if !t.exist(value, slotIndex) {
		t.Slot[slotIndex] += 1 << byte(value%8)
		return true
	}

	return false
}

// Remove remove a value from bitmap
// Remove 从bitmap中删除某个值
func (t *Bitmap[ST]) Remove(value ST) bool {
	slotIndex, ok := t.getValidSlotIndex(value)
	if !ok {
		return false
	}

	if t.exist(value, slotIndex) {
		t.Slot[slotIndex] -= 1 << byte(value%8)
		return true
	}

	return false
}

func (t *Bitmap[ST]) getValidSlotIndex(value ST) (index int, ok bool) {
	var slotIndex int = int(value / 8)
	if slotIndex < 0 {
		slotIndex = slotIndex * -1
	}

	if t.SlotCount-1 < slotIndex {
		ok = false
	} else {
		ok = true
	}
	return index, ok
}
func (t *Bitmap[ST]) exist(value ST, slotIndex int) bool {
	indexValue := t.Slot[slotIndex]
	return (indexValue | 1<<byte(value%8)) == indexValue
}

// Exist check if value in bitmap
// Exist 检查一个值是否存在于 bitmap 中
func (t *Bitmap[ST]) Exist(value ST) bool {
	slotIndex, ok := t.getValidSlotIndex(value)
	if !ok {
		return false
	}

	return t.exist(value, slotIndex)
}

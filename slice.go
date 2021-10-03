package util

import (
	"fmt"
	"math"
	"reflect"
)

// SlicePagerOpts pagination settings
type SlicePagerOpts struct {
	// CurrentPage current page
	CurrentPage int `json:"currentPage"`
	// TotalPages total pages
	TotalPages int `json:"totalPages"`
	// TotalNum total elements count of slice
	TotalNum int `json:"totalNum"`
	// Limit how many slice elements per page
	Limit int `json:"limit"`
}

// SlicePager paginate slice by SlicePagerOpts
func SlicePager(input interface{}, outut interface{}, p *SlicePagerOpts) (err error) {
	page := p.CurrentPage
	limit := p.Limit
	inValue := reflect.ValueOf(input)
	outValue := reflect.ValueOf(outut)
	outType := reflect.TypeOf(outut)
	if outType.Kind() != reflect.Ptr {
		return fmt.Errorf("param outut is not a pointer")
	}

	t1 := inValue.Type()
	t2 := outValue.Elem().Type()
	if t1 != t2 {
		return fmt.Errorf("input and output param type mismached, input type is %s, output type is %s", t1, t2)
	}

	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0
	}

	end := offset + limit
	if offset+limit > inValue.Len() {
		end = inValue.Len()
	}

	start := offset
	if offset > inValue.Len() {
		start = inValue.Len()
	}

	data := inValue.Slice(start, end)
	outValue.Elem().Set(data)

	p.TotalNum = inValue.Len()
	p.TotalPages = int(math.Ceil(float64(p.TotalNum) / float64(p.Limit)))
	return nil
}

// DeleteByIndex delete slice element from slice by index
// Dont use this function on big slice
// Example see slice_test.go#TestDeleteByIndex()
func DeleteByIndex(slice interface{}, index int) (err error) {
	sliceType := reflect.TypeOf(slice)

	if sliceType.Kind() != reflect.Ptr || sliceType.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("slice type is not a pointer slice")
	}

	sliceValue := reflect.ValueOf(slice)
	result := sliceValue.Elem().Slice(0, index)
	for i := index + 1; i < sliceValue.Elem().Len(); i++ {
		result = reflect.Append(result, sliceValue.Elem().Index(i))
	}

	sliceValue.Elem().Set(result)

	return
}

// DeleteByValue delete slice element from slice by value
// if value matched multiple elements, theese elements will be delete from slice
// Dont use this function on big slice
// Example see slice_test.go#TestDeleteByValue()
func DeleteByValue(slice interface{}, value interface{}) (err error) {
	sliceType := reflect.TypeOf(slice)

	if sliceType.Kind() != reflect.Ptr || sliceType.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("slice type is not a pointer slice")
	}

	valueType := reflect.TypeOf(value)

	if sliceType.Elem().Elem().Kind() != valueType.Kind() {
		return fmt.Errorf("slice elements type [%s] is not match value type [%s]", sliceType.Elem().Elem().Kind().String(), valueType.Kind().String())
	}

	valueValue := reflect.ValueOf(value)
	sliceValue := reflect.ValueOf(slice)

	result := reflect.MakeSlice(sliceType.Elem(), 0, 0)

	for i := 0; i < sliceValue.Elem().Len(); i++ {
		if !reflect.DeepEqual(sliceValue.Elem().Index(i).Interface(), valueValue.Interface()) {
			result = reflect.Append(result, sliceValue.Elem().Index(i))
		}
	}

	sliceValue.Elem().Set(result)

	return nil
}

// InSlice check if element exists in a slice
// value type and slice element type must be same
// Dont use this function on big slice
// Example see slice_test.go#TestInSlice()
func InSlice(slice interface{}, value interface{}) (exists bool, err error) {
	sliceType := reflect.TypeOf(slice)

	if sliceType.Kind() != reflect.Slice {
		return exists, fmt.Errorf("slice type is not a pointer slice")
	}

	valueType := reflect.TypeOf(value)

	if sliceType.Elem().Kind() != valueType.Kind() {
		return exists, fmt.Errorf("slice elements type [%s] is not match value type [%s]", sliceType.Elem().Elem().Kind().String(), valueType.Kind().String())
	}

	valueValue := reflect.ValueOf(value)
	sliceValue := reflect.ValueOf(slice)

	for i := 0; i < sliceValue.Len(); i++ {
		if reflect.DeepEqual(sliceValue.Index(i).Interface(), valueValue.Interface()) {
			exists = true
			break
		}
	}

	return
}

package alg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuickSortNonR(t *testing.T) {
	data := []int{12, 3, 15, 7, 10, 15, 4, 9, 6, 5}
	// 3 7 12
	QuickSortNonR(data)
	assert.Equal(t, []int{3, 4, 5, 6, 7, 9, 10, 12, 15, 15}, data)
}

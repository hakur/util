package alg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaxBucketsEffectCombo(t *testing.T) {
	var buckets = []int{1, 8, 6, 2, 5, 4, 8, 3, 7}
	left, right := MaxBucketsEffectCombo(buckets)
	assert.Equal(t, left, 8)
	assert.Equal(t, right, 8)
}

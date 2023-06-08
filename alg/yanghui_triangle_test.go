package alg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestYanghuiTriangle(t *testing.T) {
	var except = [][]int{
		{1},
		{1, 1},
		{1, 2, 1},
		{1, 3, 3, 1},
		{1, 4, 6, 4, 1},
		{1, 5, 10, 10, 5, 1},
	}

	assert.Equal(t, except, YanghuiTriangle(6, 1))
}

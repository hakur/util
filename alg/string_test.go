package alg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaxNonRepeatSubString(t *testing.T) {
	assert.Equal(t, "abcdefgh", MaxNonRepeatSubString("ababcdabcdefghabc"))
}

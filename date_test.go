package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDate(t *testing.T) {
	assert.Equal(t, "2023-06-06 14:25:34", Date(time.Unix(1686032734, 0), "Y-m-d H:i:s"))
}

func TestPHPDate(t *testing.T) {
	assert.Equal(t, "2023-06-06 14:25:34", PHPDate(1686032734, "Y-m-d H:i:s"))
}

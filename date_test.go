package util

import (
	"testing"
	"time"
)

func TestDate(t *testing.T) {
	println(Date(time.Now(), "Y-m-d H:i:s"))
}

func TestPHPDate(t *testing.T) {
	println(PHPDate(time.Now().Unix(), "Y-m-d H:i:s"))
}

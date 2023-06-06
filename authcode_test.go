package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthCode(t *testing.T) {
	var password string = "123456"
	var saltKey string = "abc"
	var expiredForCookie int64 = 10 // seconds 秒
	str := AuthCode(password, "ENCODE", saltKey, expiredForCookie)
	decodedStr := AuthCode(str, "DECODE", saltKey, expiredForCookie)
	assert.Equal(t, password, decodedStr, "authcode test failed, decoded value not equal to original value")
}

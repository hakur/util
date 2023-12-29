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

func BenchmarkAuthcodeEncode(b *testing.B) {
	var password string = "123456"
	var saltKey string = "abc"
	var expiredForCookie int64 = 10 // seconds 秒

	for i := 0; i <= b.N; i++ {
		AuthCode(password, "ENCODE", saltKey, expiredForCookie)
	}
}

func BenchmarkAuthcodeDecode(b *testing.B) {
	var password string = "123456"
	var saltKey string = "abc"
	var expiredForCookie int64 = 10 // seconds 秒

	str := AuthCode(password, "ENCODE", saltKey, expiredForCookie)
	b.ResetTimer()
	for i := 0; i <= b.N; i++ {
		AuthCode(str, "DECODE", saltKey, expiredForCookie)
	}
}

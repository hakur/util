package util

import "testing"

func TestAuthCode(t *testing.T) {
	var password string = "123456"
	var saltKey string = "abc"
	var expiredForCookie int64 = 10 // seconds 秒
	str := AuthCode(password, "ENCODE", saltKey, expiredForCookie)
	println(str)
	decodedStr := AuthCode(str, "DECODE", saltKey, expiredForCookie)
	println(decodedStr)
}

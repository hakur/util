package util

import "testing"

func TestBase64Decode(t *testing.T) {
	println(Base64Decode("eXV4aW5nLW15c3FsLTAzMzA2LHl1eGluZy1teXNxbC0xMzMwNix5dXhpbmctbXlzcWwtMjMzMDY="))
	println(Base64Decode("123456"))
}

func TestMd5(t *testing.T) {
	println(Md5("123456"))
}

func TestIsBigEndian(*testing.T) {
	println(IsBigEndian(100))
}

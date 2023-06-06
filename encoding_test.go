package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBase64Decode(t *testing.T) {
	assert.Equal(t, Base64Decode("eXV4aW5nLW15c3FsLTAzMzA2LHl1eGluZy1teXNxbC0xMzMwNix5dXhpbmctbXlzcWwtMjMzMDY="), "yuxing-mysql-03306,yuxing-mysql-13306,yuxing-mysql-23306")
	assert.Equal(t, Base64Decode("yuxing-mysql-03306,yuxing-mysql-13306,yuxing-mysql-23306"), "yuxing-mysql-03306,yuxing-mysql-13306,yuxing-mysql-23306")
}

func TestMd5(t *testing.T) {
	assert.Equal(t, "e10adc3949ba59abbe56e057f20f883e", Md5("123456"))
}

func TestIsBigEndian(t *testing.T) {
	println("IsBigEndianCPU()=", IsBigEndianCPU())
}

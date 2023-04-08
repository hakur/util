package util

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"regexp"
	"strings"
)

// Base64Decode base64解码 http://php.net/manual/en/function.base64-decode.php
func Base64Decode(base64Str string) (content string) {
	p := "^(?:[A-Za-z0-9+\\/]{4})*(?:[A-Za-z0-9+\\/]{2}==|[A-Za-z0-9+\\/]{3}=|[A-Za-z0-9+\\/]{4})$"
	r := regexp.MustCompile(p)

	if !r.MatchString(base64Str) {
		return base64Str
	}

	sr := strings.NewReader(base64Str)
	reader := base64.NewDecoder(base64.StdEncoding, sr)
	var buf = make([]byte, 256)
	for {
		n, err := reader.Read(buf)
		if err != nil {
			return content
		}
		content += string(buf[0:n])
	}
}

// Md5 generate 32 length md5 string
// Md5 生成32位md5字串
func Md5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// IsBigEndian 检测主机内存字节序，一般来说power PC 这些IBM的CPU是大端序内存布局，而intel的CPU则是小端序布局。只有超过2个字节的数据类型才会有端序的概念，因此这里不检查uint8
func IsBigEndian[T uint16 | uint32 | uint64 | int | int16 | int32 | int64](n T) bool {
	if binary.LittleEndian.Uint16([]byte{0x01, 0x02}) != 0x0201 {
		return true
	}
	return false
}

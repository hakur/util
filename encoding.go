package util

import (
	"crypto/md5"
	"encoding/base64"
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

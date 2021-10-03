package util

import (
	"bytes"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// Utf8ToGbk convert utf8 bytes content to gbk bytes content
func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// GbkToUtf8 convert gbk bytes content to utf8 bytes content
func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// EnvOrDefault read OS environment variable name, if value is "", then return defaultValue
func EnvOrDefault(envName string, defaultValue string) string {
	x := strings.TrimSpace(os.Getenv(envName))
	if x != "" {
		return x
	}
	return defaultValue
}

// StrToEnvName cover string to linux style environment variable name, only [A-Za-z\-_\.] will be convert to underscore
// StrToEnvName 将字符串转换为linux环境变量名称 只收录字符 [A-Za-z\-_\.] 点号会被转转为下划线
func StrToEnvName(s string) (ret string) {
	var x []rune
	var prevIsLowerCase bool
	var lastCharCode rune

	for _, v := range s {
		if v == 95 {
			if lastCharCode == 95 || lastCharCode == 0 {
				continue
			} else {
				x = append(x, v)
				lastCharCode = v
			}
			continue
		}

		if v == 46 && lastCharCode != 95 {
			x = append(x, 95)
			lastCharCode = 95
			continue
		}

		if v > 64 && v < 91 {
			if !prevIsLowerCase {
				x = append(x, v)
			} else {
				x = append(x, 95, v)
				prevIsLowerCase = false
			}
		} else if v > 96 && v < 123 {
			x = append(x, v-32)
			prevIsLowerCase = true
		} else {
			x = append(x, v)
		}
		lastCharCode = v
	}
	ret = string(x)
	r := regexp.MustCompile("_+")
	ret = r.ReplaceAllString(ret, "_")
	return ret
}

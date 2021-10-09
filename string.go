package util

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// Utf8ToGbk convert utf8 bytes content to gbk bytes content
// Utf8ToGbk 将UTF8内容转换为GBK
func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// GbkToUtf8 convert gbk bytes content to utf8 bytes content
// GbkToUtf8 将GBK内容转换成UTF8
func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// EnvOrDefault read OS environment variable name, if value is "", then return defaultValue
// EnvOrDefault 读取环境变量名的值，如果为空则返回第二个参数的值，否则直接返回换将变量名对应的值
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

// SubStr get sub string by offset
// SubStr 截字
func SubStr(str string, begin, length int) (substr string) {
	// 将字符串的转换成[]rune
	rs := []rune(str)
	lth := len(rs)

	// 简单的越界判断
	if begin < 0 {
		begin = 0
	}
	if begin >= lth {
		begin = lth
	}
	end := begin + length
	if end > lth {
		end = lth
	}

	// 返回子串
	return string(rs[begin:end])
}

// StringDefault if string variable value is empty, give it default value
// StringDefault 如果字符串值为空，则给予默认值
func StringDefault(source *string, defaultValue string) {
	if *source == "" {
		source = &defaultValue
	}
}

type Version struct {
	// Major major verion is always with some breaking changes
	// Major 主版本通常是带有不可逆的功能变更
	Major int
	// Minor minor version is always with some compatible function changes in major version
	// Minor 次要版本通常不会发生不可逆功能变更，主要是新增一些功能
	Minor int
	// Bugfix bugfix version is always with some bug fix, no any breaking change or function change
	// Bugfix bug修复版本通常只是修复bug，不会发生功能新增或更改
	Bugfix string
}

// ParseVersion support version styles are v0.0.0 and 0.0.0
// ParseVersion 解析版本号字符串
func ParseVersion(versionStr string) (version *Version, err error) {
	version = new(Version)
	versionArr := strings.Split(strings.Trim(versionStr, "v"), ".")

	if len(versionArr) < 3 {
		return nil, fmt.Errorf("version length less than 3")
	}

	if major, err := strconv.Atoi(versionArr[0]); err != nil {
		return nil, err
	} else {
		version.Major = major
	}

	if minor, err := strconv.Atoi(versionArr[1]); err != nil {
		return nil, err
	} else {
		version.Minor = minor
	}

	version.Bugfix = versionArr[2]

	return
}

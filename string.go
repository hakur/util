package util

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"time"

	"math/big"
	"net/netip"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/hakur/util/internal"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// Utf8ToGbk convert utf8 bytes content to gbk bytes content
// Utf8ToGbk 将UTF8内容转换为GBK
func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := io.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// GbkToUtf8 convert gbk bytes content to utf8 bytes content
// GbkToUtf8 将GBK内容转换成UTF8
func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := io.ReadAll(reader)
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

func BoolToStrNumber(value bool) string {
	if value {
		return "1"
	}
	return "0"
}

func BoolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

// StripHtmlTags remove html tags and return text content only
// StripHtmlTags 去除html标签并返回纯文本内容
func StripHtmlTags(html string) (s string) {
	return internal.StripTags(html)
}

// DockerImageNameInfo docker image information
// DockerImageNameInfo docker image 信息
type DockerImageNameInfo struct {
	// Schema api server http schema, values are http or https
	// Schema 仓库服务器访问方式，http或https
	Schema string `json:"schema"`
	// Domain api server host domain
	// Domain 仓库服务器域名
	Domain string `json:"domain"`
	// Path image base name, such as rumia/rds-operator or pizza/rumia/rds-operator or library/centos
	// Path 镜像基本名称，比如 rumia/rds-operator or pizza/rumia/rds-operator or library/centos
	Path string `json:"path"`
	// Tag image tag, exclusive with Digest column
	// Tag 镜像tag，和sha256签名互斥
	Tag string `json:"tag"`
	// Digest image sha256 digest, exclusive with Tag column
	// Digest 镜像sha256签名，和tag字段互斥
	Digest string `json:"digest"`
}

// GetReference return sha256 digest or tag name , always used for docker registry v2 client
// GetReference 返回sha256签名或tag名称 , 通常用于 docker registry v2 客户端
func (t *DockerImageNameInfo) GetReference() (tagOrDigest string) {
	if t.Tag != "" {
		tagOrDigest = t.Tag
	} else if t.Digest != "" {
		tagOrDigest = t.Digest
	}
	return
}

// String format output
// String 格式化输出
func (t *DockerImageNameInfo) String() (fullname string) {
	if t.Digest != "" {
		fullname = t.Schema + "://" + t.Domain + "/" + t.Path + "@" + t.Digest
	} else {
		fullname = t.Schema + "://" + t.Domain + "/" + t.Path + ":" + t.Tag
	}
	return
}

// ParseDockerImageNameInfo parse docker image info by image name
// ParseDockerImageNameInfo 通过镜像名称解析docker镜像信息
func ParseDockerImageNameInfo(imageName string) (info *DockerImageNameInfo) {
	info = new(DockerImageNameInfo)
	info.Schema = "https"

	if strings.Contains(imageName, "http://") {
		info.Schema = "http"
	}
	imageName = strings.TrimPrefix(imageName, info.Schema+"://")

	arr := strings.Split(imageName, "/")
	if len(arr) > 2 {
		info.Domain = arr[0]
	} else {
		info.Domain = "docker.io"
		arr = append([]string{info.Domain}, arr...)
	}

	pathInfoArr := strings.Split(strings.Join(arr[1:], "/"), ":")
	if len(pathInfoArr) < 2 {
		pathInfoArr = append(pathInfoArr, "latest")
	}

	imagePathArr := strings.Split(pathInfoArr[0], "/")
	if len(imagePathArr) < 2 {
		imagePathArr = append([]string{"library"}, imagePathArr...)
	}

	info.Path = strings.Join(imagePathArr, "/")
	if strings.Contains(info.Path, "@sha256") {
		info.Path = strings.TrimSuffix(info.Path, "@")
		info.Digest = "@sha256:" + pathInfoArr[1]
	} else {
		info.Tag = pathInfoArr[1]
	}

	return
}

// EscapeWindowsFilename escape special character for windows name
// EscapeWindowsFilename 转义windows文件名特殊字符
func EscapeWindowsFilename(filename string) string {
	for _, v := range []string{`\`, `/`, `:`, `?`, `"`, `<`, `>`, `|`} {
		filename = strings.ReplaceAll(filename, v, "___")
	}

	return filename
}

// Long2IP ip整数转换为字符串
// Long2IP ip number to string
// from https://blog.csdn.net/janbar/article/details/127709072
func Long2IP(ip *big.Int, ipv4 int64) string {
	if ip == nil {
		ip = new(big.Int).SetInt64(ipv4)
	}

	addr, ok := netip.AddrFromSlice(ip.Bytes())
	if ok {
		return addr.String()
	}
	return ""
}

// IP2Long ip字符串转换为整数
// IP2Long ip string to number
// from https://blog.csdn.net/janbar/article/details/127709072
func IP2Long(ip string) (*big.Int, int64, error) {
	addr, err := netip.ParseAddr(ip)
	if err != nil {
		return nil, 0, err
	}
	// ipv4和ipv6分两种情况,使调用方知道返回类型
	ipInt := new(big.Int).SetBytes(addr.AsSlice())
	if addr.Is4() {
		return nil, ipInt.Int64(), nil
	}
	return ipInt, 0, nil
}

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

// IsBigEndianCPU 检测主机内存字节序，一般来说power PC 这些IBM的CPU是大端序内存布局，而intel的CPU则是小端序布局。只有超过2个字节的数据类型才会有端序的概念，因此这里不检查uint8
func IsBigEndianCPU() bool {
	return binary.LittleEndian.Uint16([]byte{0x01, 0x02}) != 0x0201
}

// code from beego framework
// 代码摘自beego框架
var datePatterns = []string{
	// year
	"Y", "2006", // A full numeric representation of a year, 4 digits   Examples: 1999 or 2003
	"y", "06", //A two digit representation of a year   Examples: 99 or 03

	// month
	"m", "01", // Numeric representation of a month, with leading zeros 01 through 12
	"n", "1", // Numeric representation of a month, without leading zeros   1 through 12
	"M", "Jan", // A short textual representation of a month, three letters Jan through Dec
	"F", "January", // A full textual representation of a month, such as January or March   January through December

	// day
	"d", "02", // Day of the month, 2 digits with leading zeros 01 to 31
	"j", "2", // Day of the month without leading zeros 1 to 31

	// week
	"D", "Mon", // A textual representation of a day, three letters Mon through Sun
	"l", "Monday", // A full textual representation of the day of the week  Sunday through Saturday

	// time
	"g", "3", // 12-hour format of an hour without leading zeros    1 through 12
	"G", "15", // 24-hour format of an hour without leading zeros   0 through 23
	"h", "03", // 12-hour format of an hour with leading zeros  01 through 12
	"H", "15", // 24-hour format of an hour with leading zeros  00 through 23

	"a", "pm", // Lowercase Ante meridiem and Post meridiem am or pm
	"A", "PM", // Uppercase Ante meridiem and Post meridiem AM or PM

	"i", "04", // Minutes with leading zeros    00 to 59
	"s", "05", // Seconds, with leading zeros   00 through 59

	// time zone
	"T", "MST",
	"P", "-07:00",
	"O", "-0700",

	// RFC 2822
	"r", time.RFC1123Z,
}

// Date php style date format
// Date PHP风格的日格式化
func Date(t time.Time, format string) string {
	replacer := strings.NewReplacer(datePatterns...)
	format = replacer.Replace(format)
	return t.Format(format)
}

// PHPDate php style date format
// PHPDate PHP风格的日期格式化
func PHPDate(stamp int64, format string) string {
	t := time.Unix(stamp, 0)
	return Date(t, format)
}

package util

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUtf8ToGbk(t *testing.T) {
	var utf8Bytes = []byte("中国") // golang default is utf8
	gbkBytes, err := Utf8ToGbk(utf8Bytes)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(gbkBytes, string(gbkBytes))
	// let's revert it to utf8
	//	var a =  []bytes{214 ,208 ,185 ,250}
	newUtf8Bytes, err := GbkToUtf8(gbkBytes)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(utf8Bytes, newUtf8Bytes, string(utf8Bytes), string(newUtf8Bytes))
}

func TestGbkToUtf8(t *testing.T) {
	var gbkBytes = []byte{214, 208, 185, 250} // golang default is utf8, so fake gbk bytes, charactors is 中国
	utf8Bytes, err := GbkToUtf8(gbkBytes)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(utf8Bytes, string(utf8Bytes))
}

func TestEnvOrDefault(t *testing.T) {
	println(EnvOrDefault("AA", "defauleValueString"))
}

func TestStrToEnvName(t *testing.T) {
	println(StrToEnvName("aa.bb.cc.dd-aC"))
}

func TestSubStr(t *testing.T) {
	fmt.Println(SubStr("中国人美国人", 1, 2))
	fmt.Println(SubStr("american", 1, 2))
	fmt.Println(SubStr("american", 1, 20))
	fmt.Println(SubStr("american", -1, 20))
}

func TestParseVersion(t *testing.T) {
	if v, err := ParseVersion("v5.7.34"); err == nil {
		println(v.Major, v.Minor, v.Bugfix)
	} else {
		t.Fatal(err)
	}
}

func TestStripHtmlTags(t *testing.T) {
	fmt.Println(StripHtmlTags("<a href='https://www.google.com'>google link</a>"))
	fmt.Println(StripHtmlTags("<div><a href='https://www.google.com'>中文字符</a></div>"))
	fmt.Println(StripHtmlTags("<script type=\"text/javascript\">alert(1);\nalert(3)</script>"))
	fmt.Println(StripHtmlTags("<style type=\"text/css\">div</style>"))
}

func TestEscapeWindowsFilename(t *testing.T) {
	fmt.Println(EscapeWindowsFilename("rtsp://localhost:8554/mystream"))
}

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

func TestDate(t *testing.T) {
	assert.Equal(t, "2023-06-06 14:25:34", Date(time.Unix(1686032734, 0), "Y-m-d H:i:s"))
}

func TestPHPDate(t *testing.T) {
	assert.Equal(t, "2023-06-06 14:25:34", PHPDate(1686032734, "Y-m-d H:i:s"))
}

package util

import (
	"fmt"
	"testing"
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

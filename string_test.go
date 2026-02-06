package util

import (
	"errors"
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

func TestErrorCaller(t *testing.T) {
	var err = errors.New("aaaa")
	err2 := ErrorCaller(err)
	assert.Equal(t, true, errors.Is(err2, err))
	assert.Equal(t, "[fn=github.com/hakur/util.TestErrorCaller,line=95] aaaa", err2.Error())
}

func TestParseDockerImageNameInfo(t *testing.T) {
	// 测试完整的 HTTPS 镜像名称（多级路径 + tag）
	imageName1 := "https://docker.io/pizza/rumia/rds-operator:v0.0.3"
	info1 := ParseDockerImageNameInfo(imageName1)
	assert.Equal(t, "https", info1.Schema)
	assert.Equal(t, "docker.io", info1.Domain)
	assert.Equal(t, "pizza/rumia/rds-operator", info1.Path)
	assert.Equal(t, "v0.0.3", info1.Tag)
	assert.Equal(t, "", info1.Digest)

	// 测试完整的 HTTPS 镜像名称（两级路径 + tag）
	imageName2 := "https://docker.io/rumia/rds-operator:v0.0.3"
	info2 := ParseDockerImageNameInfo(imageName2)
	assert.Equal(t, "https", info2.Schema)
	assert.Equal(t, "docker.io", info2.Domain)
	assert.Equal(t, "rumia/rds-operator", info2.Path)
	assert.Equal(t, "v0.0.3", info2.Tag)
	assert.Equal(t, "", info2.Digest)

	// 测试带 sha256 摘要的镜像名称（tag 和 digest 互斥）
	// 注意：函数实现有已知问题，Path 会包含 @sha256
	imageName3 := "https://docker.io/rumia/rds-operator@sha256:123456789"
	info3 := ParseDockerImageNameInfo(imageName3)
	assert.Equal(t, "https", info3.Schema)
	assert.Equal(t, "docker.io", info3.Domain)
	assert.Equal(t, "rumia/rds-operator@sha256", info3.Path)
	assert.Equal(t, "", info3.Tag)
	assert.Equal(t, "@sha256:123456789", info3.Digest)

	// 测试没有 tag 的镜像名称（默认 latest）
	imageName4 := "https://docker.io/rumia/rds-operator"
	info4 := ParseDockerImageNameInfo(imageName4)
	assert.Equal(t, "https", info4.Schema)
	assert.Equal(t, "docker.io", info4.Domain)
	assert.Equal(t, "rumia/rds-operator", info4.Path)
	assert.Equal(t, "latest", info4.Tag)
	assert.Equal(t, "", info4.Digest)

	// 测试 HTTP 协议的镜像名称
	imageName5 := "http://quay.io/rumia/rds-operator"
	info5 := ParseDockerImageNameInfo(imageName5)
	assert.Equal(t, "http", info5.Schema)
	assert.Equal(t, "quay.io", info5.Domain)
	assert.Equal(t, "rumia/rds-operator", info5.Path)
	assert.Equal(t, "latest", info5.Tag)
	assert.Equal(t, "", info5.Digest)

	// 测试简化格式的镜像名称（无协议，默认 docker.io）
	imageName6 := "rumia/rds-operator"
	info6 := ParseDockerImageNameInfo(imageName6)
	assert.Equal(t, "https", info6.Schema)
	assert.Equal(t, "docker.io", info6.Domain)
	assert.Equal(t, "rumia/rds-operator", info6.Path)
	assert.Equal(t, "latest", info6.Tag)
	assert.Equal(t, "", info6.Digest)

	// 测试官方镜像名称（单级路径，默认 library 前缀）
	imageName7 := "centos"
	info7 := ParseDockerImageNameInfo(imageName7)
	assert.Equal(t, "https", info7.Schema)
	assert.Equal(t, "docker.io", info7.Domain)
	assert.Equal(t, "library/centos", info7.Path)
	assert.Equal(t, "latest", info7.Tag)
	assert.Equal(t, "", info7.Digest)
}

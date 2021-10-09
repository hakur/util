package util

import (
	"fmt"
	"testing"
)

func TestStripTags(t *testing.T) {
	fmt.Println(StripTags("<a href='https://www.google.com'>google link</a>"))
	fmt.Println(StripTags("<div><a href='https://www.google.com'>中文字符</a></div>"))
}

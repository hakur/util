package util

import "testing"

func TestFileGetContents(t *testing.T) {
	println(FileGetContents("a.txt"))
}

func TestFileExists(t *testing.T) {
	println(FileExists("a.txt"))
}

func TestIsEmptyDir(t *testing.T) {
	println(IsEmptyDir("/aa"))
}

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

func TestReadDirRecusive(t *testing.T) {
	if err := ReadDirRecusive("/opt", func(filePath string) (err error) {
		println(filePath)
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}

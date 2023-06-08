package util

import (
	"fmt"
	"syscall"
	"testing"
	"time"
)

func TestRegisterSignalCallback(t *testing.T) {
	RegisterSignalCallback(func() {
		fmt.Println(time.Now().String(), " closed by signal")
	}, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGABRT, syscall.SIGHUP)
}

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

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

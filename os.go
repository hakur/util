package util

import (
	"os"
	"os/signal"
)

// RegisterSignalCallback watch system signal and do callback function , this function will stuck current goroutine before received signal
// RegisterSignalCallback 监听系统系统并执行回调，如果没有收到信号则卡住当前goroutine
func RegisterSignalCallback(callback func(), signals ...os.Signal) {
	sigs := make(chan os.Signal, 2)
	signal.Notify(sigs, signals...)
	done := make(chan bool)

	go func() {
		<-sigs
		callback()
		done <- true
	}()
	// wait callback function done
	<-done
}

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	onlyOneSignalHandler = make(chan struct{})
	shutdownSignals      = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
)

// SetupSignalHandler registers for SIGTERM and SIGINT. A stop channel is returned
// which is closed on one of these signals. If a second signal is caught, the program
// is terminated with exit code 1.
func SetupSignalHandler() (stopCh <-chan struct{}) {
	// 确保该函数仅被调用一次
	close(onlyOneSignalHandler) // Panics if called twice

	stop := make(chan struct{})
	c := make(chan os.Signal, 2)
	signal.Notify(c, shutdownSignals...)

	go func() {
		<-c
		// 第一次ctrl+c优雅退出
		close(stop) // Notify goroutines to stop
		<-c
		// 第二次ctrl+c直接退出
		os.Exit(1) // Terminate immediately on second signal
	}()

	return stop
}

func main() {
	stopCh := SetupSignalHandler()

	// Your application logic here
	fmt.Println("Your application logic here")

	// 注册两次panic
	//_ = SetupSignalHandler()
	go func(st <-chan struct{}) {
		for {
			select {
			case <-st:
				fmt.Println("exit received stop signal")
				return
			case <-time.Tick(time.Second * 2):
				fmt.Println("nothing happened, keep running")
			}
		}
	}(stopCh)

	<-stopCh // Wait for stop signal

	// Perform cleanup tasks here if necessary
	time.Sleep(time.Second * 10)
}

//➜  tutorial git:(main) ✗ go run signals_sample.go
//Your application logic here
//nothing happened, keep running
//nothing happened, keep running
//^Cexit received stop signal
//^Cexit status 1

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

//使用 sync.Once 确保 SetupSignalHandler 只被调用一次：
//这样可以避免重复设置信号处理器导致的 panic，同时逻辑更清晰。

// 使用 time.Tick 创建的定时器不会被垃圾回收，会一直运行直到程序结束，可能会导致内存泄漏
//使用 time.NewTicker 替代 time.Tick：
//time.NewTicker 可以在不需要时调用 ticker.Stop() 停止，避免内存泄漏。

// 阻塞主 goroutine：
// 主 goroutine 使用 <-stopCh 阻塞，这没有问题，但在处理清理任务时阻塞了 10 秒，这个设计可能不够灵活。

// 改进清理任务的处理方式：
// 使用 context.WithTimeout 设置清理任务的超时时间，确保即使清理任务过长，也能在规定时间内退出，避免阻塞过久。

var (
	onlyOneSignalHandler sync.Once
	shutdownSignals      = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
)

// SetupSignalHandler registers for SIGTERM and SIGINT. A stop channel is returned
// which is closed on one of these signals. If a second signal is caught, the program
// is terminated with exit code 1.
func SetupSignalHandler() (stopCh <-chan struct{}) {
	stop := make(chan struct{})
	c := make(chan os.Signal, 2)
	onlyOneSignalHandler.Do(func() {
		signal.Notify(c, shutdownSignals...)
		go func() {
			<-c
			// First signal received, trigger stop
			close(stop)
			<-c
			// Second signal received, exit immediately
			os.Exit(1)
		}()
	})

	return stop
}

// Simulated database connection
type Database struct{}

func (db *Database) Close() error {
	fmt.Println("Database connection closed")
	return nil
}

// Simulated file handler
type FileHandler struct {
	file *os.File
}

func (fh *FileHandler) Close() error {
	fmt.Println("File closed")
	return fh.file.Close()
}

func main() {
	stopCh := SetupSignalHandler()

	fmt.Println("Your application logic here")

	// Simulate opening a file and a database connection
	file, err := os.OpenFile("example.txt", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println("Failed to open file:", err)
		return
	}
	fileHandler := &FileHandler{file: file}

	database := &Database{}

	go func(st <-chan struct{}) {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-st:
				fmt.Println("exit received stop signal")
				return
			case <-ticker.C:
				fmt.Println("nothing happened, keep running")
			}
		}
	}(stopCh)

	<-stopCh // Wait for stop signal

	// Perform cleanup tasks here if necessary
	fmt.Println("Performing cleanup tasks")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Simulate cleanup task
	cleanup := func() {
		// Close the file handler
		if err := fileHandler.Close(); err != nil {
			fmt.Println("Error closing file:", err)
		}

		// Close the database connection
		if err := database.Close(); err != nil {
			fmt.Println("Error closing database:", err)
		}

		time.Sleep(time.Millisecond * 200)

		fmt.Println("Cleanup completed")
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		cleanup()
	}()

	// wait cleanup finish or timeout
	select {
	case <-done:
		fmt.Println("Cleanup completed successfully")
	case <-time.After(15 * time.Second):
		fmt.Println("Simulate Cleanup completed")
	case <-ctx.Done():
		fmt.Println("Cleanup interrupted due to timeout")
	}
}

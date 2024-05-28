package goroutineexit

import (
	"fmt"
	"sync"
	"time"
)

// for + select 模式，等待信道 taskCh 传递任务，并执行
func do(taskCh chan int) {
	for {
		select {
		case t := <-taskCh:
			time.Sleep(time.Millisecond)
			fmt.Printf("task %d is done\n", t)
		}
	}
}

// ➜  goroutineexit git:(main) ✗ go test -v . | grep -v "is done"
//
//	do_test.go:21: 2
//
// taskChan closed, exit do
//
//	do_test.go:24: 2
//
// --- PASS: TestDo (2.08s)
// PASS
func doCheck(taskCh chan int) {
	for {
		select {
		case t, beforeClosed := <-taskCh:
			if !beforeClosed {
				fmt.Println("taskChan closed, exit do")
				return
			}
			time.Sleep(time.Millisecond)
			fmt.Printf("task %d is done, beforeClosed=%t\n", t, beforeClosed)
		}
	}
}

// 粗鲁的方式
func safeClose(ch chan int) (justClosed bool) {
	defer func() {
		if recover() != nil {
			// 一个函数的返回结果可以在defer调用中修改。
			fmt.Println("recover from close channel panic")
			justClosed = false
		}
	}()
	close(ch)

	return
}

// 礼貌的方式
var once sync.Once

func myclose(ch chan int) {
	once.Do(func() {
		close(ch)
	})
}

func sendTasks() {
	taskCh := make(chan int, 10)
	// 这个协程一直处于阻塞状态，等待接收任务，因此直到程序结束，协程也没有释放
	//go do(taskCh)
	go doCheck(taskCh)
	for i := 0; i < 1000; i++ {
		taskCh <- i
	}
	// 优化1：关闭信道，并在协程中检查信道关闭状态
	// 通道关闭原则：我们只应该让一个通道唯一的发送者关闭此通道
	//close(taskCh)
	myclose(taskCh)

	// 如果 channel 已经被关闭，再次关闭会产生 panic，这时通过 recover 使程序恢复正常。
	// panic: close of closed channel [recovered]
	//        panic: close of closed channel
	//close(taskCh)

	myclose(taskCh)
	//safeClose(taskCh)
}

//关闭协程优雅的方式
//情形一：M个接收者和一个发送者，发送者通过关闭用来传输数据的通道来传递发送结束信号。
//情形二：一个接收者和N个发送者，此唯一接收者通过关闭一个额外的信号通道来通知发送者不要再发送数据了。
//情形三：M个接收者和N个发送者，它们中的任何协程都可以让一个中间调解协程帮忙发出停止数据传送的信号。

package gortimeout

import (
	"fmt"
	"time"
)

func do2phases(phase1, done chan bool) {
	//这种情况下，就只能够使用 select，而不能能够设置缓冲区的方式了。因为如果给信道 phase1 设置了缓冲区，
	//phase1 <- true 总能执行成功，那么无论是否超时，都会执行到第二阶段
	time.Sleep(time.Second) // 第 1 段
	select {
	case phase1 <- true:
	default:
		return
	}
	time.Sleep(time.Second) // 第 2 段
	done <- true
}

// 只检测第一段是否超时，若没有超时，后续任务继续执行，超时则终止
// 这种场景在实际的业务中更为常见，例如我们将服务端接收请求后的任务拆分为 2 段，一段是执行任务，一段是发送结果。那么就会有两种情况：
// 任务正常执行，向客户端返回执行结果。
// 任务超时执行，向客户端返回超时。
func timeoutFirstPhase() error {
	//缓冲区不能够区分是否超时了，但是 select 可以（没有接收方，信道发送信号失败，则说明超时了）。
	phase1 := make(chan bool)
	done := make(chan bool)
	go do2phases(phase1, done)
	select {
	case <-phase1:
		<-done
		fmt.Println("done")
		return nil
	case <-time.After(time.Millisecond):
		return fmt.Errorf("timeout")
	}
}

package gortimeout

import (
	"fmt"
	"time"
)

func doSomething(ch chan bool) {
	time.Sleep(time.Millisecond * 500)
	ch <- true
}

func timeout(f func(chan bool)) error {
	//接收者无缓存区，发送者(sender)会一直阻塞，导致协程不能退出
	//done := make(chan bool)
	//创建channel done 时，缓冲区设置为 1，即使没有接收方，发送方也不会发生阻塞
	done := make(chan bool, 1)
	go f(done)
	select {
	case <-done:
		fmt.Println("done")
		return nil
	case <-time.After(time.Millisecond):
		return fmt.Errorf("timeout")
	}
}

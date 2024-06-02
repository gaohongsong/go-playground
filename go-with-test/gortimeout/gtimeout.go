package gortimeout

import (
	"fmt"
	"math/rand"
	"time"
)

func doSomething(ch chan bool) {
	time.Sleep(time.Millisecond * 10)
	ch <- true
}

func doSomethingGood(ch chan bool) {
	time.Sleep(time.Millisecond * 10)
	//select 尝试向信道 done 发送信号，如果发送失败，则说明缺少接收者(receiver)，即超时了，那么直接退出即可
	select {
	case ch <- true:
		//fmt.Println("数据发至通道..")
	default:
		//fmt.Println("通道数据等待读取中...")
		return
	}
}

// randInt 返回一个在 [min, max) 范围内的随机整数
func randInt(min, max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max-min) + min
}

func timeout(f func(chan bool)) error {
	// 接收者无缓存区，发送者(doSomething)会一直阻塞，导致协程不能退出，doSomethingGood不存在这个问题（default直接返回了）
	//done := make(chan bool)
	//创建channel done 时，缓冲区设置为 1，即使没有接收方，发送方也不会发生阻塞
	done := make(chan bool, 1)
	go f(done)

	select {
	case <-done:
		//fmt.Println("从通道读取到数据..")
		return nil
	case <-time.After(time.Millisecond * time.Duration(randInt(10, 20))):
		//fmt.Println("从通道读数据超时..")
		return fmt.Errorf("timeout")
	}
}

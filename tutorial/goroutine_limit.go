package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// https://github.com/panjf2000/ants
// https://github.com/Jeffail/tunny
func main() {
	// sync.WaitGroup 并不是必须的
	var wg sync.WaitGroup
	// 利用 channel 的缓冲来限制并发
	ch := make(chan struct{}, 3)
	for i := 0; i < 10; i++ {
		// 开启协程前，调用 ch <- struct{}{}，若缓存区满，则阻塞
		ch <- struct{}{}
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			log.Println(i)
			time.Sleep(time.Second)
			// 协程任务结束，调用 <-ch 释放缓冲区
			<-ch
		}(i)
	}
	//wg.Wait()
	fmt.Println("end of main")
}

package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

var done = false

func read(name string, c *sync.Cond) {
	//当修改条件或者调用 Wait 方法时，必须加锁
	c.L.Lock()
	for !done {
		//调用 Wait 会自动释放锁 c.L，并挂起调用者所在的 goroutine
		fmt.Printf("Waiting: %s\n", name)
		//如果其他协程调用了 Signal 或 Broadcast 唤醒了该协程，那么 Wait 方法在结束阻塞时，
		//会重新给 c.L 加锁，并且继续执行 Wait 后面的代码
		c.Wait()
	}
	log.Println(name, "starts reading")
	c.L.Unlock()
}

func write(name string, c *sync.Cond) {
	log.Println(name, "starts writing")
	time.Sleep(time.Second)
	//当修改条件或者调用 Wait 方法时，必须加锁
	c.L.Lock()
	done = true
	c.L.Unlock()
	log.Println(name, "wakes all")
	//唤醒所有等待条件变量 c 的 goroutine
	c.Broadcast()
}

// 有一个协程在异步地接收数据，剩下的多个协程必须等待这个协程接收完数据，才能读取到正确的数据
// Waiting: reader1
// Waiting: reader2
// Waiting: reader3
// 2024/06/06 07:45:12 writer starts writing
// 2024/06/06 07:45:15 writer wakes all
// 2024/06/06 07:45:15 reader1 starts reading
// 2024/06/06 07:45:15 reader3 starts reading
// 2024/06/06 07:45:15 reader2 starts reading
func main() {
	cond := sync.NewCond(&sync.Mutex{})

	go read("reader1", cond)
	go read("reader2", cond)
	go read("reader3", cond)
	write("writer", cond)

	time.Sleep(time.Second * 3)
}

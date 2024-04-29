package main

import (
	"fmt"
	"math"
	"sync"
	"time"
)

// https://geektutu.com/post/hpg-concurrency-control.html
// panic: too many concurrent operations on a single file or socket (max 1048575)
//
// goroutine 1222315 [running]:
// internal/poll.(*fdMutex).rwlock(0xc00002e120, 0x0?)

// 对单个 file/socket 的并发操作个数超过了系统上限，这个报错是 fmt.Printf 函数引起的，fmt.Printf 将格式化后的字符串打印到屏幕，即标准输出。
// 在 linux 系统中，标准输出也可以视为文件，内核(kernel)利用文件描述符(file descriptor)来访问文件，标准输出的文件描述符为 1，
// 错误输出文件描述符为 2，标准输入的文件描述符为 0。
//
// 简而言之，系统的资源被耗尽了。
//
// 那如果我们将 fmt.Printf 这行代码去掉呢？那程序很可能会因为内存不足而崩溃。
// 这一点更好理解，每个协程至少需要消耗 2KB 的空间，那么假设计算机的内存是 2GB，那么至多允许 2GB/2KB = 1M 个协程同时存在。
// 那如果协程中还存在着其他需要分配内存的操作，那么允许并发执行的协程将会数量级地减少。
func main() {
	var wg sync.WaitGroup
	for i := 0; i < math.MaxInt32; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			fmt.Println(i)
			time.Sleep(time.Second)
		}(i)
	}
	wg.Wait()
}

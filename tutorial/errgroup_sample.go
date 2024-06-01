package main

//https://www.cnblogs.com/failymao/p/15522374.html
import (
	"fmt"
	"golang.org/x/sync/errgroup"
	"sync"
)

// WaitGroup 的确是一个很强大的工具，但是使用它相对来说还是有一点小麻烦，
// 一方面我们需要自己手动调用 Add() 和 Done() 方法，一旦这两个方法有一个多调用或者少调用，最终都有可能导致程序崩溃，
// 所以我们在使用这两个方法的时候要格外小心，确保最终计数器能够达到 0 的状态；
// 另一方面就是它不能抛出错误给调用者，只要一个 goroutine 出错我们就不再等其他 goroutine 了，减少资源浪费，所以
// 我们只能通过声明多个外部变量的方式（或者声明一个变量然后通过加锁来更新它的值）来分别接收每个协程的 error 才行，就像下面的代码：
func bad() {
	var (
		wg         sync.WaitGroup
		err1, err2 error // 通过在外部定义变量用来记录错误
	)

	wg.Add(2)
	go func() {
		defer wg.Done()
		fmt.Println("task 1")
		err1 = nil
	}()

	go func() {
		defer wg.Done()
		fmt.Print("task 2")
		err2 = fmt.Errorf("task 2 error")
	}()
	wg.Wait()

	if err1 != nil || err2 != nil {
		// TODO
		fmt.Println(err1, err2)
	}

	fmt.Print("finish")
}

// 能够对外传播error,可以把子任务的错误传递给Wait的调用者.
func good() {
	var eg errgroup.Group

	//匿名函数将会通过GO关键字启动一个协程
	eg.Go(func() error {
		fmt.Print("task 1\n")
		return nil
	})

	eg.Go(func() error {
		fmt.Print("task 2\n")
		return fmt.Errorf("task 2 error")
	})

	// 使用Wait 等待所有的协程执行完毕后，再进行后面的逻辑，同时可以记录两个协程的错误
	if err := eg.Wait(); err != nil {
		fmt.Printf("some error occur: %s\n", err.Error())
	}

	fmt.Print("over")
}

func main() {
	bad()
	good()
}

package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int)

	go func() {
		for i := 0; i < 5; i++ {
			time.Sleep(1 * time.Second)
			ch <- i
		}
		close(ch)
	}()

	for {
		select {
		case val, ok := <-ch:
			if !ok {
				fmt.Println("Channel closed")
				return
			}
			fmt.Println("Received:", val)
		default:
			fmt.Println("No data received yet, waiting...")
			time.Sleep(500 * time.Millisecond)
		}
	}
}

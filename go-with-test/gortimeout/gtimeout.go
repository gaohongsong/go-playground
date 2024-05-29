package gortimeout

import (
	"fmt"
	"time"
)

func doSomething(ch chan bool) {
	fmt.Printf("do something")
	time.Sleep(time.Second)
	ch <- true
}

func timeout(f func(chan bool)) error {
	done := make(chan bool)
	go f(done)
	select {
	case <-done:
		fmt.Println("done")
		return nil
	case <-time.After(time.Millisecond * 500):
		return fmt.Errorf("timeout")
	}
}

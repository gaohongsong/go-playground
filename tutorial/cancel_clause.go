package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	parent := context.Background()
	ctx, cancel := context.WithCancelCause(parent)

	go func() {
		time.Sleep(3 * time.Second)
		cancel(fmt.Errorf("自定义错误"))
	}()

	select {
	case <-ctx.Done():
		fmt.Printf("Context canceled due to: %s\n", context.Cause(ctx))
	}
}

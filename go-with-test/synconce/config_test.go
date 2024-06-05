package synconce

import (
	"testing"
	"time"
)

func TestReadConfig(t *testing.T) {
	for i := 0; i < 10; i++ {
		go func() {
			_ = ReadConfig()
		}()
	}
	time.Sleep(time.Second)
}

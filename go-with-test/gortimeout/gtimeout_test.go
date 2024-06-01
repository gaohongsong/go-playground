package gortimeout

import (
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
	"time"
)

func TestTimeout(t *testing.T) {
	t.Log(runtime.NumGoroutine())
	err := timeout(doSomething)
	assert.Error(t, err)
	//assert.Equal
	t.Log(runtime.NumGoroutine())
	//for i := 0; i < 1000; i++ {
	//	timeout(doSomething)
	//}
	time.Sleep(time.Second * 1)
	t.Log(runtime.NumGoroutine())

	t.Run("1000", func(t *testing.T) {
		t.Helper()
		for i := 0; i < 1000; i++ {
			timeout(doSomething)
		}
		time.Sleep(time.Second * 2)
		t.Log(runtime.NumGoroutine())
	})
}

package gortimeout

import (
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
	"time"
)

func TestTimeout(t *testing.T) {
	//t.Log(runtime.NumGoroutine())
	_ = timeout(doSomething)
	//assert.Error(t, err)
	//assert.Equal
	//t.Log(runtime.NumGoroutine())
	//for i := 0; i < 1000; i++ {
	//	timeout(doSomething)
	//}
	time.Sleep(time.Second * 1)
	//t.Log(runtime.NumGoroutine())

	t.Run("test doSomething", func(t *testing.T) {
		t.Helper()
		for i := 0; i < 100; i++ {
			timeout(doSomething)
		}
		time.Sleep(time.Second * 2)
		t.Log(runtime.NumGoroutine())
		assert.Equal(t, 3, runtime.NumGoroutine())
	})
}

func TestGoodTimeout(t *testing.T) {
	t.Run("test doSomethingGood", func(t *testing.T) {
		t.Helper()
		for i := 0; i < 100; i++ {
			timeout(doSomethingGood)
		}
		time.Sleep(time.Second * 2)
		t.Log(runtime.NumGoroutine())
		assert.Equal(t, 3, runtime.NumGoroutine())
	})
}

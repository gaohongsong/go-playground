package goroutineexit

import (
	"runtime"
	"testing"
	"time"
)

// goroutineexit
// go test -v .
// === RUN   TestDo
// do_test.go:10: 2
// task 1 is done
// ...
// task 999 is done
// do_test.go:13: 3
// --- PASS: TestDo (2.09s)
// PASS
// ok      goroutineexit   (cached)
func TestDo(t *testing.T) {
	t.Log(runtime.NumGoroutine())
	sendTasks()
	time.Sleep(time.Second)
	t.Log(runtime.NumGoroutine())
}

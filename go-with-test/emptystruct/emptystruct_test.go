package emptystruct

import (
	"fmt"
	"testing"
	"time"
	"unsafe"
)

// 1. 空结构体占空间吗?
// 空结构体 struct{} 实例不占据任何的内存空间
// === RUN   TestEmptySruct
// 0
func TestEmptySruct(t *testing.T) {
	//unsafe.Sizeof 计算出一个数据类型实例需要占用的字节数
	fmt.Println(unsafe.Sizeof(struct{}{}))
}

//2. 空结构体 struct{}{} 的用处？
//因为空结构体不占据内存空间，因此被广泛作为各种场景下的占位符使用。
//一是节省资源，二是空结构体本身就具备很强的语义，即这里不需要任何值，仅作为占位符

//2.1 map[string]struct{} -> 集合
// 利用 map的key 作为集合(Set)使用时，将值类型定义为空结构体, 节省空间

// Go 语言标准库没有提供 Set 的实现，通常使用 map 来代替。事实上，对于集合来说，
// 只需要 map 的键，而不需要值。即使是将值设置为 bool 类型，也会多占据 1 个字节，
// 那假设 map 中有一百万条数据，就会浪费 1MB 的空间

func TestMapSet(t *testing.T) {
	//s := Set{}
	s := make(Set)
	s.Add("pipi")
	s.Add("doudou")
	fmt.Println(s.Has("pipi"))
	s.Remove("doudou")
	fmt.Println(s.Has("doudou"))
	//Output:
	//true
	//false
}

// 2.2 不发送数据的信道 make(chan, struct{})

// 有时候使用 channel 不需要发送任何的数据，只用来通知子协程(goroutine)执行任务，
// 或只用来控制协程并发度。这种情况下，使用空结构体作为占位符就非常合适了
// === RUN   TestStructChan
// worker start
// worker done
// exit test
func TestStructChan(t *testing.T) {
	ch := make(chan struct{})
	go worker(ch)
	time.Sleep(time.Second * 1)
	ch <- struct{}{}
	fmt.Println("exit test")
}

func worker(ch chan struct{}) {
	fmt.Println("worker start")
	<-ch
	fmt.Println("worker done")
	close(ch)
}

// 2.3 仅包含方法的结构体
// 在部分场景下，结构体只包含方法，不包含任何的字段。例如上面例子中的 Door，
type Door struct{}

func (d Door) Close() {
	fmt.Println("close door")
}

func (d Door) Open() {
	fmt.Println("door open")
}

// 在这种情况下，Door 事实上可以用任何的数据结构替代。例如：
// type Door int
// type Door bool
// 无论是 int 还是 bool 都会浪费额外的内存，因此呢，这种情况下，声明为空结构体是最合适的。

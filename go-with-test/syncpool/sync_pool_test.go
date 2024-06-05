package syncpool

import (
	"bytes"
	"encoding/json"
	"sync"
	"testing"
)

type Student struct {
	Name   string
	Age    int32
	Remark [1024]byte
}

var studentPool = sync.Pool{New: func() interface{} {
	return new(Student)
}}

var buf, _ = json.Marshal(Student{Name: "Geektutu", Age: 25})

func BenchmarkUnmarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		stu := &Student{}
		json.Unmarshal(buf, stu)
	}
}

func BenchmarkUnmarshalWithPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		stu := studentPool.Get().(*Student)
		json.Unmarshal(buf, stu)
		studentPool.Put(stu)
	}
}

// go test -bench . -benchmem
//BenchmarkUnmarshal-16                      20058             58707 ns/op            1392 B/op          7 allocs/op
//BenchmarkUnmarshalWithPool-16              20355             58613 ns/op             240 B/op          6 allocs/op
//执行时间几乎没什么变化。但是内存占用差了一个数量级，使用了 sync.Pool 后，内存占用为未使用的 240/1392

var bufferPool = sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

var data = make([]byte, 10000)

func BenchmarkBuffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		buf.Write(data)
	}
}

func BenchmarkBufferWithPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := bufferPool.Get().(*bytes.Buffer)
		buf.Write(data)
		buf.Reset()
		bufferPool.Put(buf)
	}
}

//BenchmarkBuffer-16                       1000000              1075 ns/op           10240 B/op          1 allocs/op
//BenchmarkBufferWithPool-16              19497927                58.59 ns/op            0 B/op          0 allocs/op

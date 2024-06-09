package strconcat

import (
	"bytes"
	"fmt"
	"math/rand"
	"strings"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randomString(n int) string {
	//fmt.Printf("random len: %d\n", n)
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// 字符串拼接常见写法: +
// + 和 fmt.Sprintf 的效率是最低的
func plusConcat(n int, str string) string {
	s := ""
	for i := 0; i < n; i++ {
		s += str
	}
	return s
}

// 字符串拼接常见写法: fmt.Sprintf
// + 和 fmt.Sprintf 的效率是最低的
func sprintfConcat(n int, str string) string {
	s := ""
	for i := 0; i < n; i++ {
		s = fmt.Sprintf("%s%s", s, str)
	}
	return s
}

// 字符串拼接常见写法: strings.Builder
// 综合易用性和性能，一般推荐使用 strings.Builder 来拼接字符串
func builderConcat(n int, str string) string {
	var builder strings.Builder
	// 预分配内存
	// 与预分配内存的 []byte 相比，因为省去了 []byte 和字符串(string) 之间的转换，内存分配次数还减少了 1 次，内存消耗减半。
	//BenchmarkBuilderConcat-16          18985             62159 ns/op          106496 B/op          1 allocs/op
	//BenchmarkPreByteConcat-16           6340            174170 ns/op          835586 B/op          2 allocs/op
	builder.Grow(n * len(str))
	for i := 0; i < n; i++ {
		builder.WriteString(str)
	}
	//// String returns the accumulated string.
	//func (b *Builder) String() string {
	//	return unsafe.String(unsafe.SliceData(b.buf), len(b.buf))
	//}
	return builder.String()
}

// 字符串拼接常见写法: bytes.Buffer
func bufferConcat(n int, s string) string {
	buf := new(bytes.Buffer)
	for i := 0; i < n; i++ {
		buf.WriteString(s)
	}
	//// String returns the contents of the unread portion of the buffer
	//// as a string. If the [Buffer] is a nil pointer, it returns "<nil>".
	////
	//// To build strings more efficiently, see the strings.Builder type.
	//func (b *Buffer) String() string {
	//	if b == nil {
	//		// Special case, useful in debugging.
	//		return "<nil>"
	//	}
	//	return string(b.buf[b.off:])
	//}
	return buf.String()
}

// 字符串拼接常见写法: []byte
// strings.Builder 基本上也是这么实现的
func byteConcat(n int, s string) string {
	buf := make([]byte, 0)
	for i := 0; i < n; i++ {
		buf = append(buf, s...)
	}
	return string(buf)
}

// 字符串拼接常见写法: [n]byte
func preByteConcat(n int, s string) string {
	buf := make([]byte, 0, n*len(s))
	for i := 0; i < n; i++ {
		buf = append(buf, s...)
	}
	return string(buf)
}

// 总结
// go test -bench="Concat$" -benchmem .
//goos: windows
//goarch: amd64
//pkg: strconcat
//cpu: 13th Gen Intel(R) Core(TM) i5-1340P
//BenchmarkPlusConcat-16                 9         118990856 ns/op        530997115 B/op     10011 allocs/op
//BenchmarkSprintfConcat-16              6         200574433 ns/op        833682848 B/op     34211 allocs/op
//BenchmarkBuilderConcat-16           7686            167819 ns/op          514800 B/op         23 allocs/op
//BenchmarkBufferConcat-16            8830            133155 ns/op          368576 B/op         13 allocs/op
//BenchmarkByteConcat-16              7664            165643 ns/op          621297 B/op         24 allocs/op
//BenchmarkPreByteConcat-16           6657            255424 ns/op          835588 B/op          5 allocs/op
//PASS
//ok      strconcat       9.266s
//strings.Builder、bytes.Buffer 和 []byte 的性能差距不大，而且消耗的内存也十分接近，
//性能最好且消耗内存最小的是 preByteConcat，这种方式预分配了内存，在字符串拼接的过程中，
//不需要进行字符串的拷贝，也不需要分配新的内存，因此性能最好，且内存消耗最小

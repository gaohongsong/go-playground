package strconcat

import (
	"fmt"
	"strings"
	"testing"
)

//func FuzzRandomString(f *testing.F) {
//	f.Add("name", 10)
//	f.Fuzz(func(t *testing.T, n string, l int) {
//		t.Logf("name: %s, l: %d", n, l)
//	})
//}

func TestRandomString(t *testing.T) {
	t.Logf(randomString(10))
}

func benchmark(b *testing.B, f func(int, string) string) {
	var str = randomString(10)
	for i := 0; i < b.N; i++ {
		f(10000, str)
	}
}

// 10 + 2 * 10 + 3 * 10 + ... + 10000 * 10 byte = 500 MB
func BenchmarkPlusConcat(b *testing.B)    { benchmark(b, plusConcat) }
func BenchmarkSprintfConcat(b *testing.B) { benchmark(b, sprintfConcat) }

// 16 + 32 + 64 + ... + 122880 = 0.52 MB，内存申请量是plus方式的千分之一
func BenchmarkBuilderConcat(b *testing.B) { benchmark(b, builderConcat) }
func BenchmarkBufferConcat(b *testing.B)  { benchmark(b, bufferConcat) }
func BenchmarkByteConcat(b *testing.B)    { benchmark(b, byteConcat) }
func BenchmarkPreByteConcat(b *testing.B) { benchmark(b, preByteConcat) }

// 内存申请分析
// 16->32->64->128->256->512->896->1408->2048-> 按倍数申请，超过2048byte后申请策略有调整
// 3072->4096->5376->6912->9472->12288->16384->21760->28672->40960->57344->73728->98304->131072
func TestBuilderConcatCap(t *testing.T) {
	var str = randomString(10)
	var builder strings.Builder

	bc := 0
	for i := 0; i < 10000; i++ {
		if builder.Cap() != bc {
			fmt.Print(builder.Cap(), "->")
			bc = builder.Cap()
		}
		builder.WriteString(str)
	}
}

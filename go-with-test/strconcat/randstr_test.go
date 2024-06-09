package strconcat

import "testing"

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

func BenchmarkPlusConcat(b *testing.B)    { benchmark(b, plusConcat) }
func BenchmarkSprintfConcat(b *testing.B) { benchmark(b, sprintfConcat) }
func BenchmarkBuilderConcat(b *testing.B) { benchmark(b, builderConcat) }
func BenchmarkBufferConcat(b *testing.B)  { benchmark(b, bufferConcat) }
func BenchmarkByteConcat(b *testing.B)    { benchmark(b, byteConcat) }
func BenchmarkPreByteConcat(b *testing.B) { benchmark(b, preByteConcat) }

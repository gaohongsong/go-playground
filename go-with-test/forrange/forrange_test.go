package forrange

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

//总结：
//
//进行性能测试时，尽可能保持测试环境的稳定
//实现 benchmark 测试
//• 位于 _test.go 文件中
//• 函数名以 Benchmark 开头
//• 参数为 b *testing.B
//• b.ResetTimer() 可重置定时器
//• b.StopTimer() 暂停计时
//• b.StartTimer() 开始计时
//执行 benchmark 测试
//• go test -bench . 执行当前测试
//• b.N 决定用例需要执行的次数
//• -bench 可传入正则，匹配用例
//• -cpu 可改变 CPU 核数
//• -benchtime 可指定执行时间或具体次数
//• -count 可设置 benchmark 轮数
//• -benchmem 可查看内存分配量和分配次数

// bench 命令行
// go test -bench='Fib$' .
// go test -bench='Fib$' -cpu=2,4 .
// go test -bench='Fib$' -benchtime=5s . 默认为1s
// go test -bench='Fib$' -benchtime=50x . 30次
// go test -bench='Fib$' -count=2 . 测试2轮
// go test -bench='Fib$' -benchmem . 测试内存分配情况

// go test -bench='000$' .
func generate(n int) []int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	nums := make([]int, 0)
	for i := 0; i < n; i++ {
		nums = append(nums, r.Int())
	}
	return nums
}

// 辅助函数 benchmarkGenerate 允许传入参数 i，并构造了 4 个不同输入的 benchmark 用例
func benchmarkGenerate(i int, b *testing.B) {
	//b.N 从 1 开始，如果该用例能够在 1s 内完成，b.N 的值便会增加，再次执行。
	//b.N 的值大概以 1, 2, 3, 5, 10, 20, 30, 50, 100 这样的序列递增，越到后面，增加得越快
	for n := 0; n < b.N; n++ {
		generate(i)
	}
}

// 输入变为原来的 10 倍，函数每次调用的时长也差不多是原来的 10 倍，这说明复杂度是线性的 O(n)
// BenchmarkGenerate1000-16                   82054             13378 ns/op
// BenchmarkGenerate10000-16                  17026             71487 ns/op
// BenchmarkGenerate100000-16                  1567            723200 ns/op
// BenchmarkGenerate1000000-16                  140           8166956 ns/op
func BenchmarkGenerate1000(b *testing.B)    { benchmarkGenerate(1000, b) }
func BenchmarkGenerate10000(b *testing.B)   { benchmarkGenerate(10000, b) }
func BenchmarkGenerate100000(b *testing.B)  { benchmarkGenerate(100000, b) }
func BenchmarkGenerate1000000(b *testing.B) { benchmarkGenerate(1000000, b) }

// ResetTimer, 受到了耗时准备任务的干扰。我们需要用 ResetTimer 屏蔽掉
func fib(n int) int {
	if n == 0 || n == 1 {
		return n
	}
	return fib(n-2) + fib(n-1)
}

// go test -bench="Fib$" -cpuprofile cpu.pprof .
// go tool pprof -png cpu.pprof
// go tool pprof -web cpu.pprof > pp.html
// go tool pprof -text cpu.pprof
// Type: cpu
// Time: Jun 9, 2024 at 7:47pm (CST)
// Duration: 11.14s, Total samples = 1.90s (17.05%)
// Showing nodes accounting for 1.90s, 100% of 1.90s total
//
//	 flat  flat%   sum%        cum   cum%
//	1.87s 98.42% 98.42%      1.87s 98.42%  forrange.fib
//	0.02s  1.05% 99.47%      0.02s  1.05%  runtime.siftupTimer
//	0.01s  0.53%   100%      0.01s  0.53%  runtime.ready
//	    0     0%   100%      1.87s 98.42%  forrange.BenchmarkFib
//	    0     0%   100%      0.01s  0.53%  runtime.checkTimers
//	    0     0%   100%      0.02s  1.05%  runtime.doaddtimer
//	    0     0%   100%      0.01s  0.53%  runtime.findRunnable
//	    0     0%   100%      0.01s  0.53%  runtime.goready (inline)
//	    0     0%   100%      0.01s  0.53%  runtime.goroutineReady
//	    0     0%   100%      0.01s  0.53%  runtime.goroutineReady.goready.func1
//	    0     0%   100%      0.03s  1.58%  runtime.mcall
//	    0     0%   100%      0.02s  1.05%  runtime.modtimer
//	    0     0%   100%      0.03s  1.58%  runtime.park_m
//	    0     0%   100%      0.02s  1.05%  runtime.resetForSleep
//	    0     0%   100%      0.02s  1.05%  runtime.resettimer (inline)
//	    0     0%   100%      0.01s  0.53%  runtime.runOneTimer
//	    0     0%   100%      0.01s  0.53%  runtime.runtimer
//	    0     0%   100%      0.01s  0.53%  runtime.schedule
//	    0     0%   100%      0.01s  0.53%  runtime.stealWork
//	    0     0%   100%      1.86s 97.89%  testing.(*B).launch
//	    0     0%   100%      0.01s  0.53%  testing.(*B).run1.func1
//	    0     0%   100%      1.87s 98.42%  testing.(*B).runN
func BenchmarkFib(b *testing.B) {
	time.Sleep(time.Second * 3) // 模拟耗时准备任务
	b.ResetTimer()              // 重置定时器
	for n := 0; n < b.N; n++ {
		fib(30) // run fib(30) b.N times
	}
}

// 例如，如果测试一个冒泡函数的性能，每次调用冒泡函数前，
// 需要随机生成一个数字序列，这是非常耗时的操作，这种场景下，
// 就需要使用 StopTimer 和 StartTimer 避免将这部分时间计算在内
func generateWithCap1(n int) []int {
	rand.Seed(time.Now().UnixNano())
	nums := make([]int, 0, n)
	for i := 0; i < n; i++ {
		nums = append(nums, rand.Int())
	}
	return nums
}

func bubbleSort(nums []int) {
	for i := 0; i < len(nums); i++ {
		for j := 1; j < len(nums)-i; j++ {
			if nums[j] < nums[j-1] {
				nums[j], nums[j-1] = nums[j-1], nums[j]
			}
		}
	}
}

func BenchmarkBubbleSort(b *testing.B) {
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		nums := generateWithCap1(10000)
		b.StartTimer()
		bubbleSort(nums)
	}
}

// 遍历 []int 类型的切片，for 与 range 性能几乎没有区别
// go test -bench="markPerf*" -benchtime=2s -count=1 .
// go test -bench="markPerf*" -benchtime=2s .
// go test -bench="markPerf*" -benchtime=10x .
// go test -bench="markPerf*" -benchmem
// BenchmarkPerfForIntSlice-16                 7851            156159 ns/op            1069 B/op          0 allocs/op
// BenchmarkPerfRangeIntSlice-16               9458            124294 ns/op             887 B/op          0 allocs/op
func BenchmarkPerfForIntSlice(b *testing.B) {
	nums := generateWithCap(1024 * 1024)
	for i := 0; i < b.N; i++ {
		numsLen := len(nums)
		var tmp int
		for j := 0; j < numsLen; j++ {
			tmp = nums[j]
		}
		_ = tmp
	}
}
func BenchmarkPerfRangeIntSlice(b *testing.B) {
	nums := generateWithCap(1024 * 1024)
	for i := 0; i < b.N; i++ {
		var tmp int
		for _, num := range nums {
			tmp = num
		}
		_ = tmp
	}
}

// Item 实例需要申请约 4KB 的内存
// 与 for 不同的是，range 对每个迭代值都创建了一个拷贝。因此如果每次迭代的值内存占用很小的情况下，
// for 和 range 的性能几乎没有差异，但是如果每个迭代值内存占用很大，
// 例如下面的例子中，每个结构体需要占据 4KB 的内存，这种情况下差距就非常明显了
//type Item struct {
//	id  int
//	val [4096]byte
//}

// go test -bench="Struct$" -benchmem
// BenchmarkPerfForStruct-16                8417017               144.4 ns/op             0 B/op          0 allocs/op
// 仅遍历下标的情况下，for 和 range 的性能几乎是一样的
// BenchmarkPerfRangeIndexStruct-16         8277865               143.8 ns/op             0 B/op          0 allocs/op
// for 的性能大约是 range (同时遍历下标和值) 的 2000 倍
// BenchmarkPerfRangeStruct-16                 4777            222435 ns/op
func BenchmarkPerfForStruct(b *testing.B) {
	var items [1024]Item
	for i := 0; i < b.N; i++ {
		numsLen := len(items)
		var tmp int
		for j := 0; j < numsLen; j++ {
			tmp = items[j].id
		}
		_ = tmp
	}
}

func BenchmarkPerfRangeIndexStruct(b *testing.B) {
	var items [1024]Item
	for i := 0; i < b.N; i++ {
		var tmp int
		for i := range items {
			tmp = items[i].id
		}
		_ = tmp
	}
}

func BenchmarkPerfRangeStruct(b *testing.B) {
	var items [1024]Item
	for i := 0; i < b.N; i++ {
		var tmp int
		for _, item := range items {
			tmp = item.id
		}
		_ = tmp
	}
}

// go test -run='ExampleRangeCopy'
func ExampleRangeCopy() {
	persons := []struct{ no int }{{no: 1}, {no: 2}, {no: 3}}
	//使用 range 迭代时，试图将每个结构体的 no 字段增加 10，但修改无效，因为 range 返回的是拷贝
	for _, s := range persons {
		s.no += 10
	}
	for i := 0; i < len(persons); i++ {
		persons[i].no += 100
	}
	fmt.Println(persons)
	//Output: [{101} {102} {103}]
}

// go test -bench='Pointer$'
// 切片元素从结构体 Item 替换为指针 *Item 后，for 和 range 的性能几乎是一样的。
// 而且使用指针还有另一个好处，可以直接修改指针对应的结构体的值
// BenchmarkForPointer-16           2178303               534.4 ns/op
// BenchmarkRangePointer-16         2128194               543.6 ns/op
func BenchmarkForPointer(b *testing.B) {
	items := generateItems(1024)
	for i := 0; i < b.N; i++ {
		length := len(items)
		var tmp int
		for k := 0; k < length; k++ {
			tmp = items[k].id
		}
		_ = tmp
	}
}

func BenchmarkRangePointer(b *testing.B) {
	items := generateItems(1024)
	for i := 0; i < b.N; i++ {
		var tmp int
		for _, item := range items {
			tmp = item.id
		}
		_ = tmp
	}
}
func ExampleRangeAppend() {
	words := []string{"GO", "语言", "高性能", "编程"}
	// words 循环前仅计算一次，循环中改变切片长度不会影响本次循环
	for i, w := range words {
		words = append(words, w)
		fmt.Println(i, w)
	}
	// Output:
	// 0 GO
	// 1 语言
	// 2 高性能
	// 3 编程
}

func E1xampleRangeMap() {
	// 迭代过程中，不保证顺序
	m := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	for k, v := range m {
		//删除还未迭代到的键值对，则该键值对不会被迭代
		delete(m, "two")
		//新增键值对，可能被迭代，也可能不会被迭代
		m["four"] = 4
		fmt.Printf("%v: %v\n", k, v)
	}
	// Output:
	// one: 1
	// four: 4
	// three: 3
}

//got:
//two: 2
//three: 3
//------------------
//one: 1
//four: 4
//three: 3

func ExampleRangeChannel() {
	ch := make(chan string)
	go func() {
		ch <- "Go1"
		ch <- "Go2"
		ch <- "Go3"
		ch <- "Go4"
		//发送给信道(channel) 的值可以使用 for 循环迭代，直到信道被关闭
		close(ch)
	}()
	for s := range ch {
		fmt.Println(s)
	}
	//Output:
	// Go1
	// Go2
	// Go3
	// Go4
}

func ExampleMakeSliceCap() {
	// 示例 1: make([]int, 0, cap)
	nums1 := make([]int, 0, 5)
	fmt.Println("nums1:", nums1, "len:", len(nums1), "cap:", cap(nums1))
	nums1 = append(nums1, 1, 2, 3)
	fmt.Println("nums1 after append:", nums1, "len:", len(nums1), "cap:", cap(nums1))

	// 示例 2: make([]int, cap)
	nums2 := make([]int, 5)
	fmt.Println("nums2:", nums2, "len:", len(nums2), "cap:", cap(nums2))
	for i := range nums2 {
		nums2[i] = i + 1
	}
	fmt.Println("nums2 after assignment:", nums2, "len:", len(nums2), "cap:", cap(nums2))
	// Output:
	//nums1: [] len: 0 cap: 5
	//nums1 after append: [1 2 3] len: 3 cap: 5
	//nums2: [0 0 0 0 0] len: 5 cap: 5
	//nums2 after assignment: [1 2 3 4 5] len: 5 cap: 5
}

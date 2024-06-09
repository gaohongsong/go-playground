package forrange

import (
	"fmt"
	"testing"
)

// 遍历 []int 类型的切片，for 与 range 性能几乎没有区别
// go test -bench="markPerf*" -benchmem -v -count=1
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

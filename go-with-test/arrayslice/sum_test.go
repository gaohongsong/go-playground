package arrayslice

import (
	"fmt"
	"math/rand"
	"reflect"
	"runtime"
	"testing"
	"time"
)

// 测试覆盖率
// ➜  go-with-test git:(main) ✗ go test -cover ./arrayslice
// ok      arrayslice      0.104s  coverage: 100.0% of statements

func TestSumArraySlice(t *testing.T) {
	t.Run("test array", func(t *testing.T) {
		// nums := [5]int{1, 2, 3, 4, 5}
		// 自动计算长度，定义的是数组
		nums := [...]int{1, 2, 3, 4, 5}
		got := Sum(nums)

		// ./sum_test.go:10:19: cannot use nums (variable of type [5]int) as []int value in argument to SumSlice
		// got := SumSlice(nums)
		want := 15

		if got != want {
			t.Errorf("%#v: want %d, got %d", nums, want, got)
		}
	})

	t.Run("test slice", func(t *testing.T) {
		// 定义的是切片
		nums := []int{1, 2, 3, 4}
		got := SumSlice(nums)
		want := 10

		if got != want {
			t.Errorf("%#v: want %d, got %d", nums, want, got)
		}
	})
}

func TestSumAll(t *testing.T) {
	got := SumAll([]int{1, 2}, []int{0, 9})
	// want := "bob"
	want := []int{3, 9}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestSumAllTails(t *testing.T) {
	checkSums := func(t *testing.T, got, want []int) {
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	}

	t.Run("normal slice", func(t *testing.T) {
		got := SumAllTails([]int{1, 2}, []int{0, 9})
		want := []int{2, 9}
		checkSums(t, got, want)
	})

	t.Run("empty slice", func(t *testing.T) {
		got := SumAllTails([]int{}, []int{3, 4, 5})
		want := []int{0, 9}
		checkSums(t, got, want)
	})
}

func TestSquare(t *testing.T) {
	a := [...]int{1, 2, 3}
	square(&a)
	fmt.Println(a)
	if a[1] != 4 && a[2] != 9 {
		t.Errorf("failed, a[1]=%d,a[2]=%d", a[1], a[2])
	}
}

func printLenCap(nums []int) {
	fmt.Printf("%v: len=%d, cap=%d\n", nums, cap(nums), len(nums))
}

// === RUN   TestSliceLenAndCap
// [1]: len=1, cap=1
// [1 2]: len=2, cap=2
// [1 2 3]: len=4, cap=3
// [1 2 3 4]: len=4, cap=4
// [1 2 3 4 5]: len=8, cap=5
// --- PASS: TestSliceLenAndCap (0.00s)
// 往切片中不断地增加新的元素。如果超过了当前切片的容量，就需要分配新的内存，并将当前切片所有的元素拷贝到新的内存块上
// 容量在比较小的时候，一般是以 2 的倍数扩大的，例如 2 4 8 16 …，当达到 2048 时，会采取新的策略
func TestSliceLenAndCap(t *testing.T) {
	nums := []int{1}
	printLenCap(nums)
	nums = append(nums, 2)
	printLenCap(nums)
	nums = append(nums, 3)
	printLenCap(nums)
	nums = append(nums, 4)
	printLenCap(nums)
	nums = append(nums, 5)
	printLenCap(nums)
}

// https://geektutu.com/post/hpg-slice.html
// === RUN   TestSliceAppend
// [1 2 3 4 5]: len=8, cap=5
// [3 4]: len=6, cap=2
// [1 2 3 4 50]: len=8, cap=5
// [3 4 50 60]: len=6, cap=4
// --- PASS: TestSliceAppend (0.00s)
// PASS
func TestSliceAppend(t *testing.T) {
	nums := make([]int, 0, 8)
	nums = append(nums, 1, 2, 3, 4, 5)
	nums2 := nums[2:4]
	printLenCap(nums)  // len: 5, cap: 8 [1 2 3 4 5]
	printLenCap(nums2) // len: 2, cap: 6 [3 4]

	// nums2 增加 2 个元素 50 和 60 后，将底层数组下标 [4] 的值改为了 50，下标[5] 的值置为 60
	// nums2[2]=50 会覆盖 nums[4]=5
	nums2 = append(nums2, 50, 60)
	printLenCap(nums)  // len: 5, cap: 8 [1 2 3 4 50]
	printLenCap(nums2) // len: 4, cap: 6 [3 4 50 60]
}

// 在已有切片的基础上进行切片，不会创建新的底层数组。因为原来的底层数组没有发生变化，内存会一直占用，直到没有变量引用该数组。
// 因此很可能出现这么一种情况，原切片由大量的元素构成，但是我们在原切片的基础上切片，虽然只使用了很小一段，
// 但底层数组在内存中仍然占据了大量空间，得不到释放。比较推荐的做法，使用 copy 替代 re-slice
func lastNumsBySlice(origin []int) []int {
	return origin[len(origin)-2:]
}

func lastNumsByCopy(origin []int) []int {
	result := make([]int, 2)
	copy(result, origin[len(origin)-2:])
	return result
}

func generateWithCap(n int) []int {
	rand.Seed(time.Now().UnixNano())
	nums := make([]int, 0, n)
	for i := 0; i < n; i++ {
		nums = append(nums, rand.Int())
	}
	return nums
}

func printMem(t *testing.T) {
	t.Helper()
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)
	t.Logf("%.2f MB", float64(rtm.Alloc)/1024./1024.)
}

// generateWithCap 用于随机生成 n 个 int 整数，64位机器上，一个 int 占 8 Byte，128 * 1024 个整数恰好占据 1 MB 的空间
func testLastChars(t *testing.T, f func([]int) []int) {
	t.Helper()
	ans := make([][]int, 0)
	for k := 0; k < 100; k++ {
		origin := generateWithCap(128 * 1024) // 1M
		ans = append(ans, f(origin))
		//如果我们在循环中，显示地调用 runtime.GC()，效果会更加地明显
		//=== RUN   TestLastCharsByCopy
		//    sum_test.go:189: 0.25 MB
		//--- PASS: TestLastCharsByCopy (0.21s)
		//runtime.GC()
	}
	printMem(t)
	_ = ans
}

// go test -run=^TestLastChars  -v
// === RUN   TestLastCharsBySlice
//
//	sum_test.go:168: 100.25 MB
//
// --- PASS: TestLastCharsBySlice (0.27s)
// === RUN   TestLastCharsByCopy
//
//	sum_test.go:169: 1.25 MB
//
// --- PASS: TestLastCharsByCopy (0.21s)
// astNumsBySlice 耗费了 100.25 MB 内存，申请的 100 个 1 MB 大小的内存没有被回收。因为切片虽然只使用了最后 2 个元素，
// 但是因为与原来 1M 的切片引用了相同的底层数组，底层数组得不到释放
func TestLastCharsBySlice(t *testing.T) { testLastChars(t, lastNumsBySlice) }

// lastNumsByCopy 仅消耗了 1.25 MB 的内存。这是因为，通过 copy，指向了一个新的底层数组，当 origin 不再被引用后，内存会被垃圾回收
func TestLastCharsByCopy(t *testing.T) { testLastChars(t, lastNumsByCopy) }

func foo(a []int) {
	//新切片 a 增加了 8 个元素，原切片对应的底层数组不够放置这 8 个元素，
	//因此申请了新的空间来放置扩充后的底层数组。这个时候新切片和原切片指向的底层数组就不是同一个了
	a = append(a, 1, 2, 3, 4, 5, 6, 7, 8)
	a[0] = 200
}

func foo1(a *[]int) {
	*a = append(*a, 1, 2, 3, 4, 5, 6, 7, 8)
	(*a)[0] = 200
}

// === RUN   TestSliceFoo
// [1 2]
// [200 2 1 2 3 4 5 6 7 8]
func TestSliceFoo(t *testing.T) {
	a := []int{1, 2}
	foo(a)
	fmt.Println(a)
	foo1(&a)
	fmt.Println(a)
}

// 设置返回值，将新切片返回并赋值给 main 函数中的变量 a
func bar(a []int) []int {
	a = append(a, 1, 2, 3, 4, 5, 6, 7, 8)
	a[0] = 200
	return a
}

// === RUN   TestSliceBar
// [200 2 1 2 3 4 5 6 7 8]
func TestSliceBar(t *testing.T) {
	a := []int{1, 2}
	//传参时拷贝了新的切片
	a = bar(a)
	fmt.Println(a)
}

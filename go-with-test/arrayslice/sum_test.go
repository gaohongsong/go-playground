package arrayslice

import (
	"fmt"
	"reflect"
	"testing"
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

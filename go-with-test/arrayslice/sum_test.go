package arrayslice

import (
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

package forrange

import (
	"fmt"
)

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

func ExampleRangeMap() {
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

package arrayslice

import "fmt"

// 数组
// 切片
// 多种方式的切片初始化
// 切片的容量是 固定 的，但是你可以使用 append 从原来的切片中创建一个新切片
// 如何获取部分切片
// 使用 len 获取数组和切片的长度
// 使用测试代码覆盖率的工具
// reflect.DeepEqual 的妙用和对代码类型安全性的影响
// https://go.dev/blog/slices-intro

func Sum(nums [5]int) int {
	sum := 0
	for _, num := range nums {
		sum += num
	}
	return sum
}

func SumSlice(nums []int) int {
	sum := 0
	for _, num := range nums {
		sum += num
	}
	return sum
}

// SumAll 可变参数 paramName ...type
func SumAll(numsToSum ...[]int) (sums []int) {
	// make 可以在创建切片的时候指定我们需要的长度和容量
	// sums = make([]int, len(numsToSum))

	for i, nums := range numsToSum {
		fmt.Println(i)
		// append 能为切片追加一个新值
		sums = append(sums, SumSlice(nums))
		// sums[i] = SumSlice(nums)
	}

	return
}

// SumAllTails 把每个切片的尾部元素相加（尾部的意思就是除去第一个元素以外的其他元素）
func SumAllTails(numbersToSum ...[]int) []int {
	var sums []int
	for _, numbers := range numbersToSum {
		sumTail := 0
		if len(numbers) > 0 {
			sumTail = SumSlice(numbers[1:])
		}
		sums = append(sums, sumTail)
	}
	return sums
}

// 为了避免复制数组，一般会传递指向数组的指针
func square(arr *[3]int) {
	for i, num := range *arr {
		(*arr)[i] = num * num
	}
}

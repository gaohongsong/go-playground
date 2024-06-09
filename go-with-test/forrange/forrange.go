package forrange

import (
	"math/rand"
	"time"
)

// 直接上结论
// range 在迭代过程中返回的是迭代值的拷贝，如果每次迭代的元素的内存占用很低，那么 for 和 range 的性能几乎是一样，例如 []int。
// 但是如果迭代的元素内存占用较高，例如一个包含很多属性的 struct 结构体，那么 for 的性能将显著地高于 range，
// 有时候甚至会有上千倍的性能差异。对于这种场景，建议使用 for，如果使用 range，建议只迭代下标，
// 通过下标访问迭代值，这种使用方式和 for 就没有区别了。
// 如果想使用 range 同时迭代下标和值，则需要将切片/数组的元素改为指针，才能不影响性能
func generateWithCap(cap int) []int {
	// 空切片，供追加元素
	nums := make([]int, 0, cap)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < cap; i++ {
		nums = append(nums, r.Int())
	}
	return nums
}

// Item 实例需要申请约 4KB 的内存
// 与 for 不同的是，range 对每个迭代值都创建了一个拷贝。因此如果每次迭代的值内存占用很小的情况下，
// for 和 range 的性能几乎没有差异，但是如果每个迭代值内存占用很大，
// 例如下面的例子中，每个结构体需要占据 4KB 的内存，这种情况下差距就非常明显了
type Item struct {
	id  int
	val [4096]byte
}

func generateItems(cap int) []*Item {
	items := make([]*Item, 0, cap)
	for i := 0; i < cap; i++ {
		items = append(items, &Item{id: i})
	}
	return items
}

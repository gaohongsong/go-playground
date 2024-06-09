package memalign

import (
	"fmt"
	"testing"
	"unsafe"
)

// https://geektutu.com/post/hpg-struct-alignment.html
// 合理的内存对齐可以提高内存读写的性能，并且便于实现变量操作的原子性。
// 计算结构体占用的空间 unsafe.Sizeof
// 一个结构体实例所占据的空间等于各字段占据空间之和，再加上内存对齐的空间大小
type Args struct {
	// x64 -> 8byte -> 16byte (Args)
	num1 int
	num2 int
}

type Flag struct {
	// 2byte
	num1 int16
	// 4byte
	num2 int32
}

// === RUN   TestStructSize
// 16
// 8
func TestStructSize(t *testing.T) {
	// 16
	fmt.Println(unsafe.Sizeof(Args{}))
	// 8, 为什么不是 2+4，多出来的2byte是内存对齐的结果
	fmt.Println(unsafe.Sizeof(Flag{}))

	//对于任意类型的变量 x ，unsafe.Alignof(x) 至少为 1
	//对于 array 数组类型的变量 x，unsafe.Alignof(x) 等于构成数组的元素类型的对齐倍数
	//对于 struct 结构体类型的变量 x，计算 x 每一个字段 f 的 unsafe.Alignof(x.f)，
	//unsafe.Alignof(x) 等于其中的最大值。

	//Args{} 的对齐倍数是 8，Args{} 两个字段占据 16 字节，是 8 的倍数，无需占据额外的空间对齐
	fmt.Println(unsafe.Alignof(Args{}))
	//Flag{} 的对齐倍数是 4，因此 Flag{} 占据的空间必须是 4 的倍数，因此，6 内存对齐后是 8 字节
	fmt.Println(unsafe.Alignof(Flag{}))
}

// 合理布局减少内存占用
// 顺序会对 struct 的大小产生影响，在对内存特别敏感的结构体的设计上，
// 我们可以通过调整字段的顺序，减少内存的占用
// a a b b c c c c -> 8 VS 4 倍数 -> 8
type demo1 struct {
	a int8  // 1
	b int16 // 2
	c int32 // 4
}

// a a a a c c c c b b -> 10 vs 4（结构体取字段最大对齐值） 倍数 -> 12
type demo2 struct {
	a int8  // 1
	c int32 // 4
	b int16 // 2
}

func TestAlignOrder(t *testing.T) {
	fmt.Println(unsafe.Sizeof(demo1{})) // 8
	fmt.Println(unsafe.Sizeof(demo2{})) // 12
}

//空 struct{} 的对齐
//struct{} 作为其他 struct 最后一个字段时，需要填充额外的内存保证安全

// 8 ？？？
type demo3 struct {
	// 4
	c int32
	// 当 struct{} 作为结构体最后一个字段时，需要内存对齐
	//因为如果有指针指向该字段, 返回的地址将在结构体之外，
	//如果此指针一直存活不释放对应的内存，就会有内存泄露的问题（该内存不因结构体释放而释放）
	// 4 额外填充了 4 字节的空间
	a struct{} // 0
}

// 4
type demo4 struct {
	// 0
	a struct{} // 空 struct{} 大小为 0，作为其他 struct 的字段时，一般不需要内存对齐
	// 4
	c int32
}

func TestAlignOrder2(t *testing.T) {
	fmt.Println(unsafe.Sizeof(demo3{})) // 8
	fmt.Println(unsafe.Sizeof(demo4{})) // 4
}

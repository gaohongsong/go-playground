package main

import (
	"math/rand"
	"os"
	"runtime/pprof"
	"time"
)

func bubbleSort(nums []int) {
	for i := 0; i < len(nums)-1; i++ {
		for j := i + 1; j < len(nums); j++ {
			if nums[i] > nums[j] {
				nums[i], nums[j] = nums[j], nums[i]
			}
		}
	}
}

func generate(n int) []int {
	nums := make([]int, 0)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < n; i++ {
		nums = append(nums, r.Intn(100))
	}
	return nums
}

// web查看
// go tool pprof -http=:9999 cpu.pprof
// 命令行查看
// go tool pprof cpu.pprof
// Duration: 3.74s, Total samples = 3.50s (93.54%)
// Entering interactive mode (type "help" for commands, "o" for options)
// (pprof) top
// Showing nodes accounting for 3.49s, 99.71% of 3.50s total
// Dropped 22 nodes (cum <= 0.02s)
//
//	 flat  flat%   sum%        cum   cum%
//	3.46s 98.86% 98.86%      3.49s 99.71%  main.bubbleSort (inline)
//	0.03s  0.86% 99.71%      0.03s  0.86%  runtime.asyncPreempt
//	    0     0% 99.71%      3.49s 99.71%  main.main
//	    0     0% 99.71%      3.49s 99.71%  runtime.main
//
// (pprof) gif
func main() {
	f, _ := os.OpenFile("cpu.pprof", os.O_CREATE|os.O_RDWR, 0644)
	defer f.Close()
	pprof.StartCPUProfile(f)

	//pprof.StartCPUProfile(os.Stdout)
	defer pprof.StopCPUProfile()

	n := 10
	for i := 0; i < 5; i++ {
		nums := generate(n)
		bubbleSort(nums)
		n *= 10
	}
}

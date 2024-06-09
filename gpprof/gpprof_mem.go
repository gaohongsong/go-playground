package main

import (
	"github.com/pkg/profile"
	"math/rand"
	"strings"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// 字符串是不可变的，因为将两个字符串拼接时，相当于是产生新的字符串，
// 如果当前的空间不足以容纳新的字符串，则会申请更大的空间，
// 将新字符串完全拷贝过去，这消耗了 2 倍的内存空间。
// 在这 100 次拼接的过程中，会产生多次字符串拷贝，从而消耗大量的内存
// (pprof) top --cum
// Showing nodes accounting for 524.73kB, 91.28% of 574.87kB total
// Dropped 63 nodes (cum <= 2.87kB)
// Showing top 10 nodes out of 26
//
//	    flat  flat%   sum%        cum   cum%
//	       0     0%     0%   573.59kB 99.78%  runtime.main
//	       0     0%     0%   573.48kB 99.76%  main.main
//	524.61kB 91.26% 91.26%   546.63kB 95.09%  main.concat
//	  0.09kB 0.016% 91.27%    26.59kB  4.63%  github.com/pkg/profile.Start
//	       0     0% 91.27%    26.50kB  4.61%  github.com/pkg/profile.Start.func2
//	       0     0% 91.27%    26.50kB  4.61%  log.(*Logger).output
//	       0     0% 91.27%    26.50kB  4.61%  log.Printf (inline)
//	       0     0% 91.27%    23.64kB  4.11%  sync.(*Once).Do (inline)
//	       0     0% 91.27%    23.64kB  4.11%  sync.(*Once).doSlow
//	  0.03kB 0.0054% 91.28%    23.60kB  4.11%  log.formatHeader
func concat(n int) string {
	s := ""
	for i := 0; i < n; i++ {
		s += randomString(n)
	}
	return s
}

// (pprof) top --cum
// Showing nodes accounting for 44.36kB, 46.98% of 94.41kB total
// Dropped 42 nodes (cum <= 0.47kB)
// Showing top 10 nodes out of 39
//
//	   flat  flat%   sum%        cum   cum%
//	      0     0%     0%    93.13kB 98.64%  main.main
//	      0     0%     0%    93.13kB 98.64%  runtime.main
//	      0     0%     0%    66.29kB 70.21%  main.concatGood
//	44.27kB 46.88% 46.88%    44.27kB 46.88%  strings.(*Builder).WriteString (inline)
//	 0.09kB 0.099% 46.98%    26.59kB 28.17%  github.com/pkg/profile.Start
//	      0     0% 46.98%    26.50kB 28.07%  github.com/pkg/profile.Start.func2
//	      0     0% 46.98%    26.50kB 28.07%  log.(*Logger).output
//	      0     0% 46.98%    26.50kB 28.07%  log.Printf (inline)
//	      0     0% 46.98%    23.64kB 25.04%  sync.(*Once).Do (inline)
//	      0     0% 46.98%    23.64kB 25.04%  sync.(*Once).doSlow
func concatGood(n int) string {
	var s strings.Builder
	for i := 0; i < n; i++ {
		s.WriteString(randomString(n))
	}
	return s.String()
}

// go run .\gpprof_mem.go
// go tool pprof -http=:9999 /tmp/profile215959616/mem.pprof
func main() {
	//pkg/profile 封装了 runtime/pprof 的接口，使用起来更简单
	defer profile.Start(profile.MemProfile, profile.MemProfileRate(1)).Stop()
	concat(100)
	//concatGood(100)
}

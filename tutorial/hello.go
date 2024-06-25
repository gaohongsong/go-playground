package main

import (
	"fmt"
	"strings"
	"time"
)

func main() {
	fmt.Println("hello world")
	v1 := []string{}
	fmt.Println(strings.Join(v1, ","))
	fmt.Println("aaa")

	now := time.Now()
	hourAgo := now.Add(-time.Hour)
	fmt.Println(hourAgo)
}

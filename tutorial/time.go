package main

import (
	"fmt"
	"time"
)

func main() {
	end := time.Now()
	start := end.AddDate(0, 0, -1)

	endTimestamp := end.Unix()
	startTimestamp := start.Unix()

	fmt.Println("Start Timestamp:", startTimestamp)
	fmt.Println("End Timestamp:", endTimestamp)
}

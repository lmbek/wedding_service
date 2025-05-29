package main

import (
	"fmt"
	"time"
)

func main() {
	i := 0
	for {
		i++
		fmt.Printf("Hej %v\n", i)
		time.Sleep(1 * time.Second)
	}
}

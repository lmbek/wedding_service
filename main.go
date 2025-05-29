package main

import (
	"fmt"
	"time"
)

var Running bool

func Start() {
	Running = true
	fmt.Println("Started")
}

func Stop() {
	Running = false
	fmt.Println("Stopped")
}

func Run() {
	for Running {
		fmt.Println("Running...")
		time.Sleep(100 * time.Millisecond) // shorter delay for test speed
	}
}

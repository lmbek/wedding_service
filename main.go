package main

import (
	"fmt"
	"time"
)

var Running bool

func main() {
	Start()

	go Run()

	// Let it run for 1 second
	time.Sleep(2 * time.Second)

	Stop()
}

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

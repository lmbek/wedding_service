package main

import (
	"testing"
	"time"
)

func Test_main(t *testing.T) {
	go func() {
		time.Sleep(2 * time.Second)
		httpServer.Close()
		httpsServer.Close()
	}()

	main()
}

func Test_start(t *testing.T) {
	go func() {
		time.Sleep(2 * time.Second)
		httpServer.Close()
		httpsServer.Close()
	}()

	err := start()
	if err != nil {
		t.Errorf("err: %v", err)
	}
}

package main

import (
	"testing"
	"time"
)

func TestRunLifecycle(t *testing.T) {
	Start()

	done := make(chan struct{})
	go func() {
		Run()
		close(done)
	}()

	time.Sleep(300 * time.Millisecond) // Let it "run" a bit

	if !Running {
		t.Error("Expected Running to be true")
	}

	Stop()

	select {
	case <-done:
		// success
	case <-time.After(1 * time.Second):
		t.Fatal("Run() did not exit after Stop()")
	}

	if Running {
		t.Error("Expected Running to be false after Stop()")
	}
}

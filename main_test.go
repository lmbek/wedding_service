package main

import (
	"errors"
	"net/http"
	"testing"
	"time"
)

func Test_main(t *testing.T) {
	// TODO - use synctest
	go func() {
		time.Sleep(2 * time.Second)
		httpServer.Close()
		httpsServer.Close()
	}()

	main()
}

func Test_start(t *testing.T) {
	// TODO - use synctest
	go func() {
		time.Sleep(2 * time.Second)
		httpServer.Close()
		httpsServer.Close()
	}()

	err := start()
	if err != nil && !errors.Is(err, http.ErrServerClosed) { // if the server is closed, we probably forced it
		t.Errorf("err: %v", err)
	}
}

package main

import (
	"fmt"
	"net/http"
	"testing"
	"time"
	"wedding_service/env"
	"wedding_service/webserver"
)

func Example_main() {
	err := env.Init()
	if err != nil {
		fmt.Printf("could not env.Init(): %v", err)
		return
	}

	w, err := webserver.NewWebserver()
	if err != nil {
		fmt.Printf("could not webserver.NewWebserver(): %v", err)
		return
	}

	err = w.ListenAndServe()
	if err != nil {
		fmt.Printf("%s", fmt.Errorf("could not w.ListenAndServe(): %w", err))
		return
	}
}

func Test_main(t *testing.T) {
	err := env.Init()
	if err != nil {
		fmt.Printf("could not env.Init(): %v", err)
		return
	}

	go func() {
		couldRequest := waitForWebserverGetResponse(t, 0)
		if !couldRequest {
			t.Errorf("could not request webserver withing allowed time")
		}

		if mainWebserver != nil {
			mainWebserver.Close()
		}
	}()

	main()
	mainWebserver = nil
}

func Test_initEnv(t *testing.T) {
	initEnv()
	defer env.Reset()

	t.Run("run invalid env file", func(t *testing.T) {
		tempDir := t.TempDir()
		t.Chdir(tempDir)
		err := initEnv()
		if err == nil {
			t.Errorf("err should not be nil: %s", err)
			return
		}
	})
}

func Test_createWebserver(t *testing.T) {
	env.Init()
	defer env.Reset()

	_, err := createWebserver()
	if err != nil {
		t.Errorf("could not create webserver: %s", err)
		return
	}

	t.Run("invalid webserver configuration", func(t *testing.T) {
		defer env.Reset()
		env.Env.CertPath = "invalid path"
		env.Env.KeyPath = "invalid path"

		_, err := createWebserver()
		if err == nil {
			t.Errorf("err should not be nil")
			return
		}
	})
}

func Test_startWebserver(t *testing.T) {
	w, err := webserver.NewWebserver()
	if err != nil {
		t.Errorf("could not create new webserver: %s", err)
		return
	}

	go func() {
		couldRequest := waitForWebserverGetResponse(t, 0)
		if !couldRequest {
			t.Errorf("could not request webserver withing allowed time")
		}

		if w != nil {
			w.Close()
		}
	}()

	err = startWebserver(w)
	if err != nil {
		t.Errorf("could not start webserver: %s", err)
		return
	}

	t.Run("invalid configuration", func(t *testing.T) {
		defer env.Reset()
		env.Env.HttpPort = "-1"
		env.Env.HttpsPort = "-1"

		w, err := webserver.NewWebserver()
		if err != nil {
			t.Errorf("could not create new webserver: %s", err)
			return
		}

		go func() {
			couldRequest := waitForWebserverGetResponse(t, 0)
			if !couldRequest {
				t.Errorf("could not request webserver withing allowed time")
			}

			if w != nil {
				w.Close()
			}
		}()

		err = startWebserver(w)
		if err == nil {
			t.Errorf("should not be nil")
			return
		}
	})
}

func waitForWebserverGetResponse(t *testing.T, retryNr int) bool {
	client := &http.Client{}

	addr := "http://localhost:" + env.Env.HttpPort
	resp, err := client.Get(addr)
	if err != nil {
		// retry for up to 5 seconds
		if retryNr < 5000 {
			time.Sleep(1 * time.Millisecond)
			retryNr++
			return waitForWebserverGetResponse(t, retryNr) // <-- FIX: return the recursive call result
		}
		t.Errorf("Request failed after %d milliseconds full of retries: %s", retryNr, err)
		return false
	}
	defer resp.Body.Close()
	return true
}

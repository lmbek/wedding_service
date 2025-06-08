package main

import (
	"fmt"
	"testing"
	"testing/synctest"
	"time"
	"wedding_service/env"
	"wedding_service/webserver"
)

func Example_main() {
	env.Init()

	w, err := webserver.NewWebserver()
	if err != nil {
		fmt.Printf("could not make new webserver: %v", err)
		return
	}

	err = w.ListenAndServe()
	if err != nil {
		fmt.Printf("%s", fmt.Errorf("could not ListenAndServe: %w", err))
		return
	}
}

func Test_main(t *testing.T) {
	env.Init()

	// we use synctest to simulate the time.Sleep, but without waiting 500ms
	synctest.Run(func() {
		// For main func testing we need to initialize a specific main webserver to test on
		w, err := webserver.NewWebserver()
		if err != nil {
			t.Errorf("could not make new webserver: %v", err)
			return
		}
		// global mainWebserver gets set, so that we can test on it
		mainWebserver = w

		// close before we start to close the server instantly (use mock testing to test endpoints instead)
		err = mainWebserver.Close()
		if err != nil {
			t.Errorf("could not close webserver: %v", err)
			return
		}

		main()
		resetMainWebserver()

		// this is for having the prints from the test be in the correct place (verbose testing)
		time.Sleep(500 * time.Millisecond)
	})

}

func Test_start(t *testing.T) {
	env.Init()
	// start func tested by running main func test

	t.Run("no webserver test", func(t *testing.T) {
		err := start(nil)
		if err == nil {
			t.Errorf("expected an error but got nil")
		}
	})

	t.Run("invalid env for webserver", func(t *testing.T) {
		// change env "mode" to test it with invalid mode (no certificates to be handled)
		defer env.Reset()
		env.Env.HttpPort = "-1"
		env.Env.HttpsPort = "-1"

		w, err := webserver.NewWebserver()
		if err != nil {
			t.Errorf("could not create new webserver: %v", err)
			return
		}

		err = start(w)
		if err == nil {
			t.Errorf("should not be able to start webserver: %v", err)
			return
		}
	})
}

func BenchmarkStart(b *testing.B) {
	env.Init()
	synctest.Run(func() {
		// use no info debugging
		defer env.Reset()
		env.Env.DebugLevel = 2

		for i := 0; i < b.N; i++ {
			w, err := webserver.NewWebserver()
			if err != nil {
				b.Errorf("could not create new webserver: %v", err)
				return
			}

			err = w.Close()
			if err != nil {
				b.Errorf("could not close webserver: %v", err)
				return
			}

			err = start(w)
			if err != nil {
				b.Errorf("could not start webserver: %v", err)
				return
			}
		}
		time.Sleep(500 * time.Millisecond)
	})
}

func Test_initMainWebserver(t *testing.T) {
	env.Init()
	m, err := initMainWebserver()
	if err != nil {
		t.Errorf("could not init webserver: %v", err)
		return
	}

	if m == nil {
		t.Errorf("mainWebserver should not be nil: %v", err)
		return
	}

	if mainWebserver == nil {
		t.Errorf("mainWebserver should not be nil: %v", err)
		return
	}

	// reset global variables - this is needed before testing mainWebserver again
	resetMainWebserver()

	// test fail scenario:
	t.Run("fail scenario", func(t *testing.T) {
		// change env "mode" to test it with invalid mode (no certificates to be handled)
		defer env.Reset()
		env.Env.Mode = ""

		_, err := initMainWebserver()
		if err == nil {
			t.Errorf("initMainWebserver should have failed: %v", err)
			return
		}
		// reset global variables - this is needed before testing mainWebserver again
		resetMainWebserver()
	})
}

func resetMainWebserver() {
	mainWebserver = nil
}

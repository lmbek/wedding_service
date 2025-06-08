package main

import (
	"errors"
	"fmt"
	"log"
	"wedding_service/env"
	"wedding_service/webserver"
)

var mainWebserver webserver.Webserver

// main starts the program (ignore errors as they are already handled by the other functions in main package
func main() {
	env.Init()
	m, _ := initMainWebserver()
	_ = start(m)
}

// start handles the eventual error internally and prints it (main should ignore the error)
// testWebserver is used only for the tests, for normal use we should just give nil
func start(w webserver.Webserver) (err error) {
	if w == nil {
		return errors.New("webserver is nil")
	}

	err = w.ListenAndServe()
	if err != nil {
		log.Printf("%s", fmt.Errorf("could not ListenAndServe: %w", err))
		return err
	}
	return nil
}

func initMainWebserver() (webserver.Webserver, error) {
	if mainWebserver != nil {
		return mainWebserver, nil
	}

	m, err := webserver.NewWebserver()
	if err != nil {
		if env.IsDebugErrorsEnabled() {
			log.Printf("%s", fmt.Errorf("could not create new webserver: %w", err))
		}
		return nil, err
	}
	mainWebserver = m
	return mainWebserver, nil
}

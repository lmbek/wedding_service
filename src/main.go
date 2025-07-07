package main

import (
	"flag"
	"fmt"
	"log"
	"wedding_service/env"
	"wedding_service/webserver"
)

// TODO: rework this main.go into a much smaller version and reduce amount of testing needed. Simplify things
var mainWebserver webserver.Webserver

func main() {
	flag.Parse()
	err := initEnv()
	if err != nil {
		fmt.Printf("err initEnv(): %v\n", err)
		return
	}
	mainWebserver, _ = createWebserver()
	startWebserver(mainWebserver)
}

func initEnv() error {
	err := env.Init()
	if err != nil {
		log.Printf("could not initEnv in main: %s\n", err)
		return err
	}
	return nil
}

func createWebserver() (webserver.Webserver, error) {
	w, err := webserver.NewWebserver()
	if err != nil {
		log.Printf("%s", err)
		return nil, err
	}

	return w, nil
}

func startWebserver(w webserver.Webserver) error {
	err := w.ListenAndServe()
	if err != nil {
		log.Printf("%s", err)
		return err
	}

	return nil
}

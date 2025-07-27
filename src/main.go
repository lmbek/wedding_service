package main

import (
	"log"
	"wedding_service/config"
	"wedding_service/webserver"
	"wedding_service/webserver/website/frontend"
)

func main() {
	var err error
	var c config.Config
	var f frontend.Frontend
	var w webserver.Webserver

	// Initialize the environment
	c, err = config.NewConfig()
	if err != nil {
		log.Fatalf("Error initializing environment: %v", err)
		return
	}

	// Add a frontend
	f, err = frontend.NewFrontend(c.FrontendPath(), c.HotReloadEnabled())
	if err != nil {
		log.Fatalf("Error initializing frontend: %v", err)
		return
	}

	// Add a webserver with the frontend
	w, err = webserver.NewWebserver(c, f)
	if err != nil {
		log.Printf("Error creating webserver: %v", err)
		return
	}

	// Listen to requests (blocking call)
	err = w.ListenAndServe()
	if err != nil {
		log.Printf("Error starting listening on webserver: %v", err)
		return
	}
	defer w.Close()
}

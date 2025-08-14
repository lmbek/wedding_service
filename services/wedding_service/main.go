package main

import (
	"log"
	"wedding_service/config"
	"wedding_service/webserver"
	"wedding_service/webserver/website/frontend"
)

func main() {

	c, err := config.NewConfig()
	if err != nil {
		log.Println(err)
	}

	// Initialize the frontend
	f, err := frontend.NewFrontend(c.FrontendPath(), c.HotReloadEnabled())
	if err != nil {
		log.Println(err)
	}

	// Initialize the webserver
	w, err := webserver.NewWebserver(c, f)
	if err != nil {
		log.Println(err)
	}

	err = w.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}

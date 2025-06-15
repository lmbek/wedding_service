package webserver

import (
	"fmt"
	"net/http"
	"wedding_service/env"
)

func newHttpServer(port string) *http.Server {
	var addr string

	if env.IsModeDevelopment() {
		addr = fmt.Sprintf("localhost:%s", port)
	} else {
		addr = fmt.Sprintf(":%s", port)
	}

	return &http.Server{
		Addr: addr,
	}
}

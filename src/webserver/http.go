package webserver

import (
	"fmt"
	"net/http"
	"wedding_service/buildtag"
)

func newHttpServer(port string) *http.Server {
	var addr string
	addr = fmt.Sprintf(":%s", port)

	if buildtag.IsDevelopment() {
		addr = fmt.Sprintf("localhost:%s", port)
	}

	return &http.Server{
		Addr: addr,
	}
}

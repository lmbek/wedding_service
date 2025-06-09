package webserver

import (
	"fmt"
	"net/http"
)

func newHttpServer(port string) *http.Server {
	return &http.Server{
		Addr: fmt.Sprintf(":%s", port),
	}
}

package webserver

import (
	"fmt"
	"net/http"
)

func newHttpServer(port string) *http.Server {
	var addr string
	addr = fmt.Sprintf(":%s", port)

	return &http.Server{
		Addr: addr,
	}
}

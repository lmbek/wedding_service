package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"wedding_service/webserver"
)

var httpServer = &http.Server{Addr: fmt.Sprintf("localhost:%s", os.Getenv("WEDDING_SERVICE_HTTP_PORT"))}
var httpsServer = &http.Server{Addr: fmt.Sprintf("localhost:%s", os.Getenv("WEDDING_SERVICE_HTTPS_PORT"))}

func main() {
	start()
}

func start() error {
	// TODO - less parameters
	return webserver.Start(httpServer, httpsServer, filepath.Join("certificate", os.Getenv("LOCALHOST_CERT")), filepath.Join("certificate", os.Getenv("LOCALHOST_KEY")))
}

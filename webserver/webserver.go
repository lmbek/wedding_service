package webserver

import (
	"errors"
	"fmt"
	"net/http"
)

func Start(httpServer *http.Server, httpsServer *http.Server, certPath string, keyPath string) error {
	err := setupHttpsServer(httpsServer, certPath, keyPath)
	if err != nil {
		return err
	}

	return ListenAndServe(httpServer, httpsServer)
}

func ListenAndServe(httpServer *http.Server, httpsServer *http.Server) error {
	errChannel := make(chan error, 2)

	go func() {
		fmt.Println("Listening on http://" + httpServer.Addr)
		err := httpServer.ListenAndServe()
		errChannel <- err
	}()

	go func() {
		fmt.Println("Listening on https://" + httpsServer.Addr)
		err := httpsServer.ListenAndServeTLS("", "")
		errChannel <- err
	}()

	for i := 0; i < 2; i++ {
		err := <-errChannel
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}

	return nil
}

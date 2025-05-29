package webserver

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
)

func Start() error {

	return errors.New("Not implemented")
}

func ListenAndServe(httpServer *http.Server, httpsServer *http.Server) error {
	var wg sync.WaitGroup
	wg.Add(2)

	var httpErr, httpsErr error

	go func() {
		defer wg.Done()
		fmt.Println("Listening on http://" + httpServer.Addr)
		err := httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("HTTP server error: %s\n", err)
			httpErr = err
		}
	}()

	go func() {
		defer wg.Done()
		fmt.Println("Listening on https://" + httpsServer.Addr)
		err := httpsServer.ListenAndServeTLS("", "")
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("HTTPS server error: %s\n", err)
			httpsErr = err
		}
	}()

	wg.Wait()

	if httpErr != nil {
		return httpErr
	}
	if httpsErr != nil {
		return httpsErr
	}
	return nil
}

package webserver

import (
	"net/http"
	"testing"
	"time"
)

func TestStart_CertError(t *testing.T) {
	httpServer := &http.Server{Addr: ":0"} // Use :0 to get random free port
	httpsServer := &http.Server{Addr: ":0"}

	// Intentionally bad cert/key paths to trigger error in setupHttpsServer
	err := Start(httpServer, httpsServer, "invalid-cert.pem", "invalid-key.pem")
	if err == nil {
		t.Errorf("Expected error due to invalid cert/key paths, got nil")
	}
}

func TestStart_Success(t *testing.T) {
	httpServer := &http.Server{Addr: ":0"}
	httpsServer := &http.Server{Addr: ":0"}

	// Use valid test cert and key files here (you need to provide these in your testdata folder)
	certPath := "../certificate/localhost_wedding_service.crt"
	keyPath := "../certificate/localhost_wedding_service.key"

	// Run Start in a goroutine because ListenAndServe blocks
	done := make(chan error, 1)
	go func() {
		err := Start(httpServer, httpsServer, certPath, keyPath)
		done <- err
	}()

	// Wait a little to let servers start
	time.Sleep(100 * time.Millisecond)

	// Shutdown servers to stop ListenAndServe blocking
	if err := httpServer.Close(); err != nil {
		t.Errorf("Failed to close httpServer: %v", err)
	}
	if err := httpsServer.Close(); err != nil {
		t.Errorf("Failed to close httpsServer: %v", err)
	}

	// Wait for Start to return after shutdown
	err := <-done
	if err != nil {
		t.Errorf("Start returned error: %v", err)
	}
}

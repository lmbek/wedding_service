package webserver

import (
	"net"
	"testing"
	"time"
	"wedding_service/env"
)

func TestNewWebserver(t *testing.T) {
	t.Chdir("..")
	env.Init()
	defer env.Reset()

	_, err := NewWebserver()
	if err != nil {
		t.Errorf("could not create new webserver: %s", err)
	}

	t.Run("test with no mode", func(t *testing.T) {
		defer env.Reset()
		// try to set https port to something wrong before running the NewWebserver again
		env.Env.Mode = ""
		_, err := NewWebserver()
		if err == nil {
			t.Errorf("err should not be nil")
		}
	})
}

func TestWebserver_ListenAndServe(t *testing.T) {
	t.Chdir("..")
	env.Init()
	defer env.Reset()

	ws, _ := NewWebserver()
	w, _ := ws.(*webserver)

	go func() {
		time.Sleep(2 * time.Second)
		w.Close()
	}()

	w.ListenAndServe()
}

func TestWebserver_listenHTTPS(t *testing.T) {
	env.Init()

	t.Run("test listenHTTPS", func(t *testing.T) {
		t.Chdir("..")
		defer env.Reset()

		ws, _ := NewWebserver()
		w, _ := ws.(*webserver)

		go func() {
			time.Sleep(2 * time.Second)
			w.Close()
		}()

		err := w.listenHTTPS()
		if err != nil {
			t.Errorf("could not listen on HTTPS: %s", err)
		}
	})

	t.Run("test httpsServer.ListenAndServeTLS error", func(t *testing.T) {
		t.Chdir("..")
		env.Env.HttpsPort = "8443"
		defer env.Reset()

		// occupy the server
		ln, _ := net.Listen("tcp", ":8443")
		defer ln.Close()

		ws, _ := NewWebserver()
		w, _ := ws.(*webserver)
		defer w.Close()

		// ListenAndServeTLS throws error because port is used
		err := w.listenHTTPS()
		if err == nil {
			t.Errorf("err should not be nil")
		}
	})
}

func TestWebserver_listenHTTP(t *testing.T) {
	env.Init()

	t.Run("test listenHTTP", func(t *testing.T) {
		t.Chdir("..")
		defer env.Reset()

		ws, _ := NewWebserver()
		w, _ := ws.(*webserver)

		go func() {
			time.Sleep(2 * time.Second)
			w.Close()
		}()

		err := w.listenHTTP()
		if err != nil {
			t.Errorf("could not listen on HTTP: %s", err)
		}
	})

	t.Run("test httpServer.ListenAndServe error", func(t *testing.T) {
		t.Chdir("..")
		env.Env.HttpPort = "8080"
		defer env.Reset()

		// occupy the server
		ln, _ := net.Listen("tcp", ":8080")
		defer ln.Close()

		ws, _ := NewWebserver()
		w, _ := ws.(*webserver)
		defer w.Close()

		// ListenAndServe throws error because port is used
		err := w.listenHTTP()
		if err == nil {
			t.Errorf("err should not be nil")
		}
	})
}

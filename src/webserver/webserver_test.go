package webserver

//func TestNewWebserver(t *testing.T) {
//	t.Chdir("..")
//	env.Init()
//	defer env.Reset()
//
//	_, err := NewWebserver(frontend.NewFrontend())
//	if err != nil {
//		t.Errorf("could not create new webserver: %s", err)
//	}
//
//	t.Run("invalid acme settings", func(t *testing.T) {
//		defer env.Reset()
//		env.Env.Hostnames = nil
//
//		_, err := NewWebserver(frontend.NewFrontend())
//		if err == nil {
//			t.Errorf("should give an error")
//		}
//	})
//
//	t.Run("invalid webserver settings", func(t *testing.T) {
//		defer env.Reset()
//		env.Env.CertPath = "invalid_cert"
//		env.Env.KeyPath = "invalid_key"
//
//		_, err := NewWebserver(frontend.NewFrontend())
//		if err == nil {
//			t.Errorf("should give an error")
//		}
//	})
//}
//
//func TestWebserver_ListenAndServe(t *testing.T) {
//	t.Chdir("..")
//	env.Init()
//	defer env.Reset()
//
//	w := createNewWebserver(t)
//
//	go func() {
//		defer w.Close()
//		couldRequest := waitForWebserverGetResponse(t, 0)
//		if !couldRequest {
//			t.Errorf("could not request webserver withing allowed time")
//		}
//	}()
//	err := w.ListenAndServe()
//	if err != nil {
//		t.Errorf("could not ListenAndServe: %s", err)
//	}
//}
//
//func TestWebserver_listenHTTPS(t *testing.T) {
//	t.Chdir("..")
//	env.Init()
//
//	t.Run("test listenHTTPS", func(t *testing.T) {
//		env.Env.HttpsPort = "8443"
//		defer env.Reset()
//
//		w := createNewWebserver(t)
//		w.Close()
//
//		err := w.listenHTTPS()
//		if err != nil {
//			t.Errorf("could not listen on HTTPS: %s", err)
//		}
//	})
//
//	t.Run("test httpsServer.ListenAndServeTLS error", func(t *testing.T) {
//		defer env.Reset()
//		env.Env.HttpsPort = "-1"
//
//		w := createNewWebserver(t)
//
//		// ListenAndServeTLS throws error because port is used
//		err := w.listenHTTPS()
//		if err == nil {
//			t.Errorf("err should not be nil")
//		}
//	})
//}
//
//func TestWebserver_listenHTTP(t *testing.T) {
//	t.Chdir("..")
//	env.Init()
//
//	t.Run("test listenHTTP", func(t *testing.T) {
//		env.Env.HttpPort = "8080"
//		defer env.Reset()
//
//		w := createNewWebserver(t)
//
//		w.Close()
//
//		err := w.listenHTTP()
//		if err != nil {
//			t.Errorf("could not listen on HTTP: %s", err)
//		}
//	})
//
//	t.Run("test httpServer.ListenAndServe error", func(t *testing.T) {
//		env.Env.HttpPort = "-1"
//		defer env.Reset()
//
//		w := createNewWebserver(t)
//
//		// ListenAndServe throws error because port is invalid
//		err := w.listenHTTP()
//		if err == nil {
//			t.Errorf("err should not be nil")
//		}
//	})
//}
//
//func createNewWebserver(t *testing.T) *webserver {
//	newWebserver, err := NewWebserver(frontend.NewFrontend())
//	if err != nil {
//		t.Errorf("could not create new webserver: %s", err)
//	}
//	w, couldCast := newWebserver.(*webserver)
//	if !couldCast {
//		t.Errorf("%s", err)
//	}
//
//	return w
//}
//
//func waitForWebserverGetResponse(t *testing.T, retryNr int) bool {
//	client := &http.Client{}
//
//	addr := "http://localhost:" + env.Env.HttpPort
//	resp, err := client.Get(addr)
//	if err != nil {
//		// retry for up to 5 seconds
//		if retryNr < 5000 {
//			time.Sleep(1 * time.Millisecond)
//			retryNr++
//			return waitForWebserverGetResponse(t, retryNr) // <-- FIX: return the recursive call result
//		}
//		t.Errorf("Request failed after %d milliseconds full of retries: %s", retryNr, err)
//		return false
//	}
//	defer resp.Body.Close()
//	return true
//}

package webserver

// TODO: develop the webserver tests

//
//func TestNewWebserver(t *testing.T) {
//	env.Init()
//	defer env.Reset()
//
//	t.Run("Development success", func(t *testing.T) {
//		// Mock UseLOCALHOST til success, hvis nødvendigt
//		// orig := certificate.UseLOCALHOST
//		// certificate.UseLOCALHOST = func() (*tls.Certificate, error) { return &tls.Certificate{}, nil }
//		// defer func() { certificate.UseLOCALHOST = orig }()
//
//		ws, err := NewWebserver()
//		if err != nil {
//			t.Fatalf("expected no error, got %v", err)
//		}
//		if ws == nil {
//			t.Fatal("expected non-nil webserver")
//		}
//	})
//
//	t.Run("Development UseLOCALHOST error", func(t *testing.T) {
//		// Mock UseLOCALHOST til fejl
//		// orig := certificate.UseLOCALHOST
//		// certificate.UseLOCALHOST = func() (*tls.Certificate, error) {
//		// 	return nil, errors.New("failed to load cert")
//		// }
//		// defer func() { certificate.UseLOCALHOST = orig }()
//
//		_, err := NewWebserver()
//		wantErr := "could not use localhost certificate: failed to load cert"
//		if err == nil || err.Error() != wantErr {
//			t.Fatalf("expected error %q, got %v", wantErr, err)
//		}
//	})
//
//	t.Run("Production success", func(t *testing.T) {
//		// Mock UseACME til success
//		// orig := certificate.UseACME
//		// certificate.UseACME = func() (*certificate.AutocertManager, error) {
//		// 	return &certificate.AutocertManager{}, nil
//		// }
//		// defer func() { certificate.UseACME = orig }()
//
//		ws, err := NewWebserver()
//		if err != nil {
//			t.Fatalf("expected no error, got %v", err)
//		}
//		if ws == nil {
//			t.Fatal("expected non-nil webserver")
//		}
//	})
//
//	t.Run("Production UseACME error", func(t *testing.T) {
//		// Mock UseACME til fejl
//		// orig := certificate.UseACME
//		// certificate.UseACME = func() (*certificate.AutocertManager, error) {
//		// 	return nil, errors.New("failed to get ACME manager")
//		// }
//		// defer func() { certificate.UseACME = orig }()
//
//		_, err := NewWebserver()
//		wantErr := "could not use acme manager: failed to get ACME manager"
//		if err == nil || err.Error() != wantErr {
//			t.Fatalf("expected error %q, got %v", wantErr, err)
//		}
//	})
//
//	t.Run("Mode not set", func(t *testing.T) {
//		// Simuler at MODE er tomt
//		env.Env.Mode = ""
//
//		_, err := NewWebserver()
//		if err == nil {
//			t.Fatal("expected error due to missing MODE")
//		}
//	})
//}
//
//func TestListenServeAndClose(t *testing.T) {
//	env.Init()
//	defer env.Reset()
//
//	// Mock UseLOCALHOST hvis nødvendigt
//	// orig := certificate.UseLOCALHOST
//	// certificate.UseLOCALHOST = func() (*tls.Certificate, error) { return &tls.Certificate{}, nil }
//	// defer func() { certificate.UseLOCALHOST = orig }()
//
//	wsIface, err := NewWebserver()
//	if err != nil {
//		t.Fatalf("failed to create webserver: %v", err)
//	}
//	ws := wsIface.(*webserver)
//
//	ws.httpServer.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		w.WriteHeader(http.StatusOK)
//	})
//	ws.httpsServer.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		w.WriteHeader(http.StatusOK)
//	})
//
//	go func() {
//		_ = ws.ListenAndServe()
//	}()
//
//	time.Sleep(300 * time.Millisecond)
//
//	checkPortOpen(t, ws.httpServer.Addr)
//	checkPortOpen(t, ws.httpsServer.Addr)
//
//	err = ws.Close()
//	if err != nil {
//		t.Fatalf("Close() returned error: %v", err)
//	}
//}
//
//func checkPortOpen(t *testing.T, addr string) {
//	conn, err := net.DialTimeout("tcp", "localhost"+addr, 500*time.Millisecond)
//	if err != nil {
//		t.Fatalf("expected server to listen on %s, but connection failed: %v", addr, err)
//	}
//	_ = conn.Close()
//}

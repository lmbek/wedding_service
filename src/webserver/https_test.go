package webserver

import (
	"testing"
	"wedding_service/certificate"
	"wedding_service/env"
)

func Test_newHttpsServer(t *testing.T) {
	t.Chdir("..")
	env.Init()
	acme, err := certificate.InitAcme()
	server, err := newHttpsServer("8443", acme)
	if err != nil {
		t.Errorf("got err: %s", err)
	}

	t.Run("use invalid cert and key", func(t *testing.T) {
		defer env.Reset()
		env.Env.CertPath = "invalid_cert"
		env.Env.KeyPath = "invalid_key"

		_, err := newHttpsServer("8443", acme)
		if err == nil {
			t.Errorf("should get an error")
		}
	})

	t.Run("testing server.TLSConfig.GetCertificate", func(t *testing.T) {
		_, err = server.TLSConfig.GetCertificate(nil)
		if err != nil {
			t.Errorf("got err: %s", err)
		}
	})

	t.Run("testing production mode", func(t *testing.T) {
		defer env.Reset()

		_, err := newHttpsServer("8443", acme)
		if err != nil {
			t.Errorf("gor err: %s", err)
		}
	})
}

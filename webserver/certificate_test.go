package webserver

import (
	"net/http"
	"testing"
	"wedding_service/env"
)

func Test_useCertificate(t *testing.T) {
	t.Chdir("..")
	env.Init()

	var httpsServer *http.Server = newHttpsServer(env.Env.HttpsPort)

	err := useCertificate(httpsServer, env.Env.CertPath, env.Env.KeyPath)
	if err != nil {
		t.Errorf("could not useCertificate %s", err)
	}

	// test if development mode works
	t.Run("development mode set", func(t *testing.T) {
		defer env.Reset()
		env.Env.Mode = "development"
		err := useCertificate(httpsServer, env.Env.CertPath, env.Env.KeyPath)
		if err != nil {
			t.Errorf("could not useCertificate for mode development: %s", err)
		}
	})

	// test if development mode with bad hostname gives err
	t.Run("development mode set and bad hostnames given", func(t *testing.T) {
		defer env.Reset()
		env.Env.Mode = "development"
		env.Env.CertPath = "invalid_cert_path"
		env.Env.KeyPath = "invalid_key_path"
		err := useCertificate(httpsServer, env.Env.CertPath, env.Env.KeyPath)
		if err == nil {
			t.Errorf("err should not be nil")
		}
	})

	t.Run("test useCertificate with production mode set", func(t *testing.T) {
		defer env.Reset()
		env.Env.Mode = "production"
		err := useCertificate(httpsServer, env.Env.CertPath, env.Env.KeyPath)
		if err != nil {
			t.Errorf("could not useCertificate: %s", err)
		}
	})
	t.Run("wrong https port", func(t *testing.T) {
		defer env.Reset()
		env.Env.Mode = "production"
		env.Env.Hostnames = nil
		err := useCertificate(httpsServer, env.Env.CertPath, env.Env.KeyPath)
		if err == nil {
			t.Errorf("err should not be nil")
		}
	})

	t.Run("test useCertificate with no mode set", func(t *testing.T) {
		defer env.Reset()
		env.Env.Mode = ""
		err := useCertificate(httpsServer, env.Env.CertPath, env.Env.KeyPath)
		if err == nil {
			t.Errorf("err should not be nil")
		}
	})

	t.Run("test useCertificate development mode certificate error", func(t *testing.T) {
		defer env.Reset()
		env.Env.Mode = "development"
		env.Env.CertPath = ""
		err := useCertificate(httpsServer, env.Env.CertPath, env.Env.KeyPath)
		if err == nil {
			t.Errorf("err should not be nil")
		}
	})

	t.Run("test useCertificate production mode certificate error", func(t *testing.T) {
		defer env.Reset()
		env.Env.Mode = "production"

	})

	t.Run("test wrapCert", func(t *testing.T) {
		_, err := wrapCert(nil)(nil)
		if err != nil {
			t.Errorf("could not wrapCert: %s", err)
			return
		}
	})
}

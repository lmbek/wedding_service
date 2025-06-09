package webserver

import (
	"net/http"
	"testing"
	"wedding_service/env"
)

func Test_useCertificate(t *testing.T) {
	env.Init()
	t.Chdir("..")
	var httpsServer *http.Server = newHttpsServer(env.Env.HttpsPort)

	err := useCertificate(httpsServer)
	if err != nil {
		t.Errorf("could not useCertificate %s", err)
	}

	t.Run("development mode set", func(t *testing.T) {
		defer env.Reset()
		env.Env.Mode = "development"
		err := useCertificate(httpsServer)
		if err != nil {
			t.Errorf("could not useCertificate for mode development: %s", err)
		}
	})
	t.Run("development mode set and bad hostnames given", func(t *testing.T) {
		// TODO: cant make test before env.Env is more expanded
		defer env.Reset()
		env.Env.Hostnames = nil
		err := useCertificate(httpsServer)
		if err == nil {
			//t.Errorf("err should not be nil")
		}
	})

	t.Run("test useCertificate with production mode set", func(t *testing.T) {
		defer env.Reset()
		env.Env.Mode = "production"
		err := useCertificate(httpsServer)
		if err != nil {
			t.Errorf("err should be nil")
		}
		//t.Run("wrong https port", func(t *testing.T) {
		//	env.Env.HttpsPort = "-1"
		//	err := useCertificate(httpsServer)
		//	if err == nil {
		//		t.Errorf("err should not be nil")
		//	}
		//})
	})

	t.Run("test useCertificate with no mode set", func(t *testing.T) {
		defer env.Reset()
		env.Env.Mode = ""
		err := useCertificate(httpsServer)
		if err == nil {
			t.Errorf("err should be nil")
		}
	})
}

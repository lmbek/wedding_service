package webserver

import (
	"testing"
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

func TestWebserver_Close(t *testing.T) {

}

func TestWebserver_ListenAndServe(t *testing.T) {

}

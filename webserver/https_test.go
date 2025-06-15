package webserver

import "testing"

func Test_newHttpsServer(t *testing.T) {
	server := newHttpsServer("8443")
	if server == nil {
		t.Errorf("server should not be nil")
	}
}

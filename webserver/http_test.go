package webserver

import "testing"

func Test_newHttpServer(t *testing.T) {
	server := newHttpServer("8080")
	if server == nil {
		t.Errorf("server should not be nil")
	}
}

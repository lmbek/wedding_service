package webserver

import (
	"testing"
	"wedding_service/env"
)

func Test_newHttpServer(t *testing.T) {
	t.Chdir("..")
	env.Init()

	newHttpServer("8080")
}

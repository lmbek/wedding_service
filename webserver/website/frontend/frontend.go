// NOTE: Normally the frontend should have a service called gompa attached
// to make frontend development, and then we just serve the out directory.
// It's important to note that in this project we only use this frontend
// for serving our assets such as images, css, js

package frontend

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"wedding_service/flags"
)

// frontendEmbedded is the embedded frontend for this application so performance is high and dependencies lower
//
//go:embed out/public/*
var frontendEmbeddedPublic embed.FS

type Frontend interface {
	GetFileSystem() http.FileSystem
}

type frontend struct {
	filesystem *fs.FS
}

// NewFrontend uses embedding or the flag for the frontend path
// example: go run . <insert flag with absolute path to directory that holds the frontend>
func NewFrontend() Frontend {
	flags.Parse()

	var frontendFS fs.FS

	// if frontend flag is set then use files dynamically
	if *flags.FrontendFlag != "" {
		frontendFS, err := fs.Sub(os.DirFS(*flags.FrontendFlag), "out/public")
		if err != nil {
			log.Fatalf("frontend path does not lead to a directory: %v\n", err)
		}
		return &frontend{filesystem: &frontendFS}
	}

	// otherwise use embedding
	frontendFS, err := fs.Sub(frontendEmbeddedPublic, "out/public")
	if err != nil {
		log.Fatalf("embedding did not work as we did not have a directory: %v\n", err)
	}
	return &frontend{filesystem: &frontendFS}
}

func (f *frontend) GetFileSystem() http.FileSystem {
	return http.FS(*f.filesystem)
}

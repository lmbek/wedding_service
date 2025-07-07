// NOTE: Normally the frontend should have a service called gompa attached
// to make frontend development, and then we just serve the out directory.
// It's important to note that in this project we only use this frontend
// for serving our assets such as images, css, js

package frontend

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"wedding_service/flags"
)

var DefaultFrontend Frontend = newFrontend()

// frontendEmbedded is the embedded frontend for this application, so performance is high and dependencies lower
//
// Example path: go run . -frontend=webserver/website/frontend/out
//
//go:embed out/public/*
var frontendEmbeddedPublic embed.FS

//go:embed out/private/*
var frontendEmbeddedPrivate embed.FS

type Frontend interface {
	Serve(w http.ResponseWriter, r *http.Request)
	GetPublicFileSystem() fs.FS
	GetPrivateFileSystem() fs.FS
}

type frontend struct {
	publicFilesystem  fs.FS
	privateFilesystem fs.FS
	handler           http.Handler
}

// NewFrontend uses embedding or the flag for the frontend path
// example: go run . <insert a flag with an absolute path to directory that holds the frontend>
func newFrontend() Frontend {
	frontendFlag := flags.LoadFrontendFlag()

	var frontendPublicFileSystem fs.FS
	var frontendPrivateFileSystem fs.FS
	var err error

	// if a frontendFlag is set, then use files dynamically
	if frontendFlag != "" {
		fmt.Println("[Frontend] Please note: hotreload enabled. Refresh browser between runs")
		fmt.Printf("Loading frontend from %s\n", frontendFlag)
		frontendPublicFileSystem, err = fs.Sub(os.DirFS(frontendFlag), "public")
		if err != nil {
			log.Fatalf("could not load public folder: %v\n", err)
		}
		frontendPrivateFileSystem, err = fs.Sub(os.DirFS(frontendFlag), "private")
		if err != nil {
			log.Fatalf("could not load private folder: %v\n", err)
		}

		// start hot reloader
		go CheckForModification(frontendFlag)
	} else {
		fmt.Println("[Frontend] Please note: hotreload disabled. Using embedded frontend! ")
		// otherwise use embedding
		frontendPublicFileSystem, err = fs.Sub(frontendEmbeddedPublic, "out/public")
		if err != nil {
			log.Fatalf("embedding did not work as we did not have a directory: %v\n", err)
		}
		frontendPrivateFileSystem, err = fs.Sub(frontendEmbeddedPrivate, "out/private")
		if err != nil {
			log.Fatalf("embedding did not work as we did not have a directory: %v\n", err)
		}
	}

	return &frontend{
		publicFilesystem:  frontendPublicFileSystem,
		privateFilesystem: frontendPrivateFileSystem,
		handler: http.FileServer(
			NewFileSystem(http.FS(frontendPublicFileSystem), false, false),
		),
	}
}

func (f *frontend) GetPublicFileSystem() fs.FS {
	return f.publicFilesystem
}

func (f *frontend) GetPrivateFileSystem() fs.FS {
	return f.privateFilesystem
}

func (f *frontend) Serve(w http.ResponseWriter, r *http.Request) {
	f.handler.ServeHTTP(w, r)
}

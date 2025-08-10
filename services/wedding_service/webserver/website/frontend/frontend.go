// NOTE: Normally the frontend should have a service called gompa attached
// to make frontend development, and then we just serve the out directory.
// It's important to note that in this project we only use this frontend
// for serving our assets such as images, css, js

package frontend

import (
	"embed"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
)

// frontendEmbedded is the embedded frontend for this application, so performance is high and dependencies lower
//
// Example path: go run . -frontend=webserver/website/frontend/out -hotreload=true
//
//go:embed out/public/*
var publicEmbeddedFS embed.FS

//go:embed out/private/*
var privateEmbeddedFS embed.FS

type Frontend interface {
	// load initializes the frontend using either embedded or file-based assets
	load(frontendPath string) (Frontend, error)
	// Serve handles HTTP requests for frontend assets
	Serve(w http.ResponseWriter, r *http.Request)
	// PublicFS returns the filesystem containing public assets
	PublicFS() fs.FS
	// PrivateFS returns the filesystem containing private assets
	PrivateFS() fs.FS
}

type frontend struct {
	hotReloadEnabled bool
	handler          http.Handler
	publicFS         fs.FS
	privateFS        fs.FS
}

// NewFrontend uses embedding or the flag for the frontend path
// example: go run . <insert a flag with an absolute path to directory that holds the frontend>
func NewFrontend(path string, hotReloadEnabled bool) (Frontend, error) {
	f := &frontend{
		hotReloadEnabled: hotReloadEnabled,
	}
	return f.load(path)
}

func (f *frontend) load(frontendPath string) (Frontend, error) {
	if frontendPath == "" {
		return f, f.UseEmbedded()
	}
	return f, f.useFileBased(frontendPath)
}

func (f *frontend) useFileBased(frontendPath string) error {
	var publicFS fs.FS
	var privateFS fs.FS
	var err error

	publicFS, err = fs.Sub(os.DirFS(frontendPath), "public")
	if err != nil {
		slog.Error("could not load public folder", slog.Any("error", err))
		return err
	}
	privateFS, err = fs.Sub(os.DirFS(frontendPath), "private")
	if err != nil {
		slog.Error("could not load private folder", slog.Any("error", err))
		return err
	}

	if f.hotReloadEnabled {
		slog.Info("[Frontend] Please note: hotreload enabled. Refresh browser between runs")

		// start hot reloader
		go CheckForModification(frontendPath)
	}

	// expose only the public part as a fileserver with a handler
	f.handler = http.FileServer(
		NewFileSystem(http.FS(publicFS), false, false),
	)
	f.publicFS = publicFS
	f.privateFS = privateFS

	return nil
}

func (f *frontend) UseEmbedded() error {
	var publicFS fs.FS
	var privateFS fs.FS
	var err error
	slog.Info("[Frontend] Please note: hotreload disabled. Using embedded frontend!")
	// otherwise use embedding
	publicFS, err = fs.Sub(publicEmbeddedFS, "out/public")
	if err != nil {
		slog.Error("embedding did not work as we did not have a directory", slog.Any("error", err))
		return err
	}
	privateFS, err = fs.Sub(privateEmbeddedFS, "out/private")
	if err != nil {
		slog.Error("embedding did not work as we did not have a directory", slog.Any("error", err))
		return err
	}

	// expose only the public part as a fileserver with a handler
	f.handler = http.FileServer(
		NewFileSystem(http.FS(publicFS), false, false),
	)
	f.publicFS = publicFS
	f.privateFS = privateFS
	return nil
}

func (f *frontend) PublicFS() fs.FS {
	return f.publicFS
}

func (f *frontend) PrivateFS() fs.FS {
	return f.privateFS
}

func (f *frontend) Serve(w http.ResponseWriter, r *http.Request) {
	f.handler.ServeHTTP(w, r)
}

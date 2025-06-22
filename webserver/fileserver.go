package webserver

import (
	"net/http"
	"wedding_service/webserver/website/frontend"
)

// FileServer is an interface for serving files.
type FileServer interface {
	Serve(w http.ResponseWriter, r *http.Request)
}

// fileServer is the concrete implementation of FileServer.
type fileServer struct {
	fileSystem http.FileSystem
}

// NewFileServer returns a new instance of FileServer
func NewFileServer() FileServer {
	return &fileServer{fileSystem: &NoIndexing{frontend.NewFrontend().GetFileSystem()}}
}

func (fs *fileServer) Serve(w http.ResponseWriter, r *http.Request) {
	http.FileServer(fs.fileSystem).ServeHTTP(w, r)
}

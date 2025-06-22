package webserver

import (
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

// NoIndexing disables directory listing unless index.html is present.
type NoIndexing struct {
	fs http.FileSystem
}

// NewNoIndexing returns a wrapped file system that disables directory listing.
func NewNoIndexing(root string) http.FileSystem {
	return &NoIndexing{fs: http.Dir(root)}
}

func (n *NoIndexing) Open(name string) (http.File, error) {
	name = path.Clean(name)

	f, err := n.fs.Open(name)
	if err != nil {
		return nil, os.ErrNotExist
	}

	stat, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, os.ErrNotExist
	}

	// Block symlinks for the requested path
	if stat.Mode()&os.ModeSymlink != 0 {
		f.Close()
		return nil, os.ErrNotExist
	}

	if !stat.IsDir() {
		return f, nil
	}

	// It's a directory: stream Readdir 1 entry at a time and check for index.html
	const batchSize = 1
	for {
		entries, err := f.Readdir(batchSize)
		if err == io.EOF {
			break
		}
		if err != nil {
			f.Close()
			return nil, os.ErrNotExist
		}

		for _, entry := range entries {
			if strings.EqualFold(entry.Name(), "index.html") {
				// Check if index.html is a symlink — block if so
				indexFile, err := n.fs.Open(path.Join(name, entry.Name()))
				if err == nil {
					defer indexFile.Close()
					indexStat, err := indexFile.Stat()
					if err == nil && indexStat.Mode()&os.ModeSymlink == 0 {
						return f, nil
					}
				}
			}
		}
	}

	f.Close()
	return nil, os.ErrNotExist
}

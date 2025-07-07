package frontend

import (
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

// FileSystem is a configurable http.FileSystem wrapper that can disable directory indexing
// and optionally allow or disallow symlinks.
type FileSystem struct {
	fs            http.FileSystem
	allowIndexing bool
	allowSymlinks bool
}

func NewFileSystem(fs http.FileSystem, allowIndexing bool, allowSymlinks bool) http.FileSystem {
	return &FileSystem{
		fs:            fs,
		allowIndexing: allowIndexing,
		allowSymlinks: allowSymlinks,
	}
}

func (f *FileSystem) Open(name string) (http.File, error) {
	name = path.Clean(name)

	file, err := f.fs.Open(name)
	if err != nil {
		return nil, os.ErrNotExist
	}

	stat, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, os.ErrNotExist
	}

	// Block symlinks for the requested path if not allowed
	if stat.Mode()&os.ModeSymlink != 0 && !f.allowSymlinks {
		file.Close()
		return nil, os.ErrNotExist
	}

	// If not a directory or directory indexing is allowed, return the file
	if !stat.IsDir() || f.allowIndexing {
		return file, nil
	}

	// Directory and indexing is disabled: check for a valid index.html
	if f.hasValidIndex(name) {
		return file, nil
	}

	file.Close()
	return nil, os.ErrNotExist
}

// hasValidIndex checks if directory contains a valid index.html file
func (f *FileSystem) hasValidIndex(dir string) bool {
	const batchSize = 1

	dirFile, err := f.fs.Open(dir)
	if err != nil {
		return false
	}
	defer dirFile.Close()

	for {
		entries, err := dirFile.Readdir(batchSize)
		if err == io.EOF {
			break
		}
		if err != nil {
			return false
		}

		for _, entry := range entries {
			if !strings.EqualFold(entry.Name(), "index.html") {
				continue
			}

			indexFile, err := f.fs.Open(path.Join(dir, entry.Name()))
			if err != nil {
				continue
			}
			defer indexFile.Close()

			indexStat, err := indexFile.Stat()
			if err != nil {
				continue
			}

			// Accept if symlinks allowed or index.html is not a symlink
			// TODO: can this be made simpler?
			if f.allowSymlinks || (indexStat.Mode()&os.ModeSymlink == 0) {
				return true
			}
		}
	}

	return false
}

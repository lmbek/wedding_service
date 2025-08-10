package frontend

import (
	"net/http"
	"os"
	"path"
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
	// Try to open the index.html file directly
	indexPath := path.Join(dir, "index.html")
	indexFile, err := f.fs.Open(indexPath)
	if err != nil {
		return false
	}
	defer indexFile.Close()

	// Check if it's a valid file
	indexStat, err := indexFile.Stat()
	if err != nil {
		return false
	}

	// Check if it's a regular file (not a directory)
	if indexStat.IsDir() {
		return false
	}

	// Accept if symlinks are allowed or if the file is not a symlink
	return f.allowSymlinks || (indexStat.Mode()&os.ModeSymlink == 0)
}

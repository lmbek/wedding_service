package main

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"sort"
	"time"
)

const (
	rootDir      = "mytestdir"
	pollInterval = 50 * time.Millisecond
)

type FileSnapshot struct {
	Path    string
	IsDir   bool
	ModTime time.Time
}

func main() {
	fmt.Println("Started watching:", rootDir)

	lastSnap, err := buildSnapshot(rootDir)
	if err != nil {
		log.Fatalln("Initial scan failed:", err)
	}

	lastMod := getLatestModTime(lastSnap)

	for {
		newSnap := waitForChange(lastSnap, lastMod, rootDir)
		log.Println("Detected change at", time.Now())
		lastSnap = newSnap
		lastMod = getLatestModTime(newSnap)
	}
}

func waitForChange(lastSnap []FileSnapshot, lastMod time.Time, root string) []FileSnapshot {
	changed := make(chan []FileSnapshot)

	go func() {
		ticker := time.NewTicker(pollInterval)
		defer ticker.Stop()

		for {
			<-ticker.C
			newSnap, err := buildSnapshot(root)
			if err != nil {
				log.Println("Snapshot error:", err)
				continue
			}

			newMod := getLatestModTime(newSnap)

			// First try fast path: modtime change
			if !newMod.Equal(lastMod) {
				changed <- newSnap
				return
			}

			// Fallback to full snapshot comparison
			if !snapshotsEqual(lastSnap, newSnap) {
				changed <- newSnap
				return
			}
		}
	}()

	return <-changed
}

func buildSnapshot(root string) ([]FileSnapshot, error) {
	var snapshot []FileSnapshot
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return nil
		}
		snapshot = append(snapshot, FileSnapshot{
			Path:    path,
			IsDir:   d.IsDir(),
			ModTime: info.ModTime(),
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Sort for stable comparison
	sort.Slice(snapshot, func(i, j int) bool {
		return snapshot[i].Path < snapshot[j].Path
	})

	return snapshot, nil
}

func getLatestModTime(snap []FileSnapshot) time.Time {
	var latest time.Time
	for _, f := range snap {
		if f.ModTime.After(latest) {
			latest = f.ModTime
		}
	}
	return latest
}

func snapshotsEqual(a, b []FileSnapshot) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Path != b[i].Path || a[i].IsDir != b[i].IsDir || !a[i].ModTime.Equal(b[i].ModTime) {
			return false
		}
	}
	return true
}

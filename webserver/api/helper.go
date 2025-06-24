package api

import (
	"errors"
	"strconv"
	"strings"
)

// parseIDFromPath extracts the ID from the URL path assuming trailing slash and numeric ID.
func parseIDFromPath(path string) (int, error) {
	segments := strings.Split(strings.Trim(path, "/"), "/")
	if len(segments) == 0 {
		return 0, errors.New("path is empty")
	}
	idStr := segments[len(segments)-1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// getNextID generates the next ID for an in-memory slice.
func getNextID() int {
	maxID := 0
	for _, p := range persons {
		if p.ID > maxID {
			maxID = p.ID
		}
	}
	return maxID + 1
}

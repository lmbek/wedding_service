package api

// getNextID generates the next ID for an in-memory slice.
func getNextID() int {
	maxID := 0
	for _, p := range persons {
		if p.Id > maxID {
			maxID = p.Id
		}
	}
	return maxID + 1
}

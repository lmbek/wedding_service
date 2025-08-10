//go:build production

package buildtag

func IsDevelopment() bool {
	return false
}

func IsDockerDev() bool {
	return false
}

func IsProduction() bool {
	return true
}

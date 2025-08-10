//go:build dockerdev

package buildtag

func IsDevelopment() bool {
	return false
}

func IsDockerDev() bool {
	return true
}

func IsProduction() bool {
	return false
}

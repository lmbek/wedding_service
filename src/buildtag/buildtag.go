//go:build !production && !dockerdev

package buildtag

func IsDevelopment() bool {
	return true
}

func IsDockerDev() bool {
	return false
}

func IsProduction() bool {
	return false
}

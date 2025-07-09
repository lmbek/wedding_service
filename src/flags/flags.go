package flags

import (
	"flag"
	"fmt"
	"sync"
)

var (
	frontendFlag = flag.String("frontend", "", "Path to frontend files")
	parseOnce    sync.Once
)

// LoadFrontendFlag parses flags (once) and returns the frontend flag value.
func LoadFrontendFlag() string {
	fmt.Println("who is the caller")
	parseOnce.Do(func() {
		flag.Parse()
		fmt.Println("parsed flag")
	})
	return *frontendFlag
}

func HotReloadEnabled() bool {
	if LoadFrontendFlag() == "" {
		return false
	}
	return true
}

func TemplateCachingEnabled() bool {
	if LoadFrontendFlag() == "" {
		return true
	}
	return false
}

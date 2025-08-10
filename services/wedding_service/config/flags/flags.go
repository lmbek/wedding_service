package flags

import (
	"flag"
)

type Flags interface {
	FrontendPath() string
	HotReloadEnabled() bool
}

type flags struct {
	frontendPath     string
	hotReloadEnabled bool
}

func NewFlags() Flags {
	frontendPathFlag := flag.String("frontend", "", "Path to frontend files")
	hotReloadFlag := flag.Bool("hotreload", false, "Hot reloading frontend")
	flag.Parse()
	return &flags{
		frontendPath:     *frontendPathFlag,
		hotReloadEnabled: *hotReloadFlag,
	}
}

func (f *flags) FrontendPath() string {
	return f.frontendPath
}

func (f *flags) HotReloadEnabled() bool {
	return f.hotReloadEnabled
}

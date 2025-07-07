package flags

import "flag"

var FrontendFlag string
var frontendFlag = flag.String("frontend", "", "Path to frontend files")

// Only parse flags once, safely
var parsed bool

func LoadFrontendFlag() string {
	if !parsed {
		flag.Parse()
		parsed = true
	}
	FrontendFlag = *frontendFlag
	return FrontendFlag
}

package flags

import "flag"

var FrontendFlag string
var frontendFlag = flag.String("frontend", "", "Path to frontend files")

func LoadFrontendFlag() string {
	FrontendFlag = *frontendFlag
	return *frontendFlag
}

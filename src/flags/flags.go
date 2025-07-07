package flags

import "flag"

var FrontendFlag string
var frontendFlag = flag.String("frontend", "", "Path to frontend files")

func parse() {
	flag.Parse()
}

func LoadFrontendFlag() string {
	parse()
	FrontendFlag = *frontendFlag
	return *frontendFlag
}

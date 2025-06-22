package flags

import "flag"

var FrontendFlag = flag.String("frontend", "", "Path to frontend files")

func Parse() {
	flag.Parse()
}

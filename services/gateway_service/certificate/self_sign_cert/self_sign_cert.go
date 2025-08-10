//go:build !windows && !linux && self_sign_cert

// TODO: the self_sign_cert ignore is for test coverage ignore. Tests can be written to remove this

// READ THIS FIRST:
// create the certificate and key files for localhost first!
// You need to generate self-signed certificates
// by running go generate

//////
//// go:generate go run ./self_sign_cert/self_sign_cert.go
//go:generate go run self_sign_cert.go

// NOTE: this file is ignored by the default build
// as it is a tool, we don't count this file in our test coverage total

package main

import (
	"fmt"
)

func main() {
	fmt.Printf("self sign cert generator: OS not supported yet, please use linux or windows for now.\n")
}

//go:build !go1.14
// +build !go1.14

package main

// Make sure people don't try to build with older versions of Go, as that will
// introduce some runtime problems (e.g. using %w) and/or give confusing "no
// such function" errors.
func init() {
	"You need Go 1.14 or newer to compile uni"
}

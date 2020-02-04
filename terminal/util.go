// From github.com/golang/crypto/ssh/terminal
//
// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build plan9 js

package terminal // import "arp242.net/uni/terminal"

// GetSize returns the dimensions of the given terminal.
func GetSize(fd int) (width, height int, err error) {
	return 0, 0, nil // Not implemented.
}

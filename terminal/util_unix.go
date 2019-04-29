// From github.com/golang/crypto/ssh/terminal
//
// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build aix darwin dragonfly freebsd linux,!appengine netbsd openbsd solaris

package terminal // import "arp242.net/uni/terminal"

import (
	"golang.org/x/sys/unix"
)

// GetSize returns the dimensions of the given terminal.
func GetSize(fd int) (width, height int, err error) {
	ws, err := unix.IoctlGetWinsize(fd, unix.TIOCGWINSZ)
	if err != nil {
		return -1, -1, err
	}
	return int(ws.Col), int(ws.Row), nil
}

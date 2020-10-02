// +build darwin linux

package application

import "golang.org/x/sys/unix"

var (
	sigInt  = unix.SIGINT
	sigQuit = unix.SIGQUIT
	sigTerm = unix.SIGTERM
)

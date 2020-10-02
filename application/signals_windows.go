// +build windows

package application

import "golang.org/x/sys/windows"

var (
	sigInt  = windows.SIGINT
	sigQuit = windows.SIGQUIT
	sigTerm = windows.SIGTERM
)

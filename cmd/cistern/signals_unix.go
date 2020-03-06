// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package main

import "os/signal"
import "syscall"

func SetupSignalHandlers() {
	signal.Ignore(syscall.SIGINT)
	// FIXME Do not ignore SIGTSTP/SIGCONT
	signal.Ignore(syscall.SIGTSTP)
}

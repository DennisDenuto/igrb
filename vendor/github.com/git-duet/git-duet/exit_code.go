// +build !windows,!plan9

package duet

import (
	"os/exec"
	"syscall"
)

func exitCode(err *exec.ExitError) int {
	return err.Sys().(syscall.WaitStatus).ExitStatus()
}

//go:build !windows

package remediation

import (
	"os/exec"
	"syscall"
	"time"
)

// prepareCommand ensures that the shell and any child processes it spawns
// run in their own process group, so that context cancellation terminates the
// whole process tree rather than leaving child processes (like sleep) holding
// open stdout/stderr pipes.
func prepareCommand(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Cancel = func() error {
		if cmd.Process != nil && cmd.Process.Pid > 0 {
			return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		}
		return nil
	}
	cmd.WaitDelay = 250 * time.Millisecond
}

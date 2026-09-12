//go:build windows

package remediation

import (
	"os/exec"
	"time"
)

func prepareCommand(cmd *exec.Cmd) {
	cmd.WaitDelay = 250 * time.Millisecond
}

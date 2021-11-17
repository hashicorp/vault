// +build windows

package child

import "os/exec"

func setSetpgid(cmd *exec.Cmd, value bool) {}

func processNotFoundErr(err error) bool {
	return false
}

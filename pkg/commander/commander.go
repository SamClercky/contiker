package commander

import (
	"os/exec"

	"github.com/SamClercky/contiker/pkg/pkgmanager"
)

type Command struct {
	program map[int]string // The names of the packages for specific OSes
	command string         // The command that will be used to check if it is installed
	execCmd exec.Cmd       // The command that will be executed
}

func (cmd *Command) Exists() bool {
	_, err := exec.LookPath(cmd.command)
	return err == nil
}

func (cmd *Command) EnsureInstalled(manager pkgmanager.PkgManager) error {
	if cmd.Exists() {
		return nil
	}

	return manager.Install(cmd.program)
}

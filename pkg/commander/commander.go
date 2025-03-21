package commander

import (
	"fmt"
	"os/exec"

	commandcontext "github.com/SamClercky/contiker/pkg/command-context"
	"github.com/SamClercky/contiker/pkg/pkgmanager"
)

type Command struct {
	program map[int]string    // The names of the packages for specific OSes
	command string            // The command that will be used to check if it is installed
	execCmd exec.Cmd          // The command that will be executed
	context commandcontext.CC // Context of the command
}

func (cmd *Command) Exists() bool {
	exists, err := cmd.context.CommandExists(cmd.command)
	if err != nil {
		fmt.Printf("[ERROR] Could not check if command exists with error: %a\n", err)
		return false
	}

	return exists
}

func (cmd *Command) EnsureInstalled(manager pkgmanager.PkgManager) error {
	if cmd.Exists() {
		return nil
	}

	return manager.Install(cmd.context, cmd.program)
}

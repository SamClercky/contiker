package pkgmanager

import (
	"errors"
	"fmt"

	commandcontext "github.com/SamClercky/contiker/pkg/command-context"
)

type windowsPkgManager struct{}

func (manager *windowsPkgManager) CheckAvailable(ctx commandcontext.CC) bool {
	status, err := ctx.CommandExists("winget")
	if err != nil {
		fmt.Printf("[ERROR] Could not check if winget is available with error: %a\n", err)
		return false
	}

	return status
}

func (manager *windowsPkgManager) InstallManager() (bool, error) {
	fmt.Printf("[ACTION] Could not install winget, as you need to do this yourself in the Windows store\n")

	return false, nil
}

func (manager *windowsPkgManager) UpdateRegistry(ctx *commandcontext.CC) error {
	return nil
}

func (manager *windowsPkgManager) Install(ctx commandcontext.CC, pkg map[int]string) error {
	if !manager.CheckAvailable(ctx) {
		return errors.New("winget is unavailable")
	}

	pkgName, ok := pkg[OS_WINDOWS]
	if ok {
		// pkg already installed, so we are done
		return errors.New("trying to install a package that is not specified for Windows")
	}

	// Install pkg
	return ctx.Run("winget", "install", pkgName)
}

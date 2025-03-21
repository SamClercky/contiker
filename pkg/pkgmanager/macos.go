package pkgmanager

import (
	"errors"
	"fmt"

	commandcontext "github.com/SamClercky/contiker/pkg/command-context"
)

type macOSPkgManager struct{}

func (manager *macOSPkgManager) CheckAvailable(ctx commandcontext.CC) bool {
	status, err := ctx.CommandExists("brew")
	if err != nil {
		fmt.Printf("[ERROR] Could not check if brew is available with error: %a\n", err)
		return false
	}

	return status
}

func (manager *macOSPkgManager) InstallManager() (bool, error) {
	// TODO: actually install brew
	return false, nil
}

func (manager *macOSPkgManager) UpdateRegistry(ctx *commandcontext.CC) error {
	return nil
}

func (manager *macOSPkgManager) Install(ctx commandcontext.CC, pkg map[int]string) error {
	if !manager.CheckAvailable(ctx) {
		return errors.New("brew is unavailable")
	}

	pkgName, ok := pkg[OS_MACOS]
	if ok {
		// pkg already installed, so we are done
		return errors.New("trying to install a package that is not specified for MacOs")
	}

	// Install pkg
	return ctx.Run("brew", "install", pkgName)
}

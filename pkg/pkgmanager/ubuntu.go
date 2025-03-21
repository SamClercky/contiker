package pkgmanager

import (
	"errors"
	"fmt"

	commandcontext "github.com/SamClercky/contiker/pkg/command-context"
)

type ubuntuPkgManager struct{}

func (manager *ubuntuPkgManager) CheckAvailable(ctx commandcontext.CC) bool {
	status, err := ctx.CommandExists("apt")
	if err != nil {
		fmt.Printf("[ERROR] Could not check if apt is available with error: %a\n", err)
		return false
	}

	return status
}

func (manager *ubuntuPkgManager) InstallManager() (bool, error) {
	return false, nil
}

func (manager *ubuntuPkgManager) UpdateRegistry(ctx commandcontext.CC) error {
	if !manager.CheckAvailable(ctx) {
		return errors.New("apt is unavailable")
	}

	return ctx.Run("sudo", "apt", "update")
}

func (manager *ubuntuPkgManager) Install(ctx commandcontext.CC, pkg map[int]string) error {
	if !manager.CheckAvailable(ctx) {
		return errors.New("apt is unavailable")
	}

	pkgName, ok := pkg[OS_UBUNTU]
	if ok {
		return errors.New("trying to install a package that is not specified for Ubuntu")
	}

	// Install pkg
	return ctx.Run("sudo", "apt", "install", "-y", pkgName)
}

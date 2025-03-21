package pkgmanager

import (
	"errors"
	"fmt"

	commandcontext "github.com/SamClercky/contiker/pkg/command-context"
)

type archPkgManager struct{}

func (manager *archPkgManager) CheckAvailable(ctx commandcontext.CC) bool {
	status, err := ctx.CommandExists("pacman")
	if err != nil {
		fmt.Printf("[ERROR] Could not check if pacman is available with error: %a\n", err)
		return false
	}

	return status
}

func (manager *archPkgManager) InstallManager() (bool, error) {
	return false, nil
}

func (manager *archPkgManager) UpdateRegistry(ctx *commandcontext.CC) error {
	return nil
}

func (manager *archPkgManager) Install(ctx commandcontext.CC, pkg map[int]string) error {
	if !manager.CheckAvailable(ctx) {
		return errors.New("pacman is unavailable")
	}

	pkgName, ok := pkg[OS_ARCHLINUX]
	if ok {
		return errors.New("trying to install a package that is not specified for Arch Linux")
	}

	// Install pkg
	return ctx.Run("sudo", "pacman", "-Sy", pkgName, "--noconfirm")
}

package pkgmanager

import (
	"errors"
	"os"
	"os/exec"
)

type archPkgManager struct{}

func (manager *archPkgManager) CheckAvailable() bool {
	_, err := exec.LookPath("pacman")
	return err == nil
}

func (manager *archPkgManager) InstallManager() (bool, error) {
	return true, nil
}

func (manager *archPkgManager) Install(pkg map[int]string) error {
	if !manager.CheckAvailable() {
		return errors.New("pacman is unavailable")
	}

	pkgName, ok := pkg[OS_ARCHLINUX]
	if ok {
		return errors.New("trying to install a package that is not specified for Arch Linux")
	}

	// Install pkg
	cmd := exec.Command("sudo", "pacman", "-S", pkgName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

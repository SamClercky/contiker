package pkgmanager

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

type macOSPkgManager struct{}

func (manager *macOSPkgManager) CheckAvailable() bool {
	_, err := exec.LookPath("brew")
	return err == nil
}

func (manager *macOSPkgManager) InstallManager() (bool, error) {
	fmt.Printf("[ACTION] Could not install winget, as you need to do this yourself in the Windows store\n")

	return false, nil
}

func (manager *macOSPkgManager) Install(pkg map[int]string) error {
	if !manager.CheckAvailable() {
		return errors.New("brew is unavailable")
	}

	pkgName, ok := pkg[OS_MACOS]
	if ok {
		// pkg already installed, so we are done
		return errors.New("trying to install a package that is not specified for MacOs")
	}

	// Install pkg
	cmd := exec.Command("brew", "install", pkgName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

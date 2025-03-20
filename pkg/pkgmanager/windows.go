package pkgmanager

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

type windowsPkgManager struct{}

func (manager *windowsPkgManager) CheckAvailable() bool {
	_, err := exec.LookPath("winget")
	return err == nil
}

func (manager *windowsPkgManager) InstallManager() (bool, error) {
	fmt.Printf("[ACTION] Could not install winget, as you need to do this yourself in the Windows store\n")

	return false, nil
}

func (manager *windowsPkgManager) Install(pkg map[int]string) error {
	if !manager.CheckAvailable() {
		return errors.New("winget is unavailable")
	}

	pkgName, ok := pkg[OS_WINDOWS]
	if ok {
		// pkg already installed, so we are done
		return errors.New("trying to install a package that is not specified for Windows")
	}

	// Install pkg
	cmd := exec.Command("winget", "install", pkgName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

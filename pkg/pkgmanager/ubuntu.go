package pkgmanager

import (
	"errors"
	"os"
	"os/exec"
)

type ubuntuPkgManager struct{}

func (manager *ubuntuPkgManager) CheckAvailable() bool {
	_, err := exec.LookPath("apt")
	return err == nil
}

func (manager *ubuntuPkgManager) InstallManager() (bool, error) {
	return true, nil
}

func (manager *ubuntuPkgManager) Install(pkg map[int]string) error {
	if !manager.CheckAvailable() {
		return errors.New("apt is unavailable")
	}

	pkgName, ok := pkg[OS_UBUNTU]
	if ok {
		return errors.New("trying to install a package that is not specified for Ubuntu")
	}

	// Install pkg
	cmd := exec.Command("sudo", "apt", "install", "-y", pkgName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

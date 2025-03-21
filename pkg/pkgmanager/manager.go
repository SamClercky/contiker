package pkgmanager

import commandcontext "github.com/SamClercky/contiker/pkg/command-context"

const (
	OS_ARCHLINUX = iota
	OS_WINDOWS
	OS_UBUNTU
	OS_FEDORA
	OS_MACOS
)

type PkgManager interface {
	// Update registry if it is possible do to so separately
	UpdateRegistry(ctx commandcontext.CC) error
	// Ensure that a specific package has been installed
	Install(ctx commandcontext.CC, pkg map[int]string) error
	// Check if current manager is available
	CheckAvailable() bool
	// Install current manager if not yet available
	InstallManager() (bool, error)
}

package pkgmanager

const (
	OS_ARCHLINUX = iota
	OS_WINDOWS
	OS_UBUNTU
	OS_FEDORA
	OS_MACOS
)

type PkgManager interface {
	// Ensure that a specific package has been installed
	Install(pkg map[int]string) error
	// Check if current manager is available
	CheckAvailable() bool
	// Install current manager if not yet available
	InstallManager() (bool, error)
}

package commandcontext

import (
	"errors"
	"os"
	"os/exec"
)

type DockerCC struct {
	containerName string
}

func NewDocker(containerName string) DockerCC {
	return DockerCC{
		containerName: containerName,
	}
}

func (cc *DockerCC) Run(name string, args ...string) error {
	cmd := exec.Command("docker", append(
		[]string{
			"exec",
			"-it",
			cc.containerName,
			name,
		}, args...)...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (cc *DockerCC) CommandExists(name string) (bool, error) {
	cmd := exec.Command("docker", "exec", "-it", cc.containerName, "which", name)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()
	code := cmd.ProcessState.ExitCode()
	switch code {
	case 0:
		return true, nil
	case 1:
		return false, nil
	default:
		return false, errors.New("couldn't check existance of command in container")
	}
}

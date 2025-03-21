package commandcontext

import (
	"os"
	"os/exec"
)

type DefaultCC struct{}

func NewDefault() DefaultCC {
	return DefaultCC{}
}

func (cc *DefaultCC) Run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (cc *DefaultCC) CommandExists(name string) (bool, error) {
	_, err := exec.LookPath(name)
	return err == nil, nil
}

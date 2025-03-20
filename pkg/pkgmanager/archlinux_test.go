package pkgmanager_test

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
)

func TestInstallPackageArchLinux(t *testing.T) {
	req := testcontainers.ContainerRequest{
		Image: "archlinux",
		Cmd:   []string{"sleep", "infinity"},
	}

	ctx := context.Background()
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	defer testcontainers.CleanupContainer(t, container)

	path := "/home/test"
	c, out_reader, err := container.Exec(ctx, []string{"pacman", "-Sy", "which", "--noconfirm"})
	// Print output
	buf := new(strings.Builder)
	_, err = io.Copy(buf, out_reader)
	require.NoError(t, err)
	// See the logs from the command execution.
	t.Log(buf.String())

	require.NoError(t, err)
	require.Zerof(t, c, "File %s should have been created successfully, expected return code 0, got %v", path, c)
}

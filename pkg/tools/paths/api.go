package paths

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// DefaultDirectoryPermissions defines the directory permissions used throughout the IAC repository
const DefaultDirectoryPermissions = 0o600

// AbsoluteRepositoryRootDirectory retrieves the root directory of the git repository xctl gets run in no matter what
// directory the user is in.
func AbsoluteRepositoryRootDirectory() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")

	stdout := bytes.Buffer{}

	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("running command: %w", err)
	}

	return strings.Trim(stdout.String(), "\n"), nil
}

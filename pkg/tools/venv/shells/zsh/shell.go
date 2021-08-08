package zsh

import (
	"fmt"
	"os"
	"path"
)

const zshConfigFilename = ".zshrc"

func (s *shell) initialize() (err error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("getting user home directory: %w", err)
	}

	f, err := s.fs.Open(path.Join(homeDir, zshConfigFilename))
	if err != nil {
		return fmt.Errorf("opening zsh config file: %w", err)
	}

	s.tmpDir, err = s.fs.TempDir("", "xctl")
	if err != nil {
		return fmt.Errorf("creating temporary directory: %w", err)
	}

	err = s.fs.WriteReader(path.Join(s.tmpDir, zshConfigFilename), f)
	if err != nil {
		return fmt.Errorf("copying zsh config: %w", err)
	}

	return nil
}

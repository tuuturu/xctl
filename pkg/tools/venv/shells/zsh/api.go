package zsh

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/afero"

	"github.com/deifyed/xctl/pkg/tools/venv"
)

func (s *shell) Command(env []string) (*exec.Cmd, error) {
	err := s.initialize()
	if err != nil {
		return nil, fmt.Errorf("initializing shell: %w", err)
	}

	cmd := exec.Command(s.shellBinPath)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	modifiedEnv := venv.MergeVariables(env, []string{fmt.Sprintf("ZDOTDIR=%s", s.tmpDir)})

	cmd.Env = modifiedEnv

	return cmd, nil
}

func (s *shell) Teardown() error {
	err := s.fs.RemoveAll(s.tmpDir)
	if err != nil {
		return fmt.Errorf("removing venv config directory: %w", err)
	}

	return nil
}

func NewZshShell(fs *afero.Afero, binPath string) venv.Shell {
	return &shell{
		shellBinPath: binPath,
		fs:           fs,
	}
}

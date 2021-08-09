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

	modifiedEnv := venv.MergeVariables(env, []string{fmt.Sprintf("ZDOTDIR=%s", s.workDir)})

	cmd.Env = modifiedEnv

	return cmd, nil
}

func NewZshShell(fs *afero.Afero, workDir, binPath string) venv.Shell {
	return &shell{
		fs:           fs,
		workDir:      workDir,
		shellBinPath: binPath,
	}
}

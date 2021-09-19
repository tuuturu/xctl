package zsh

import (
	"fmt"
	"os/exec"

	"github.com/deifyed/xctl/pkg/apis/xctl"

	"github.com/spf13/afero"

	"github.com/deifyed/xctl/pkg/tools/venv"
)

func (s *shell) Command(io xctl.IOStreams, env []string) (*exec.Cmd, error) {
	err := s.initialize()
	if err != nil {
		return nil, fmt.Errorf("initializing shell: %w", err)
	}

	cmd := exec.Command(s.shellBinPath)

	cmd.Stderr = io.Err
	cmd.Stdout = io.Out
	cmd.Stdin = io.In

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

package venv

import (
	"os/exec"

	"github.com/deifyed/xctl/pkg/apis/xctl"
)

type Shell interface {
	Command(io xctl.IOStreams, env []string) (*exec.Cmd, error)
}

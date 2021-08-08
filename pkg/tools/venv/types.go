package venv

import "os/exec"

type Shell interface {
	Command(env []string) (*exec.Cmd, error)
	Teardown() error
}

package zsh

import "github.com/spf13/afero"

type shell struct {
	fs *afero.Afero

	shellBinPath string

	tmpDir string
}

package binaries

import (
	"io"

	"github.com/spf13/afero"
)

type UnpackingFn func(io.Reader) (io.Reader, error)

type DownloadOpts struct {
	Name        string
	Version     string
	Fs          *afero.Afero
	BinariesDir string
	URL         string
	Hash        string
	UnpackingFn []UnpackingFn
}

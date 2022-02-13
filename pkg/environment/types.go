package environment

import (
	"io"

	"github.com/deifyed/xctl/pkg/cloud"
	"github.com/spf13/afero"
)

// ReconcileOpts defines required data for reconciling an environment
type ReconcileOpts struct {
	Out        io.Writer
	Err        io.Writer
	Filesystem *afero.Afero
	Provider   cloud.Provider
	Manifest   io.Reader
	Purge      bool
	Debug      bool
}

package environment

import (
	"context"
	"io"

	"github.com/spf13/afero"
)

// ReconcileOpts defines required data for reconciling an environment
type ReconcileOpts struct {
	Context    context.Context
	Out        io.Writer
	Err        io.Writer
	Filesystem *afero.Afero
	Manifest   io.Reader
	Purge      bool
}

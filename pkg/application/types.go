package application

import (
	"context"
	"io"

	"github.com/deifyed/xctl/pkg/cloud"
	"github.com/spf13/afero"
)

// ReconcileOpts defines required data for reconciling an application
type ReconcileOpts struct {
	Context             context.Context
	Out                 io.Writer
	Err                 io.Writer
	Filesystem          *afero.Afero
	Provider            cloud.Provider
	EnvironmentManifest io.Reader
	ApplicationManifest io.Reader
	Purge               bool
}

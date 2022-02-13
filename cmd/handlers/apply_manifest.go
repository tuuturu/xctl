package handlers

import (
	"bytes"
	"fmt"
	"io"

	"github.com/deifyed/xctl/pkg/application"
	"github.com/deifyed/xctl/pkg/cloud/linode"
	"github.com/deifyed/xctl/pkg/environment"

	"github.com/deifyed/xctl/pkg/apis/xctl"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

type ApplyRunEOpts struct {
	Filesystem *afero.Afero
	Io         xctl.IOStreams
	File       string
	Purge      bool
	Debug      bool
}

func ApplyRunE(opts *ApplyRunEOpts) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) (err error) {
		var originalManifestSource io.Reader

		if opts.File == "-" {
			originalManifestSource = opts.Io.In
		} else {
			originalManifestSource, err = opts.Filesystem.Open(opts.File)
			if err != nil {
				return fmt.Errorf("opening manifest file: %w", err)
			}
		}

		manifestSource := &bytes.Buffer{}
		tee := io.TeeReader(originalManifestSource, manifestSource)

		kind, err := v1alpha1.InferKindFromManifest(tee)
		if err != nil {
			return fmt.Errorf("inferring kind: %w", err)
		}

		provider := linode.NewLinodeProvider()

		err = provider.Authenticate()
		if err != nil {
			return fmt.Errorf("unable to authenticate: %w", err)
		}

		switch kind {
		case v1alpha1.ClusterKind:
			fmt.Fprintf(opts.Io.Out, "Applying cluster manifest, please wait\n\n")

			return environment.Reconcile(environment.ReconcileOpts{
				Out:        opts.Io.Out,
				Err:        opts.Io.Err,
				Filesystem: opts.Filesystem,
				Provider:   provider,
				Manifest:   manifestSource,
				Purge:      opts.Purge,
				Debug:      opts.Debug,
			})
		case v1alpha1.ApplicationKind:
			fmt.Fprintf(opts.Io.Out, "Applying application manifest %s, please wait\n\n", opts.File)

			return application.Reconcile(application.ReconcileOpts{
				Out:        opts.Io.Out,
				Filesystem: opts.Filesystem,
				Manifest:   manifestSource,
				Purge:      opts.Purge,
				Debug:      opts.Debug,
			})
		default:
			return fmt.Errorf("unknown kind %s", kind)
		}
	}
}

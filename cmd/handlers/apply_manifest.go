package handlers

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/deifyed/xctl/pkg/config"

	"github.com/deifyed/xctl/cmd/helpers"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

type ApplyRunEOpts struct {
	ExternalFilesystem *afero.Afero
	InternalFilesystem *afero.Afero
	Out                io.Writer
	File               string
	Purge              bool
}

func ApplyRunE(opts *ApplyRunEOpts) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) (err error) {
		var originalManifestSource io.Reader

		if opts.File == "-" {
			originalManifestSource = os.Stdin
		} else {
			originalManifestSource, err = opts.ExternalFilesystem.Open(opts.File)
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

		switch kind {
		case v1alpha1.ClusterKind:
			err = helpers.CopyToFs(
				opts.ExternalFilesystem,
				opts.InternalFilesystem,
				opts.File,
				config.GetAbsoluteInternalClusterManifestPath(),
			)
			if err != nil {
				return fmt.Errorf("copying cluster manifest to internal fs: %w", err)
			}

			fmt.Fprintf(opts.Out, "Applying cluster manifest, please wait\n\n")

			return handleCluster(opts.InternalFilesystem, opts.Out, opts.Purge, manifestSource)
		case v1alpha1.ApplicationKind:
			fmt.Fprintf(opts.Out, "Applying application manifest %s, please wait\n\n", opts.File)

			return handleApplication(opts.Out, opts.Purge, manifestSource)
		default:
			return fmt.Errorf("unknown kind %s", kind)
		}
	}
}

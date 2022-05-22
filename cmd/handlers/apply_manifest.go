package handlers

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/deifyed/xctl/pkg/application"
	"github.com/deifyed/xctl/pkg/cloud/linode"
	"github.com/deifyed/xctl/pkg/environment"

	"github.com/deifyed/xctl/pkg/apis/xctl"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

type ApplyRunEOpts struct {
	Filesystem         *afero.Afero
	Io                 xctl.IOStreams
	EnvironmentContext string
	File               string
	Purge              bool
}

func ApplyRunE(opts *ApplyRunEOpts) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) (err error) {
		var manifestSource io.Reader

		if opts.File == "-" {
			manifestSource = opts.Io.In
		} else {
			manifestSource, err = opts.Filesystem.Open(opts.File)
			if err != nil {
				return fmt.Errorf("opening manifest file: %w", err)
			}
		}

		rawManifest, err := io.ReadAll(manifestSource)
		if err != nil {
			return fmt.Errorf("buffering: %w", err)
		}

		kind, err := v1alpha1.InferKindFromManifest(bytes.NewReader(rawManifest))
		if err != nil {
			return fmt.Errorf("inferring kind: %w", err)
		}

		manifest := bytes.NewReader(rawManifest)

		provider := linode.NewLinodeProvider()

		err = provider.Authenticate()
		if err != nil {
			return fmt.Errorf("unable to authenticate: %w", err)
		}

		fmt.Fprintf(opts.Io.Out, "Applying %s manifest, please wait\n\n", strings.ToLower(kind))

		switch kind {
		case v1alpha1.EnvironmentKind:
			return environment.Reconcile(environment.ReconcileOpts{
				Out:        opts.Io.Out,
				Err:        opts.Io.Err,
				Filesystem: opts.Filesystem,
				Provider:   provider,
				Manifest:   manifest,
				Purge:      opts.Purge,
			})
		case v1alpha1.ApplicationKind:
			return application.Reconcile(application.ReconcileOpts{
				Out:                 opts.Io.Out,
				Filesystem:          opts.Filesystem,
				ApplicationManifest: manifest,
				Purge:               opts.Purge,
			})
		default:
			return fmt.Errorf("unknown kind %s", kind)
		}
	}
}

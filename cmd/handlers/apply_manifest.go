package handlers

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/deifyed/xctl/pkg/tools/paths"

	"github.com/deifyed/xctl/pkg/application"
	"github.com/deifyed/xctl/pkg/environment"

	"github.com/deifyed/xctl/pkg/apis/xctl"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func ApplyRunE(opts *ApplyRunEOpts) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) (err error) {
		kind, manifest, err := interpretManifest(opts.Filesystem, opts.Io.In, opts.File)
		if err != nil {
			return fmt.Errorf("interpreting manifest: %w", err)
		}

		fmt.Fprintf(opts.Io.Out, "Applying %s manifest, please wait\n\n", strings.ToLower(kind))

		absoluteRepositoryRootDirectory, err := paths.AbsoluteRepositoryRootDirectory()
		if err != nil {
			return fmt.Errorf("acquiring repository root directory: %w", err)
		}

		switch kind {
		case v1alpha1.EnvironmentKind:
			return environment.Reconcile(environment.ReconcileOpts{
				Context:    cmd.Context(),
				Out:        cmd.OutOrStdout(),
				Err:        cmd.OutOrStderr(),
				Filesystem: opts.Filesystem,
				Manifest:   manifest,
				Purge:      opts.Purge,
			})
		case v1alpha1.ApplicationKind:
			return application.Reconcile(application.ReconcileOpts{
				Context:                 cmd.Context(),
				Out:                     cmd.OutOrStdout(),
				Err:                     cmd.OutOrStderr(),
				Filesystem:              opts.Filesystem,
				RepositoryRootDirectory: absoluteRepositoryRootDirectory,
				ApplicationManifest:     manifest,
				Purge:                   opts.Purge,
			})
		default:
			return fmt.Errorf("unknown kind %s", kind)
		}
	}
}

func interpretManifest(fs *afero.Afero, in io.Reader, filepath string) (kind string, manifest io.Reader, err error) {
	var manifestSource io.Reader

	if filepath == "-" {
		manifestSource = in
	} else {
		manifestSource, err = fs.Open(filepath)
		if err != nil {
			return "", nil, fmt.Errorf("opening manifest file: %w", err)
		}
	}

	rawManifest, err := io.ReadAll(manifestSource)
	if err != nil {
		return "", nil, fmt.Errorf("buffering: %w", err)
	}

	kind, err = v1alpha1.InferKindFromManifest(bytes.NewReader(rawManifest))
	if err != nil {
		return "", nil, fmt.Errorf("inferring kind: %w", err)
	}

	return kind, bytes.NewReader(rawManifest), nil
}

type ApplyRunEOpts struct {
	Filesystem         *afero.Afero
	Io                 xctl.IOStreams
	EnvironmentContext string
	File               string
	Purge              bool
}

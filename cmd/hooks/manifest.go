package hooks

import (
	"fmt"
	"io"
	"os"

	"github.com/deifyed/xctl/pkg/tools/manifests"

	"github.com/deifyed/xctl/pkg/apis/xctl"

	"sigs.k8s.io/yaml"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

// EnvironmentManifestInitializerOpts defines required data for loading a cluster manifest
type EnvironmentManifestInitializerOpts struct {
	Io xctl.IOStreams
	Fs *afero.Afero

	EnvironmentManifest *v1alpha1.Environment

	SourcePath *string
}

// EnvironmentManifestInitializer initializes a cluster manifest struct based on a filepath or stdin
func EnvironmentManifestInitializer(opts EnvironmentManifestInitializerOpts) func(cmd *cobra.Command, args []string) error { //nolint:lll
	return func(cmd *cobra.Command, args []string) error {
		var (
			source io.Reader
			err    error
		)

		if *opts.SourcePath == "-" {
			source = opts.Io.In
		} else {
			source, err = opts.Fs.OpenFile(*opts.SourcePath, os.O_RDONLY, 0o755)
			if err != nil {
				return fmt.Errorf("opening cluster manifest: %w", err)
			}
		}

		rawManifest, err := io.ReadAll(source)
		if err != nil {
			return fmt.Errorf("reading cluster manifest: %w", err)
		}

		defaultSource, err := manifests.ResourceAsReader(v1alpha1.NewDefaultEnvironment())
		if err != nil {
			return fmt.Errorf("converting manifest to stream: %w", err)
		}

		rawDefaultSource, err := io.ReadAll(defaultSource)
		if err != nil {
			return fmt.Errorf("buffering default source: %w", err)
		}

		err = yaml.Unmarshal(rawDefaultSource, opts.EnvironmentManifest)
		if err != nil {
			return fmt.Errorf("unmarshalling default manifest: %w", err)
		}

		err = yaml.Unmarshal(rawManifest, opts.EnvironmentManifest)
		if err != nil {
			return fmt.Errorf("unmarshalling cluster manifest: %w", err)
		}

		return nil
	}
}

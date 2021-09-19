package preruns

import (
	"fmt"
	"io"
	"os"

	"github.com/deifyed/xctl/pkg/apis/xctl"

	"sigs.k8s.io/yaml"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

type ClusterManifestIniterOpts struct {
	Io xctl.IOStreams
	Fs *afero.Afero

	ClusterManifest *v1alpha1.Cluster

	SourcePath *string
}

// ClusterManifestIniter initializes a cluster manifest struct based on a filepath or stdin
func ClusterManifestIniter(opts ClusterManifestIniterOpts) func(cmd *cobra.Command, args []string) error { //nolint:lll
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

		err = yaml.Unmarshal(rawManifest, opts.ClusterManifest)
		if err != nil {
			return fmt.Errorf("unmarshalling cluster manifest: %w", err)
		}

		return nil
	}
}

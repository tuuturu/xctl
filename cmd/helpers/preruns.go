package helpers

import (
	"fmt"
	"io"
	"os"

	"sigs.k8s.io/yaml"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

// ClusterManifestIniter initializes a cluster manifest struct based on a filepath or stdin
func ClusterManifestIniter(fs *afero.Afero, sourcePath *string, clusterManifest *v1alpha1.Cluster) func(cmd *cobra.Command, args []string) error { //nolint:lll
	return func(cmd *cobra.Command, args []string) error {
		var (
			source io.Reader
			err    error
		)

		if *sourcePath == "-" {
			source = os.Stdin
		} else {
			source, err = fs.OpenFile(*sourcePath, os.O_RDONLY, 0o755)
			if err != nil {
				return fmt.Errorf("opening cluster manifest: %w", err)
			}
		}

		rawManifest, err := io.ReadAll(source)
		if err != nil {
			return fmt.Errorf("reading cluster manifest: %w", err)
		}

		err = yaml.Unmarshal(rawManifest, clusterManifest)
		if err != nil {
			return fmt.Errorf("unmarshalling cluster manifest: %w", err)
		}

		return nil
	}
}

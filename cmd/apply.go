package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/deifyed/xctl/cmd/handlers"
	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

type applyOpts struct {
	File string
}

var (
	applyCmdOpts applyOpts         //nolint:gochecknoglobals
	applyCmd     = &cobra.Command{ //nolint:gochecknoglobals
		Use:   "apply",
		Short: "applies a manifest",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			out := os.Stdout
			fs := &afero.Afero{Fs: afero.NewOsFs()}

			var manifestSource io.Reader

			if applyCmdOpts.File == "-" {
				manifestSource = os.Stdin
			} else {
				manifestSource, err = fs.Open(applyCmdOpts.File)
				if err != nil {
					return fmt.Errorf("opening manifest file: %w", err)
				}
			}

			kind, err := v1alpha1.InferKindFromManifest(manifestSource)
			if err != nil {
				return fmt.Errorf("inferring kind: %w", err)
			}

			switch kind {
			case v1alpha1.ClusterKind:
				fmt.Fprintf(out, "Applying cluster manifest, please wait\n\n")

				return handlers.HandleCluster(out, false, manifestSource)
			case v1alpha1.ApplicationKind:
				fmt.Fprintf(out, "Applying application manifest %s, please wait\n\n", applyCmdOpts.File)

				return handlers.HandleApplication(out, false, manifestSource)
			default:
				return fmt.Errorf("unknown kind %s", kind)
			}
		},
	}
)

//nolint:gochecknoinits
func init() {
	flags := applyCmd.Flags()

	flags.StringVarP(&applyCmdOpts.File, "file", "f", "-", "file to apply")

	rootCmd.AddCommand(applyCmd)
}
